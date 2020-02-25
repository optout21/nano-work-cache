// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"fmt"
	"strconv"
	"strings"
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
	timeStarted uint64
	timeAdded uint64
}

var workCache map[string]CacheEntry = map[string]CacheEntry{}

// Add a work result to the cache.  Account is optional (may be empty).
func addToCache(e rpcclient.WorkResponse, account string) {
	addToCacheInternal(CacheEntry{
		e.Hash,
		e.Work,
		e.Difficulty,
		e.Multiplier,
		account,
		"valid",
		0,
		0,
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
		0,
		0,
	})
}

func addToCacheInternal(e CacheEntry) {
	if len(e.hash) > 0 {
		workCache[e.hash] = e
	}
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

func padString(val string) string {
	if len(val) == 0 { return "_" }
	return val
}

// Convert an entry to a single-line stirng representation
func entryToString(entry CacheEntry) string {
	if len(entry.hash) == 0 { return "" }
	return fmt.Sprintf("%v %v %x %v %v %v %v %v", padString(entry.hash), padString(entry.work), entry.difficulty, entry.multiplier, 
		padString(entry.account), padString(entry.status), entry.timeStarted, entry.timeAdded)
}

// Fill cache entry from a single-line stirng represenation (parse it), see entryToString.
// Returns true on success. 
func entryLoadFromString(line string, entry *CacheEntry) bool {
	tokens := strings.Split(line, " ")
	if len(tokens) < 2 {
		// mimium hash and work values are needed; this is too short
		return false
	}
	entry.hash = tokens[0]
	entry.work = tokens[1]
	if len(tokens) >= 8 {
		diff, _ := strconv.ParseUint(tokens[2], 16, 64)
		entry.difficulty = diff
		multip, _ := strconv.ParseFloat(tokens[3], 64)
		entry.multiplier = multip
		entry.account = tokens[4]
		entry.status = tokens[5]
		timeStart, _ := strconv.ParseUint(tokens[6], 10, 64)
		timeAdded, _ := strconv.ParseUint(tokens[7], 10, 64)
		entry.timeStarted = timeStart
		entry.timeAdded = timeAdded
	}
	return true
}
