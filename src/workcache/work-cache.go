// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

type CacheEntry struct {
	hash string
	work string
	difficulty uint64
	multiplier float64
	account string
	// valid, computing
	status string
	// time started
	// time added
}

var workCache map[string]CacheEntry = map[string]CacheEntry{}

// Add a work result to the cache
func addToCache(e rpcclient.WorkResponse) {
	addToCacheInternal(CacheEntry{
		e.Hash,
		e.Work,
		e.Difficulty,
		e.Multiplier,
		"",
		"valid",
	})
}

// Mark in the cache that work request has started
func addToCacheStart(hash string) {
	addToCacheInternal(CacheEntry{
		hash,
		"",
		0,
		0,
		"",
		"computing",
	})
}
func addToCacheInternal(e CacheEntry) {
	workCache[e.hash] = e
}

func getFromCache(hash string) (CacheEntry, bool) {
	e, ok := workCache[hash]
	if (!ok) {
		// not in cache
		return e, false
	}
	// found in cache
	return e, true
}

func cacheIsValid(e CacheEntry) bool {
	if (e.status == "valid") {
		return true
	}
	return false
}

// Note: difficulty may be missing (0)
func cacheDiffIsOK(e CacheEntry, diff uint64) bool {
	//fmt.Printf("diff %x %x\n", e.difficulty, diff)
	if (diff != 0 && e.difficulty != 0 && e.difficulty < diff) {
		// but diff is smaller
		return false
	}
	// diff is OK (larger or equal)
	return true
}

// StatusCacheSize Return the current number of entries in the cache
func StatusCacheSize() int {
	return len(workCache)
}
