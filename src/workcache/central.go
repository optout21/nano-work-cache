// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"errors"
	"log"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

type WorkResponse struct {
	Hash string
	Work string
	Difficulty uint64
	Multiplier float64
	// values: 'fresh', 'cache'
	Source string
	Error error
}

var statusWorkReqCount int = 0
var statusWorkRespCount int = 0

// Start Invoked at the beginning, can perform initializations, read the cache, etc.
func Start() {
	LoadCache()
}

// GetWorkFromCache Try to get work from cache, nil is returned if not found in cache, or not valid
func GetWorkFromCache(url string, hash string, diff uint64) (bool, WorkResponse) {
	cachedEntry, ok := getFromCache(hash)
	if (!ok) { return false, WorkResponse{} }
	// found in cache
	if !cacheIsValid(cachedEntry) {
		// found in cache, but not (yet) valid
		// TODO we could wait here, to avoid starting again
		log.Println("WARNING", "Work in progress, yet starting again, hash", hash)
		return false, WorkResponse{}
	}
	if !cacheDiffIsOK(cachedEntry, diff) {
		// found but diff is smaller, must recompute
		log.Println("WARNING", "Found in cache, buf diff is smaller, hash", hash, "cdiff", cachedEntry.difficulty, "diff", diff)
		return false, WorkResponse{}
	}
	// found in cache, use it
	return true, WorkResponse {
		cachedEntry.hash,
		cachedEntry.work,
		cachedEntry.difficulty,
		cachedEntry.multiplier, 
		"cache",
		nil,
	}
}

// GetCachedWork Retrieve work for a given hash; either from cache (if exists), or computed afresh from node.
// Account is optional (may be empty).
func GetCachedWork(url string, hash string, diff uint64, account string) (WorkResponse, error) {
	// get from cache
	found, respFromCache := GetWorkFromCache(url, hash, diff)
	if (found) {
		// found in cache, use it
		return respFromCache, nil
	}
	// We need to call into RPC node for work.
	resp := getWorkSync(url, hash, diff, account)
	return resp, resp.Error
}

// getWorkSync Obtain the work now, by calling into the RPC node
// When result is obtained, it is added to cache.  Account is optional (may be empty).
func getWorkSync(url string, hash string, diff uint64, account string) WorkResponse {
	statusWorkReqCount++
	// mark start in cache
	addToCacheStart(hash)
	log.Println("Requesting work from node for hash", hash)
	// trigger work
	resp, err := rpcclient.GetWork(url, hash, diff)
	if (err != nil) {
		return WorkResponse{Error: err}
	}
	// we have response, add to cache
	if (len(resp.Hash) == 0) { resp.Hash = hash } // for the case if hash is missing in the response
	addToCache(resp, account)
	statusWorkRespCount++
	go SaveCache()
	log.Println("Work resp from node, added to cache; work_generate resp", resp)
	return WorkResponse {resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, "fresh", nil}
}

// get default difficulty -- TODO should come from RPC, cached
func GetDefaultDifficulty() uint64 {
	return 0xffffffc000000000;
}

/// First obtain frontier hash of the account, then request work for the hash (if needed), by calling GetCachedWork
func GetCachedWorkByAccount(url string, account string) WorkResponse {
	// get frontier of account
	hash, err := rpcclient.GetFrontier(url, account)
	if (err != nil) {
		return WorkResponse{Error: errors.New("Could not obtain frontier block for account " + account + ", " + err.Error())}
	}
	log.Println("Frontier block of account", account, "is", hash)
	difficulty := GetDefaultDifficulty()
	resp, _ := GetCachedWork(url, hash, difficulty, account)
	return resp
}


// StatusWorkReqCount Return the number of work requests (to node) since start (including currently pending ones)
func StatusWorkReqCount() int { return statusWorkReqCount }

// StatusWorkRespCount Return the number of work requests responses (from node) since start
func StatusWorkRespCount() int { return statusWorkRespCount }
