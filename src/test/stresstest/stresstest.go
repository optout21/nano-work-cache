// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

var url string = "http://localhost:7376"
var testDataHashCount int = 300
var delayBetweenCommandsMs time.Duration = 250

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func rpcCallPrint(url string, reqJson string, printResult bool) {
	fmt.Println("Calling url", url, "with data", reqJson)
	startTime := time.Now()
	resp, err := rpcclient.RpcCall(url, reqJson)
	if (err != nil) {
		fmt.Println("Error:", err)
		return
	}
	endTime := time.Now()
	if printResult {
		fmt.Printf("Response:  dur %v \n", endTime.Sub(startTime))
		fmt.Println(resp)
	}
}

func runWorkGenerate() {
	hash := getRandomHash()
	go rpcCallPrint(url, fmt.Sprintf(`{"action": "work_generate","hash": "%v"}`, hash), true)
}

func runPregenerateByHash() {
	hash := getRandomHash()
	rpcCallPrint(url, fmt.Sprintf(`{"action": "work_pregenerate_by_hash","hash": "%v"}`, hash), false)
}

func runAccountBalance() {
	account := getValidAccount()
	go rpcCallPrint(url, fmt.Sprintf(`{"action": "account_balance","account": "%v"}`, account), true)
}

func runAccountsBalances() {
	account := getValidAccount()
	go rpcCallPrint(url, fmt.Sprintf(`{"action": "accounts_balances","accounts": ["%v"]}`, account), true)
}

func runCommand() {
	var commandTypeRandom int = rand.Intn(100)

	// 10% work_generate
	commandTypeRandom -= 10
	if commandTypeRandom <= 0 {
		runWorkGenerate()
		return
	}

	// 3% accounts_balances
	commandTypeRandom -= 3
	if commandTypeRandom <= 0 {
		runAccountsBalances()
		return
	}

	// 3% account_balance
	commandTypeRandom -= 3
	if commandTypeRandom <= 0 {
		runAccountBalance()
		return
	}
	
	// rest: work_pregenerate_by_hash
	runPregenerateByHash()
}

func main() {
	initTestData(testDataHashCount)

	for {
		runCommand()
		time.Sleep(delayBetweenCommandsMs * time.Millisecond)
	}
}
