// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"math/rand"
)

var testHashes []string
var hexDigits string = "0123456789ABCDEF"

func randomHash() string {
	index := rand.Intn(len(testHashes))
	return testHashes[index]
}

func randomHexDigit() string {
	index := rand.Intn(16)
	return hexDigits[index:index+1]
}

func generateRandomHash() string {
	var hash string = ""
	for i := 0; i < 64; i++ {
		hash += randomHexDigit()
	}
	return hash
}

func initTestData(hashCount int) {
	testHashes = make([]string, hashCount)
	for i := 0; i < hashCount; i++ {
		testHashes[i] = generateRandomHash()
	}
	fmt.Printf("Test data: %v hashes generated\n", len(testHashes))
}
