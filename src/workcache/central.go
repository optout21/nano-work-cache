// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"errors"
	"log"
	"time"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

const (
	WorkInputHash = 0
	WorkInputAccount = 1
)

type WorkRequest struct {
	Url string
	Input int // WorkInputHash or Account
	Hash string
	Diff uint64
	Account string
}

type WorkResponse struct {
	Hash string
	Work string
	Difficulty uint64
	Multiplier float64
	// values: 'fresh', 'cache'
	Source string
	Error error
}

var maxOutRequests int = 8;
var statusWorkOutReqCount int = 0
var statusWorkOutRespCount int = 0
var statusWorkInReqCount int = 0
var statusWorkInReqFromCache int = 0

// Start Invoked at the beginning, can perform initializations, read the cache, etc.
func Start(backgroundWorkerCount int, maxOutRequestsIn int) {
	LoadCache()
	maxOutRequests = maxOutRequestsIn
	startWorkers(backgroundWorkerCount)
	go housekeepingCycle()
}

// Generate Generate work, in foreground, but for rate limiting and priority handling it goes to a pool worker.
// Account is optional, may by empty. 
// Difficulty may be 0, default will be used
func Generate(url string, hash string, difficulty uint64, account string) (WorkResponse, error) {
	req := WorkRequest{url, WorkInputHash, hash, difficulty, account}
	resp := getCachedWork(req)
	return resp, resp.Error
}

// PregenerateByHash Enqueue a pregeneration request, by hash
// Account is optional, may by empty. 
// Default difficulty will be used
func PregenerateByHash(url string, hash string, account string) {
	addPregenerateRequest(WorkRequest{url, WorkInputHash, hash, 0, account})
}

// PregenerateByAccount Enqueue a pregeneration request, by account
// Default difficulty will be used
func PregenerateByAccount(url string, account string) {
	addPregenerateRequest(WorkRequest{url, WorkInputAccount, "", 0, account})
}

func waitForCacheResult(req WorkRequest) (WorkResponse, error) {
	// TODO do with events, timeout
	for i := 0; i < 120-5; i++ {
		found, _, resp := getWorkFromCache(req)
		if found {
			return resp, resp.Error
		}
		// not found, wait
		time.Sleep(500 * time.Millisecond)
	}
	// not found
	return WorkResponse{}, errors.New("Timeout in work generation")
}

// getWorkFromCache Try to get work from cache, nil is returned if not found in cache, or not valid
// Returns if valid work found in cache
// Returns if computation is in progress
func getWorkFromCache(req WorkRequest) (bool, bool, WorkResponse) {
	cachedEntry, ok := getFromCache(req.Hash)
	if (!ok) { return false, false, WorkResponse{} }
	// found in cache
	if !cacheIsValid(cachedEntry) {
		// found in cache, but not (yet) valid
		return false, true, WorkResponse{}
	}
	if !cacheDiffIsOK(cachedEntry, req.Diff) {
		// found but diff is smaller, must recompute
		log.Println("WARNING", "Found in cache, buf diff is smaller, hash", req.Hash, "cdiff", cachedEntry.difficulty, "diff", req.Diff)
		return false, false, WorkResponse{}
	}
	// found in cache, use it
	return true, false, WorkResponse {
		cachedEntry.hash,
		cachedEntry.work,
		cachedEntry.difficulty,
		cachedEntry.multiplier, 
		"cache",
		nil,
	}
}

// getCachedWork Retrieve work for a given hash; either from cache (if exists), or computed afresh from node.
// Account is optional (may be empty).
func getCachedWork(req WorkRequest) WorkResponse {
	statusWorkInReqCount++
	// Fill difficuly if missing
	if req.Diff == 0 {
		req.Diff = GetDefaultDifficulty()
	}
	// get from cache
	found, inprogress, respFromCache := getWorkFromCache(req)
	if (found) {
		// found in cache, use it
		statusWorkInReqFromCache++
		return respFromCache
	}
	if (inprogress) {
		// computation is in progress, wait
		log.Println("WARNING", "Work in progress but requested again, waiting; hash", req.Hash)
		// wait for result
		resp, err := waitForCacheResult(req)
		if err != nil {
			return WorkResponse{Error: err}
		}
		// do not count this in the cache stat (because it had to wait)
		return resp
	}
	// We need to call into RPC node for work.
	resp := getWorkFreshSync(req)
	return resp
}

// If input is account, get frontier first
func getCachedWorkByAccountOrHash(req WorkRequest) WorkResponse {
	if req.Input == WorkInputAccount {
		hash, err := GetFrontierHash(req.Url, req.Account)
		if err != nil { return WorkResponse{Error: err} }
		req.Hash = hash
	}
	resp := getCachedWork(req)
	return resp
}

var activeWorkOutReqCount int = 0
func decActiveWorkOutReqCount() { activeWorkOutReqCount-- }

// getWorkFreshSync Obtain the work now, by calling into the RPC node
// When result is obtained, it is added to cache.  Account is optional (may be empty).
func getWorkFreshSync(req WorkRequest) WorkResponse {
	activeWorkOutReqCount++
	defer decActiveWorkOutReqCount()

	statusWorkOutReqCount++
	if activeWorkOutReqCount >= maxOutRequests {
		// too many work requests
		return WorkResponse{Error: errors.New("Overload: too many active outgoing work requests")}
	}

	// mark start in cache
	addToCacheStart(req.Hash)
	log.Printf("Requesting work from node, reqCount %v  hash %v \n", activeWorkOutReqCount, req.Hash)
	// trigger work
	timeComputed := time.Now().Unix()
	resp, err, duration := rpcclient.GetWork(req.Url, req.Hash, req.Diff)
	if (err != nil) {
		return WorkResponse{Error: err}
	}
	
	// we have response, add to cache
	if (len(resp.Hash) == 0) { resp.Hash = req.Hash } // for the case if hash is missing in the response
	addToCache(resp, req.Account, timeComputed)
	statusWorkOutRespCount++
	log.Printf("Work resp from node, added to cache; dur %v, req %v, resp %v, \n", duration, req, resp)
	return WorkResponse {resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, "fresh", nil}
}

// get default difficulty -- TODO should come from RPC, cached
func GetDefaultDifficulty() uint64 {
	return 0xffffffc000000000;
}

func GetFrontierHash(url string, account string) (string, error) {
	// get frontier of account
	hash, err := rpcclient.GetFrontier(url, account)
	if (err != nil) {
		return "", errors.New("Could not obtain frontier block for account " + account + ", " + err.Error())
	}
	log.Println("Frontier block of account", account, "is", hash)
	return hash, nil
}

// StatusWorkOutReqCount Return the number of outgoing work requests (to node) since start (including currently pending ones)
func StatusWorkOutReqCount() int { return statusWorkOutReqCount }

// StatusWorkOutRespCount Return the number of outgoing work requests responses (from node) since start
func StatusWorkOutRespCount() int { return statusWorkOutRespCount }

// StatusWorkInReqCount Return the number of incoming work requests since start
func StatusWorkInReqCount() int { return statusWorkInReqCount }

// StatusWorkInReqFromCache Return the number of incoming work requests that could be serviced from the cache
func StatusWorkInReqFromCache() int { return statusWorkInReqFromCache }

func StatusActiveWorkOutReqCount() int { return activeWorkOutReqCount }
