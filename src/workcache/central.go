// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

const (
	WorkInputHash    = 0
	WorkInputAccount = 1
)

type WorkRequest struct {
	Input   int // WorkInputHash or Account
	Hash    string
	Diff    uint64
	Account string
}

type WorkResponse struct {
	Hash       string
	Work       string
	Difficulty uint64
	Multiplier float64
	// values: 'fresh', 'cache'
	Source string
	Error  error
}

var maxOutRequests int = 0
var maxCacheAgeDays int = 0
var statusWorkOutReqCount int = 0
var statusWorkOutRespCount int = 0
var statusWorkOutDurationTotal int64 = 0
var statusWorkInReqCount int = 0
var statusWorkInReqFromCache int = 0
var statusWorkInReqError int = 0

// Start Invoked at the beginning, can perform initializations, read the cache, etc.
func Start() {
	backgroundWorkerCount := ConfigBackgroundWorkerCount()
	maxOutRequests = ConfigMaxOutRequests()
	maxCacheAgeDays = ConfigMaxCacheAgeDays()
	InitQueue()
	LoadCache()
	RemoveOldEntries(float64(maxCacheAgeDays))
	startWorkers(backgroundWorkerCount)
	go housekeepingCycle()
}

// Generate Generate work or take from cache. Generation done in foreground.
// Account is optional, may by empty.
// Difficulty may be 0, default will be used
func Generate(hash string, difficulty uint64, account string) (WorkResponse, error) {
	req := WorkRequest{WorkInputHash, hash, difficulty, account}
	resp, fromcache := getCachedWork(req)
	statusWorkInReqCount++
	if fromcache {
		statusWorkInReqFromCache++
	}
	if resp.Error != nil || !IsWorkValueValid(resp.Work) {
		statusWorkInReqError++
	}
	return resp, resp.Error
}

// PregenerateByHash Enqueue a pregeneration request, by hash
// Account is optional, may by empty.
// Default difficulty will be used
func PregenerateByHash(hash string, account string) {
	req := WorkRequest{WorkInputHash, hash, 0, account}
	// check in cache
	found, _, _ := getWorkFromCache(req)
	if found {
		// found in cache, no need to compute
		return
	}
	addPregenerateRequest(req)
}

// PregenerateByAccount Enqueue a pregeneration request, by account
// Default difficulty will be used
func PregenerateByAccount(account string) {
	req := WorkRequest{WorkInputAccount, "", 0, account}
	// check if frontier hash has work in cache
	// get frontier of account
	hash, err := GetFrontierHash(account)
	if err != nil {
		// could not get frontier, add it as fallback
		addPregenerateRequest(req)
		return
	}
	// check in cache
	found, _, _ := getWorkFromCache(WorkRequest{WorkInputHash, hash, 0, account})
	if found {
		// found in cache, no need to compute
		return
	}
	// not found, add pregenerate request, but by account
	addPregenerateRequest(req)
}

func waitForCacheResult(req WorkRequest) (WorkResponse, error) {
	// TODO do with events, timeout
	for i := 0; i < 100-1; i++ {
		found, _, resp := getWorkFromCache(req)
		if found {
			return resp, resp.Error
		}
		// not found, wait
		time.Sleep(250 * time.Millisecond)
	}
	// not found
	return WorkResponse{}, errors.New("Timeout in work generation")
}

// Check is a work value string looks valid: not empty, hex string
func IsWorkValueValid(work string) bool {
	return len(work) > 9
}

// getWorkFromCache Try to get work from cache, nil is returned if not found in cache, or not valid
// Returns if valid work found in cache
// Returns if computation is in progress
func getWorkFromCache(req WorkRequest) (bool, bool, WorkResponse) {
	cachedEntry, ok := getFromCache(req.Hash)
	if !ok {
		return false, false, WorkResponse{}
	}
	// found in cache
	if !cacheIsValid(cachedEntry) {
		// found in cache, but not (yet) valid
		return false, true, WorkResponse{}
	}
	if !IsWorkValueValid(cachedEntry.work) {
		log.Println("WARNING", "Invalid work value found in cache", cachedEntry.work, cachedEntry.hash)
		return false, false, WorkResponse{}
	}
	if !cacheDiffIsOK(cachedEntry, req.Diff) {
		// found but diff is smaller, must recompute
		log.Println("WARNING", "Found in cache, buf diff is smaller, hash", req.Hash, "cdiff", cachedEntry.difficulty, "diff", req.Diff)
		return false, false, WorkResponse{}
	}
	// found in cache, use it
	return true, false, WorkResponse{
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
// Return response and true if it is taken from cache
func getCachedWork(req WorkRequest) (WorkResponse, bool) {
	// Fill difficuly if missing
	if req.Diff == 0 {
		req.Diff = rpcclient.GetDifficultyCached()
	}
	// get from cache
	found, inprogress, respFromCache := getWorkFromCache(req)
	if found {
		// found in cache, use it
		return respFromCache, true
	}
	if inprogress {
		// computation is in progress, wait
		log.Println("WARNING", "Work in progress but requested again, waiting; hash", req.Hash)
		// wait for result
		resp, err := waitForCacheResult(req)
		if err != nil {
			// non-success (timeout), do not count as cache success
			return WorkResponse{Error: err}, false
		}
		return resp, true
	}
	// We need to call into RPC node for work.
	resp := getWorkFreshSync(req)
	return resp, false
}

// If input is account, get frontier first
func getCachedWorkByAccountOrHash(req WorkRequest) WorkResponse {
	if req.Input == WorkInputAccount {
		hash, err := GetFrontierHash(req.Account)
		if err != nil {
			return WorkResponse{Error: err}
		}
		req.Hash = hash
	}
	resp, _ := getCachedWork(req)
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
	if maxOutRequests > 0 {
		if activeWorkOutReqCount >= maxOutRequests {
			// too many work requests
			return WorkResponse{Error: fmt.Errorf("Overload: too many active outgoing work requests %v %v", activeWorkOutReqCount, maxOutRequests)}
		}
	}

	// mark start in cache
	addToCacheStart(req.Hash)
	log.Printf("Requesting work from node, reqCount %v  hash %v \n", activeWorkOutReqCount, req.Hash)
	// trigger work
	timeComputed := time.Now().Unix()
	resp, err, duration := rpcclient.GetWork(req.Hash, req.Diff)
	if err != nil {
		return WorkResponse{Error: err}
	}

	// we have response, add to cache
	if len(resp.Hash) == 0 {
		resp.Hash = req.Hash
	} // for the case if hash is missing in the response
	addToCache(resp, req.Account, timeComputed)
	statusWorkOutRespCount++
	statusWorkOutDurationTotal += duration.Milliseconds()
	log.Printf("Work resp from node, added to cache; dur %v, req %v, resp %v, \n", duration, req, resp)
	return WorkResponse{resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, "fresh", nil}
}

func GetFrontierHash(account string) (string, error) {
	// get frontier of account
	hash, err := rpcclient.GetFrontier(account)
	if err != nil {
		return "", errors.New("Could not obtain frontier block for account " + account + ", " + err.Error())
	}
	log.Println("Frontier block of account", account, "is", hash)
	return hash, nil
}

// StatusWorkOutReqCount Return the number of outgoing work requests (to node) since start (including currently pending ones)
func StatusWorkOutReqCount() int { return statusWorkOutReqCount }

// StatusWorkOutRespCount Return the number of outgoing work requests responses (from node) since start
func StatusWorkOutRespCount() int { return statusWorkOutRespCount }

// statusWorkOutDurationAvg Return the average duration in ms of the outgoing work requests
func StatusWorkOutDurationAvg() int {
	if statusWorkOutRespCount == 0 {
		return 0
	}
	return int(float32(statusWorkOutDurationTotal) / float32(statusWorkOutRespCount))
}

// StatusWorkInReqCount Return the number of incoming work requests since start
func StatusWorkInReqCount() int { return statusWorkInReqCount }

// StatusWorkInReqFromCache Return the number of incoming work requests that could be serviced from the cache
func StatusWorkInReqFromCache() int { return statusWorkInReqFromCache }

// StatusWorkInReqError Return the number of incoming work requests that were returned with error
func StatusWorkInReqError() int { return statusWorkInReqError }

func StatusActiveWorkOutReqCount() int { return activeWorkOutReqCount }
