// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

const nodeUrl = "https://nano-rpc.trustwalletapp.com"

// Return success, duration
func getWork(hash string) (bool, time.Duration) {
	log.Printf("Requesting work from node, hash %v \n", hash)
	// trigger work
	resp, err, duration := rpcclient.GetWork(nodeUrl, hash, 0)
	if err != nil {
		log.Printf("Work resp from node FAIL, dur %v, err %v \n", duration, err)
		return false, duration
	}
	log.Printf("Work resp from node ok, dur %v, resp %v \n", duration, resp)
	return true, duration
}

func pregenerate(hash string) {
	action := "work_pregenerate_by_hash"
	body := "{\"action\":\"" + action + "\", \"hash\": \"" + hash + "\"}"
	log.Printf("Pregenerating for hash %v \n", hash)
	rpcclient.MakeGenericCall(nodeUrl, body)
}

// Do a simulated account cycle:
// - request balance
// - wait some seconds
// - ge work
// - in case of failure, repeat until ok (or max 10 tries)
// Return
// - success
// - number of iterations (usually 1)
// - duration of last work get (usually single)
// - total duration
func makeAccountWorkCycle(hash string, initialDelay int) (bool, int, time.Duration, time.Duration) {
	log.Printf("Cycle: starting for hash %v \n", hash)
	timeStart := time.Now()
	pregenerate(hash)

	log.Printf("Cycle: sleep %v \n", initialDelay)
	time.Sleep(time.Duration(initialDelay) * time.Second)

	success := false
	stepCount := 0
	var duration1 time.Duration = 0
	var res bool = false
	for stepCount = 0; stepCount < 10; {
		stepCount++
		log.Printf("Cycle: step %v \n", stepCount)
		res, duration1 = getWork(hash)
		if res {
			success = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	timeStop := time.Now()
	durationTotal := timeStop.Sub(timeStart)
	//log.Printf("Cycle starting for hash %v \n", hash)
	log.Printf("Cycle: stop, success %v, steps %v, duration1 %v, durationTotal %v \n", success, stepCount, duration1, durationTotal)
	return success, stepCount, duration1, durationTotal
}

// Return a random hash (64 hex digits)
func getRandomHash() string {
	hash := ""
	for i := 0; i < 32; i++ {
		byte := rand.Intn(256)
		hexByte := strconv.FormatInt(int64(byte), 16)
		if len(hexByte) < 2 {
			hexByte = "0" + hexByte
		}
		hash += hexByte
	}
	hash = strings.ToUpper(hash)
	return hash
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	hash := getRandomHash()
	makeAccountWorkCycle(hash, 6)
}
