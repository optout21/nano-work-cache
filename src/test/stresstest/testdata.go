// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package main

import (
	"fmt"
	"math/rand"
)

// Random hash test data
var tdRandomHashes []string
// Some valid accounts
var tdValidAccounts [5]string = [5]string{
	"nano_11kfwesa9x3gsrmnjbujmrdrbc4s8sunbdfd5ttnertwjh9xm8uqs39tbsjy",
	"nano_1ms3w8m19au6ig8hw6zwygjhpnttrn73r4kisw33pj48edwearer6bd7rfei",
	"nano_1pndfqtkxb1ob1dgeszus9jcswj4xzt8z5q6ueopdrsbphrpnjsaas9hk86e",
	"nano_3jwrszth46rk1mu7rmb4rhm54us8yg1gw3ipodftqtikf5yqdyr7471nsg1k",
	"nano_3rpb7ddcd6kux978gkwxh1i1s6cyn7pw3mzdb9aq7jbtsdfzceqdt3jureju",
}

// Chose a random hash at random
func getRandomHash() string {
	index := rand.Intn(len(tdRandomHashes))
	return tdRandomHashes[index]
}

// Chose a valid account at random
func getValidAccount() string {
	index := rand.Intn(len(tdValidAccounts))
	return tdValidAccounts[index]
}

var hexDigits string = "0123456789ABCDEF"

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
	tdRandomHashes = make([]string, hashCount)
	for i := 0; i < hashCount; i++ {
		tdRandomHashes[i] = generateRandomHash()
	}
	fmt.Printf("Test data: %v random hashes generated\n", len(tdRandomHashes))
}
