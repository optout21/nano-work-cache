// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"log"
	"time"
)

var housekeepingPeriodSec int = 1 * 60

var lastCacheSaveTime int64 = 0

// Housekeeping is executed periodically.  It incudes:
// - Saving the cachefile (if it has changed since last time)
func housekeepingCycle() {
	lastCacheSaveTime = CacheUpdateTime()
	for {
		doHousekeepingCycle()
		time.Sleep(time.Second * time.Duration(housekeepingPeriodSec))
	}
}

func doHousekeepingCycle() {
	if isPersistToFileEnabled() {
		cacheUpdateTime := CacheUpdateTime()
		origLastCacheSaveTime := lastCacheSaveTime
		if cacheUpdateTime > lastCacheSaveTime {
			SaveCache()
			lastCacheSaveTime = cacheUpdateTime
			log.Printf("Cache saved, %v %v \n", origLastCacheSaveTime, cacheUpdateTime)
		}
	}
}
