// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package rpcclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	WorkResponse struct {
		Hash       string
		Work       string
		Difficulty uint64
		Multiplier float64
	}

	WorkResponseJson struct {
		Hash       string
		Work       string
		Difficulty string
		Multiplier string
	}

	AccountFrontiersRespJson struct {
		Frontiers map[string]string
	}

	ActiveDifficultyRespJson struct {
		NetworkMinimum string `json:"network_minimum"`
		NetworkCurrent string `json:"network_current"`
		Multiplier     string
	}
)

var rpcUrl string = "?"
var rpcWorkUrl string = "?"

func Init(rpcUrlIn string, rpcWorkUrlIn string) {
	rpcUrl = rpcUrlIn
	rpcWorkUrl = rpcWorkUrlIn
}

func RpcCall(url string, reqJson string) (respJson string, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(reqJson))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// work_generate.  Difficulty may be missing (0)
func GetWork(hash string, diff uint64) (WorkResponse, error, time.Duration) {
	timeStart := time.Now()
	reqJson := fmt.Sprintf(`{"action":"work_generate","hash":"%v"`, hash)
	reqJson += `,"use_peers":"true"`
	if diff != 0 {
		reqJson += fmt.Sprintf(`,"difficulty":"%x"`, diff)
	}
	reqJson += `}`
	log.Printf("Requesting work, from %v, %v \n", rpcWorkUrl, reqJson)
	respString, err := RpcCall(rpcWorkUrl, reqJson)
	var resp WorkResponse
	if err != nil {
		return resp, err, 0
	}
	// parse json
	var respStruct1 WorkResponseJson
	err = json.Unmarshal([]byte(respString), &respStruct1)
	if err != nil {
		return resp, err, 0
	}
	difficulty, err := strconv.ParseUint(respStruct1.Difficulty, 16, 64)
	if err != nil {
		// diff not present, take input (in reality actual difficulty is usually higher)
		difficulty = diff
	}
	multiplier, err := strconv.ParseFloat(respStruct1.Multiplier, 64)
	if err != nil {
		// could not read multiplier
		multiplier = 1.0
	}
	timeStop := time.Now()
	resp = WorkResponse{respStruct1.Hash, respStruct1.Work, difficulty, multiplier}
	return resp, nil, timeStop.Sub(timeStart)
}

// Get frontier blocks for accounts, accounts_frontiers
func GetFrontiers(accounts []string) (map[string]string, error) {
	reqJson := `{"action":"accounts_frontiers","accounts":["` + strings.Join(accounts[:], `","`) + `"]}`
	//fmt.Println(reqJson)
	respString, err := RpcCall(rpcUrl, reqJson)
	if err != nil {
		return nil, err
	}
	// parse json
	//fmt.Println(respString)
	var respStruct1 AccountFrontiersRespJson
	err = json.Unmarshal([]byte(respString), &respStruct1)
	if err != nil {
		return nil, err
	}
	//fmt.Println(respStruct1)
	return respStruct1.Frontiers, nil
}

// GetDifficulty Get current level of difficulty
func GetDifficulty() (string, error) {
	reqJson := `{"action": "active_difficulty"}`
	respString, err := RpcCall(rpcUrl, reqJson)
	//fmt.Println("reqJson %v respString %v \n", reqJson, respString)
	if err != nil {
		return "", err
	}
	// parse json
	var respStruct1 ActiveDifficultyRespJson
	err = json.Unmarshal([]byte(respString), &respStruct1)
	if err != nil {
		return "", err
	}
	//fmt.Println(respStruct1)
	return respStruct1.NetworkCurrent, nil
}

// Get frontier block for an account, using accounts_frontiers
func GetFrontier(account string) (string, error) {
	accounts, err := GetFrontiers([]string{account})
	if err != nil {
		return "", err
	}
	frontier := accounts[account]
	if len(frontier) == 0 {
		return "", errors.New("Could not find account in accounts_frontiers")
	}
	return frontier, nil
}

/// Make a generic call to the RPC node
func MakeGenericCall(reqJSON string) (string, error) {
	//fmt.Println(reqJson)
	respString, err := RpcCall(rpcUrl, reqJSON)
	if err != nil {
		return "", err
	}
	return respString, nil
}
