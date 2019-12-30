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
			log.Println("WARNING", "Work in progress, yet staring again, hash", hash)
		}
	}
	// not found in cache, or being computed
	// mask start in cache
	addToCacheStart(hash)
	// trigger work
	resp, err := rpcclient.GetWork(url, hash, diff)
	if (err != nil) {
		return WorkResponse{}, err
	}
	// we have response, add to cache
	addToCache(resp)
	return WorkResponse {resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, 
		"fresh"}, nil
}
