// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

const nodeUrl = "https://nano-rpc.trustwalletapp.com"
const MonitorFreq int = 300     // sec
const TimeFromBalanceToWork = 6 // sec

type Stats struct {
	Success        bool
	IterationCount int
	DurationLast   time.Duration
	DurationTotal  time.Duration
}

func (s *Stats) ToString() string {
	return fmt.Sprintf("success %v  iterCount %v  durLast %v  durTotal %v", s.Success, s.IterationCount, s.DurationLast, s.DurationTotal)
}

type CumulStats struct {
	Count          int
	Success        int
	IterationCount int
	DurationLast   time.Duration
	DurationTotal  time.Duration
}

func (s *CumulStats) Add(s2 Stats) {
	s.Count++
	if s2.Success {
		s.Success++
	}
	s.IterationCount += s2.IterationCount
	s.DurationLast += s2.DurationLast
	s.DurationTotal += s2.DurationTotal
}

func (s *CumulStats) ToString() string {
	if s.Count <= 0 {
		return "empty"
	}
	return fmt.Sprintf("cnt %v  success %v  iterCount %v  durLast %v  durTotal %v",
		s.Count,
		s.Success,
		float32(s.IterationCount)/float32(s.Count),
		float32(s.DurationLast.Milliseconds())/float32(s.Count),
		float32(s.DurationTotal.Milliseconds())/float32(s.Count))
}

// Return success, duration
func getWork(hash string) (bool, time.Duration) {
	log.Printf("Requesting work from node, hash %v \n", hash)
	// trigger work
	resp, err, duration := rpcclient.GetWork(hash, 0)
	if err != nil {
		log.Printf("Work resp from node FAIL, dur %v, err %v \n", duration, err)
		return false, duration
	}
	if len(resp.Work) < 10 {
		log.Printf("Work resp has invalid work, work %v, dur %v, resp %v \n", resp.Work, duration, resp)
		return false, duration
	}
	log.Printf("Work resp from node ok, dur %v, resp %v \n", duration, resp)
	return true, duration
}

func pregenerate(hash string) {
	action := "work_pregenerate_by_hash"
	body := "{\"action\":\"" + action + "\", \"hash\": \"" + hash + "\"}"
	log.Printf("Pregenerating for hash %v \n", hash)
	rpcclient.MakeGenericCall(body)
}

// Do a simulated account cycle:
// - request balance
// - wait some seconds
// - ge work
// - in case of failure, repeat until ok (or max 10 tries)
// Return Stats
func makeCycle(hash string) Stats {
	log.Printf("Cycle: starting for hash %v \n", hash)
	timeStart := time.Now()
	pregenerate(hash)

	toSleep := timeStart.Add(time.Duration(TimeFromBalanceToWork) * time.Second).Sub(time.Now())
	log.Printf("Cycle: sleep %v \n", toSleep)
	time.Sleep(toSleep)

	stats := Stats{false, 0, 0, 0}
	timeStartWork := time.Now()
	for stats.IterationCount = 0; stats.IterationCount < 10; {
		stats.IterationCount++
		log.Printf("Cycle: step %v \n", stats.IterationCount)
		res, dur1 := getWork(hash)
		if res {
			stats.Success = true
			stats.DurationLast = dur1
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	timeStop := time.Now()
	stats.DurationTotal = timeStop.Sub(timeStartWork)
	//log.Printf("Cycle starting for hash %v \n", hash)
	log.Printf("STAT1 %v \n", stats.ToString())
	return stats
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
	rpcclient.Init(nodeUrl, nodeUrl)

	nextCheck := time.Now()
	var cumulative CumulStats = CumulStats{}
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		now := time.Now()
		untilNextCheck := nextCheck.Sub(now).Milliseconds()
		if untilNextCheck > 0 {
			log.Printf("Sleeping %v ms ...", untilNextCheck)
			time.Sleep(time.Duration(untilNextCheck) * time.Millisecond)
		}
		nextCheck = nextCheck.Add(time.Duration(MonitorFreq) * time.Second)

		hash := getRandomHash()
		stats := makeCycle(hash)

		cumulative.Add(stats)
		log.Printf("STAT_CUMUL %v \n", cumulative.ToString())
	}
}
