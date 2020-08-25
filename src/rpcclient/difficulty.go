// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package rpcclient

import (
	"strconv"
	"time"
)

var difficulty uint64 = 0xffffffc000000000
var diffTime time.Time = time.Now().Add(-100 * time.Hour)

const cacheExpiry time.Duration = 10 * time.Minute

// GetDifficultyCached Get the current network difficulty, comes from RPC, cached for some minutes
func GetDifficultyCached(url string) uint64 {
	now := time.Now()
	age := now.Sub(diffTime)
	//fmt.Printf("diff %v age %v \n", difficulty, age)
	if difficulty > 0 && age <= cacheExpiry {
		// valid and fresh, return cached
		return difficulty
	}
	diffRpc, err := GetDifficulty(url)
	if err != nil {
		return difficulty
	}
	difficultyParsed, err := strconv.ParseUint(diffRpc, 16, 64)
	if err != nil {
		return difficulty
	}
	// store it
	difficulty = difficultyParsed
	diffTime = now
	return difficulty
}
