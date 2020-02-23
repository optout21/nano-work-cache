// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package rpcclient

import (
	"errors"
	//"fmt"
	"net/http"
	"bytes"
	"strconv"
	"strings"
	"encoding/json"
	"io/ioutil"
)

func RpcCall(url string, reqJson string) (respJson string, err error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(reqJson))
	if (err != nil) {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if (err != nil) {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		return "", err
	}
	return string(body), nil
}

type WorkResponse struct {
	Hash string
	Work string
	Difficulty uint64
	Multiplier float64
}

type WorkResponseJson struct {
	Hash string
	Work string
	Difficulty string
	Multiplier string
}

// work_generate.  Difficulty may be missing (0)
func GetWork(url string, hash string, diff uint64) (WorkResponse, error) {
	reqJson := `{"action":"work_generate","hash":"` + hash + `"`
	if (diff != 0) {
		reqJson += `,"difficulty":"` + strconv.FormatUint(diff, 16) + `"`
	}
	reqJson += `}`
	//fmt.Println(reqJson)
	respString, err := RpcCall(url, reqJson)
	var resp WorkResponse
	if (err != nil) {
		return resp, err
	}
	// parse json
	//fmt.Println(respString)
	var respStruct1 WorkResponseJson
	err = json.Unmarshal([]byte(respString), &respStruct1)
	if (err != nil) {
		return resp, err
	}
	//fmt.Println(respStruct1)
	difficulty, err := strconv.ParseUint(respStruct1.Difficulty, 16, 64)
	if (err != nil) {
		// diff not present, take input (in reality actual difficulty is usually higher)
		difficulty = diff
	}
	multiplier, err := strconv.ParseFloat(respStruct1.Multiplier, 64)
	if (err != nil) {
		// could not read multiplier
		multiplier = 1.0
	}
	resp = WorkResponse{respStruct1.Hash, respStruct1.Work, difficulty, multiplier}
	return resp, nil
}

type AccountFrontiersRespJson struct {
	Frontiers map[string]string
}

// Get frontier blocks for accounts, accounts_frontiers
func GetFrontiers(url string, accounts []string) (map[string]string, error) {
	reqJson := `{"action":"accounts_frontiers","accounts":["` + strings.Join(accounts[:], `","`) + `"]}`
	//fmt.Println(reqJson)
	respString, err := RpcCall(url, reqJson)
	if (err != nil) {
		return nil, err
	}
	// parse json
	//fmt.Println(respString)
	var respStruct1 AccountFrontiersRespJson
	err = json.Unmarshal([]byte(respString), &respStruct1)
	if (err != nil) {
		return nil, err
	}
	//fmt.Println(respStruct1)
	return respStruct1.Frontiers, nil
}

// Get frontier block for an account, using accounts_frontiers
func GetFrontier(url string, account string) (string, error) {
	accounts, err := GetFrontiers(url, []string{account})
	if (err != nil) { return "", err }
	frontier := accounts[account]
	if (len(frontier) == 0) {
		return "", errors.New("Could not find account in accounts_frontiers")
	}
	return frontier, nil
}

/// Make a generic call to the RPC node
func MakeGenericCall(url string, reqJSON string) (string, error) {
	//fmt.Println(reqJson)
	respString, err := RpcCall(url, reqJSON)
	if (err != nil) {
		return "", err
	}
	return respString, nil;
}
