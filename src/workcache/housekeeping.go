// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"log"
	"time"
)

var housekeepingPeriodSec int = 1 * 60

var lastCacheSaveTime int64 = 0
var lastAgeCheckTime int64 = 0

// Housekeeping is executed periodically.  It incudes:
// - Saving the cachefile (if it has changed since last time)
func housekeepingCycle() {
	lastCacheSaveTime = CacheUpdateTime()
	lastAgeCheckTime = CacheUpdateTime()
	for {
		doHousekeepingCycle()
		time.Sleep(time.Second * time.Duration(housekeepingPeriodSec))
	}
}

func doHousekeepingCycle() {
	cacheUpdateTime := CacheUpdateTime()

	if maxCacheAgeDays > 0 {
		if cacheUpdateTime > lastAgeCheckTime {
			log.Printf("Running cache aging (%v %v)\n", lastAgeCheckTime, cacheUpdateTime)
			RemoveOldEntries(float64(maxCacheAgeDays))
			lastAgeCheckTime = cacheUpdateTime
		}
	}

	if isPersistToFileEnabled() {
		origLastCacheSaveTime := lastCacheSaveTime
		if cacheUpdateTime > lastCacheSaveTime {
			SaveCache()
			lastCacheSaveTime = cacheUpdateTime
			log.Printf("Cache saved, %v %v \n", origLastCacheSaveTime, cacheUpdateTime)
		}
	}
}
