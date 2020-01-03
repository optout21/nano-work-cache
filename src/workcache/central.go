// Copyright © 2019-2019 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"log"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

type WorkResponse struct {
	Hash string
	Work string
	Difficulty uint64
	Multiplier float64
	// values: 'fresh', 'fromcache'
	Source string
}

func GetCachedWork(url string, hash string, diff uint64) (WorkResponse, error) {
	cachedEntry, ok := getFromCache(hash)
	if (ok) {
		if cacheIsValid(cachedEntry) {
			if cacheDiffIsOK(cachedEntry, diff) {
				// found in cache, use it
				return WorkResponse {
					cachedEntry.hash,
					cachedEntry.work,
					cachedEntry.difficulty,
					cachedEntry.multiplier, 
					"fromcache",
				}, nil
			} else {
				// found but diff is smaller, must recompute
				log.Println("WARNING", "Found in cache, buf diff is smaller, hash", hash, "cdiff", cachedEntry.difficulty, "diff", diff)
			}
		} else {
			// found in cache, but not (yet) valid
			// TODO we could wait here, to avoid starting again
			log.Println("WARNING", "Work in progress, yet starting again, hash", hash)
		}
	}
	// We need to call into RPC node for work.
	resp, err := callRpcWork(url, hash, diff)
	return resp, err
}

/// Request work from remote RPC node
func callRpcWork(url string, hash string, diff uint64) (WorkResponse, error) {
	// mark start in cache
	addToCacheStart(hash)
	// trigger work
	resp, err := rpcclient.GetWork(url, hash, diff)
	if (err != nil) {
		return WorkResponse{}, err
	}
	// we have response, add to cache
	if (len(resp.Hash) == 0) { resp.Hash = hash } // for the case if hash is missing in the response
	addToCache(resp)
	log.Println("Work resp from node, added to cache; work_generate resp", resp)
	return WorkResponse {resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, "fresh"}, nil
}

// get default difficulty -- TODO should come from RPC, cached
func GetDefaultDifficulty() uint64 {
	return 0xffffffc000000000;
}
