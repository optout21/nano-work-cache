// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/catenocrypt/nano-work-cache/rpcclient"
	"github.com/catenocrypt/nano-work-cache/workcache"
)

type blockJson struct {
	Account string
}

type workGenerateJson struct {
	Action     string
	Hash       string
	Difficulty string
}

type workPregenerateByHashJson struct {
	Action string
	Hash   string
}

type workPregenerateByAccountJson struct {
	Action  string
	Account string
}

type accountBalanceJson struct {
	Action  string
	Account string
}

type accountsBalancesJson struct {
	Action   string
	Accounts []string
}

type requestWithBlockJson struct {
	Block blockJson
}

type responseWithHashJson struct {
	Hash string
}

var enablePregeneration int = 1

/// Not the normal Json Encode way, due to the difficulty hex formatting.  Using simple string concatenation.
func workResponseToJson(resp workcache.WorkResponse) string {
	return fmt.Sprintf(`{"hash":"%v","work":"%v","difficulty":"%x","multiplier":"%v","source":"%v"}`,
		resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, resp.Source)
}

/// Proxy an incoming call to the node unmodified
func proxyCall(action string, req string) (string, error) {
	//log.Println("transparent proxying of action", action)
	respJSON, err := rpcclient.MakeGenericCall(req)
	if err != nil {
		log.Println("RPC error:", err.Error())
		return "", err
	}
	return respJSON, nil
}

func handleReqSync(action string, reqBody []byte, w http.ResponseWriter) {
	switch action {
	case "work_generate":
		var workGenerate workGenerateJson
		err := json.Unmarshal(reqBody, &workGenerate)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_generate parse error"}`)
			return
		}
		log.Println("work_generate req", workGenerate)
		var difficulty uint64 = rpcclient.GetDifficultyCached()
		if len(workGenerate.Difficulty) > 0 {
			difficultyParsed, err := strconv.ParseUint(workGenerate.Difficulty, 16, 64)
			if err != nil {
				// diff present, but could not parse
				fmt.Fprintln(w, `{"error":"work_generate difficulty parse error"}`)
				return
			}
			difficulty = difficultyParsed
		}
		// handle
		workResp, err := workcache.Generate(workGenerate.Hash, difficulty, "")
		log.Println("work_generate resp", workResp)
		if err != nil {
			fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
			return
		}
		fmt.Fprintln(w, workResponseToJson(workResp))
		break

	case "work_pregenerate_by_hash":
		var workPregenerateByHash workPregenerateByHashJson
		err := json.Unmarshal(reqBody, &workPregenerateByHash)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_pregenerate_by_hash parse error"}`)
			return
		}
		log.Println("work_pregenerate_by_hash req", workPregenerateByHash)
		var hash = workPregenerateByHash.Hash
		// start pregenerate asynchronously, regardless of enable flag
		workcache.PregenerateByHash(hash, "")
		// return response, only hash
		fmt.Fprintln(w, fmt.Sprintf(`{"hash":"%v","source":"started_in_background"}`, hash))
		return

	case "work_pregenerate_by_account":
		var workPregenerateByAccount workPregenerateByAccountJson
		err := json.Unmarshal(reqBody, &workPregenerateByAccount)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_pregenerate_by_account parse error"}`)
			return
		}
		log.Println("work_pregenerate_by_account req", workPregenerateByAccount)
		var account = workPregenerateByAccount.Account
		// get frontier of account
		hash, err := workcache.GetFrontierHash(account)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf(`{"error":"%v"}`, err.Error()))
			return
		}
		// pregenerate work asynchronously, regardless of enable flag
		workcache.PregenerateByHash(hash, account)
		// return response; account is echoed back; hash is returned; work is not available yet
		fmt.Fprintln(w, fmt.Sprintf(`{"account":"%v","hash":"%v","source":"started_in_background"}`, account, hash))
		return

	case "account_balance":
		// account_balance also triggers work_precompute in the background, and transparently proxies the call for balance
		var accountBalance accountBalanceJson
		err := json.Unmarshal(reqBody, &accountBalance)
		if err != nil {
			fmt.Fprintln(w, `{"error":"account_balance parse error"}`)
			return
		}
		//log.Println("account_balance", accountBalance)

		if enablePregeneration >= 1 {
			// get frontier and pregenerate work asynchronously
			workcache.PregenerateByAccount(accountBalance.Account)
		}

		// proxy the call
		respJSON, err := proxyCall(action, string(reqBody))
		if err != nil {
			fmt.Fprintln(w, `{"error":"RPC error: `+err.Error()+`","action":"`+action+`"}`)
			return
		}
		fmt.Fprintln(w, respJSON)
		break

	case "accounts_balances":
		// accounts_balances also triggers work_precompute (for all accounts) in the background, and transparently proxies the call for balances
		var accountsBalances accountsBalancesJson
		err := json.Unmarshal(reqBody, &accountsBalances)
		if err != nil {
			fmt.Fprintln(w, `{"error":"accounts_balances parse error"}`)
			return
		}
		//log.Println("accounts_balances", accountsBalances)

		if enablePregeneration >= 1 {
			// for all accounts get frontier and pregenerate work asynchronously
			for _, account := range accountsBalances.Accounts {
				workcache.PregenerateByAccount(account)
			}
		}

		// proxy the call
		respJSON, err := proxyCall(action, string(reqBody))
		if err != nil {
			fmt.Fprintln(w, `{"error":"RPC error: `+err.Error()+`","action":"`+action+`"}`)
			return
		}
		fmt.Fprintln(w, respJSON)
		break

	case "nano-work-cache-status-internal":
		status := getStatus()
		fmt.Fprintln(w, status)
		break

	case "block_create":
	case "block_hash":
	case "process":
		// proxy these calls unmodified, but watch the hash in the result, and trigger work computation for it in the background
		// first try to obtain account from the request
		var account string = ""
		var requestWithBlock requestWithBlockJson
		err := json.Unmarshal(reqBody, &requestWithBlock)
		if err == nil {
			account = requestWithBlock.Block.Account
			log.Println("Extracted account from request action", action, "account", account)
		}

		respJSON, err := proxyCall(action, string(reqBody))
		if err != nil {
			fmt.Fprintln(w, `{"error":"RPC error: `+err.Error()+`","action":"`+action+`"}`)
			return
		}

		// read out hash from the response
		var responseWithHash responseWithHashJson
		err = json.Unmarshal([]byte(respJSON), &responseWithHash)
		if err != nil {
			log.Println("Warning: Error reading hash from response of" + action)
		} else {
			if enablePregeneration >= 1 {
				// we have the hash, trigger work computation
				hash := responseWithHash.Hash
				if len(hash) > 0 {
					log.Println("Reqesting work from action", action, "for hash", hash, "and account", account)
					workcache.PregenerateByHash(hash, account)
				}
			}
		}
		fmt.Fprintln(w, respJSON)

	default:
		// proxy any other request unmodified
		respJSON, err := proxyCall(action, string(reqBody))
		if err != nil {
			fmt.Fprintln(w, `{"error":"RPC error: `+err.Error()+`","action":"`+action+`"}`)
			return
		}
		fmt.Fprintln(w, respJSON)
	}
}
