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
var delayBetweenCommandsMs time.Duration = 300

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
	hash := randomHash()
	go rpcCallPrint(url, fmt.Sprintf(`{"action": "work_generate","hash": "%v"}`, hash), true)
}

func runPregenerateByHash() {
	hash := randomHash()
	rpcCallPrint(url, fmt.Sprintf(`{"action": "work_pregenerate_by_hash","hash": "%v"}`, hash), false)
}

func runCommand() {
	var commandTypeRandom int = rand.Intn(100)
	// 10% work_generate
	commandTypeRandom -= 10
	if commandTypeRandom <= 0 {
		runWorkGenerate()
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
