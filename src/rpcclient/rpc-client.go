// Copyright © 2019-2019 catenocrypt.  See LICENSE file for license information.

package rpcclient

import (
	"net/http"
	"bytes"
	"strconv"
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

// difficulty may be missing (0)
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
