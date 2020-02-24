// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"encoding/json"
	"fmt"
    "log"
	"net/http"
	"strconv"
	"github.com/catenocrypt/nano-work-cache/workcache"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

type workGenerateJson struct {
	Action string
	Hash string
	Difficulty string
}

type workPregenerateByHashJson struct {
	Action string
	Hash string
}

type workPregenerateByAccountJson struct {
	Action string
	Account string
}

type accountBalanceJson struct {
	Action string
	Account string
}

type accountsBalancesJson struct {
	Action string
	Accounts []string
}

/// Not the normal Json Encode way, due to the difficult hex formatting.  Using simple string concatenation.
func workResponseToJson(resp workcache.WorkResponse) string {
	return fmt.Sprintf(`{"hash":"%v","work":"%v","difficulty":"%x","multiplier":"%v","source":"%v"}`,
		resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, resp.Source)
}

func handleReqSync(action string, respBody []byte, w http.ResponseWriter) {
	switch action {
	case "work_generate":
		var workGenerate workGenerateJson
		err := json.Unmarshal(respBody, &workGenerate)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_generate parse error"}`)
			return
		}
		log.Println("work_generate req", workGenerate)
		var difficulty uint64 = workcache.GetDefaultDifficulty()
		if (len(workGenerate.Difficulty) > 0) {
			difficultyParsed, err := strconv.ParseUint(workGenerate.Difficulty, 16, 64)
			if (err != nil) {
				// diff present, but could not parse
				fmt.Fprintln(w, `{"error":"work_generate difficulty parse error"}`)
				return
			}
			difficulty = difficultyParsed
		}
		// handle
		workResp, err := workcache.GetCachedWork(nanoNodeUrl, workGenerate.Hash, difficulty, "")
		if (err != nil) {
			fmt.Fprintln(w, `{"error":"` + err.Error() + `"}`)
			return
		}
		log.Println("work_generate resp", workResp)
		fmt.Fprintln(w, workResponseToJson(workResp))
		break

	case "work_pregenerate_by_hash":
		var workPregenerateByHash workPregenerateByHashJson
		err := json.Unmarshal(respBody, &workPregenerateByHash)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_pregenerate_by_hash parse error"}`)
			return
		}
		log.Println("work_pregenerate_by_hash req", workPregenerateByHash)
		var hash = workPregenerateByHash.Hash
		var difficulty uint64 = workcache.GetDefaultDifficulty()
		// start work asynchronously
		go workcache.GetCachedWork(nanoNodeUrl, hash, difficulty, "")
		// return response, only hash
		fmt.Fprintln(w, fmt.Sprintf(`{"hash":"%v","source":"started_in_background"}`, hash))
		return

	case "work_pregenerate_by_account":
		var workPregenerateByAccount workPregenerateByAccountJson
		err := json.Unmarshal(respBody, &workPregenerateByAccount)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_pregenerate_by_account parse error"}`)
			return
		}
		log.Println("work_pregenerate_by_account req", workPregenerateByAccount)
		var account = workPregenerateByAccount.Account
		// get frontier and start work asynchronously
		go workcache.GetCachedWorkByAccount(nanoNodeUrl, account)
		// return response, only account is echoed back, hash and work is not available yet
		fmt.Fprintln(w, fmt.Sprintf(`{"account":"%v","hash":"<retrieve_in_progress>","source":"started_in_background"}`, account))
		return

	case "account_balance":
		// account_balance also triggers work_precompute in the background, and transparently proxies the call for balance
		var accountBalance accountBalanceJson
		err := json.Unmarshal(respBody, &accountBalance)
		if err != nil {
			fmt.Fprintln(w, `{"error":"account_balance parse error"}`)
			return
		}
		log.Println("account_balance", accountBalance)

		// get frontier and start work asynchronously
		go workcache.GetCachedWorkByAccount(nanoNodeUrl, accountBalance.Account)

		// proxy the call
		respJSON, err := rpcclient.MakeGenericCall(nanoNodeUrl, string(respBody))
		if (err != nil) {
			log.Println("RPC error:", err.Error())
			fmt.Fprintln(w, `{"error":"RPC error: ` + err.Error() + `","action":"` + action + `"}`)
			return
		}
		fmt.Fprintln(w, respJSON)
		break
		
	case "accounts_balances":
		// accounts_balances also triggers work_precompute (for all accounts) in the background, and transparently proxies the call for balances
		var accountsBalances accountsBalancesJson
		err := json.Unmarshal(respBody, &accountsBalances)
		if err != nil {
			fmt.Fprintln(w, `{"error":"accounts_balances parse error"}`)
			return
		}
		log.Println("accounts_balances", accountsBalances)

		// for all accounts get frontier and start work asynchronously
		for _, account := range accountsBalances.Accounts {
			go workcache.GetCachedWorkByAccount(nanoNodeUrl, account)
		}

		// proxy the call
		respJSON, err := rpcclient.MakeGenericCall(nanoNodeUrl, string(respBody))
		if (err != nil) {
			log.Println("RPC error:", err.Error())
			fmt.Fprintln(w, `{"error":"RPC error: ` + err.Error() + `","action":"` + action + `"}`)
			return
		}
		fmt.Fprintln(w, respJSON)
		break

	case "nano-work-cache-status-internal":
		status := getStatus()
		fmt.Fprintln(w, status)
		break

	default:
		// proxy any other request unmodified
		log.Println("transaprent proxying of action", action)
		respJSON, err := rpcclient.MakeGenericCall(nanoNodeUrl, string(respBody))
		if (err != nil) {
			log.Println("RPC error:", err.Error())
			fmt.Fprintln(w, `{"error":"RPC error: ` + err.Error() + `","action":"` + action + `"}`)
			return
		}
		fmt.Fprintln(w, respJSON)
	}
}
