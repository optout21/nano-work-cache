// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// SaveCache Save the cache to file or orther persistence configured
func SaveCache() {
	filename := ConfigGetString("Main.CachePeristFileName")
	if len(filename) == 0 { return }
	saveToFile(filename)
}

// LoadCache Load the cache from file or orther persistence configured
func LoadCache() {
	filename := ConfigGetString("Main.CachePeristFileName")
	if len(filename) == 0 { return }
	loadFromFile(filename)
}

// saveToFile save cache to the given file
func saveToFile(filename string) {
	// first try to rename old file (ignore error if not possible / not exists)
	fileNameBak := filename + ".bak"
	_ = os.Remove(fileNameBak)
	_ = os.Rename(filename, fileNameBak)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error saving cache to file, could not open file", err.Error())
		return
	}
	defer file.Close()

	var cnt int = 0
	for _, entry := range workCache {
		line := entryToString(entry);
		if len(line) > 0 {
			fmt.Fprintln(file, entryToString(entry));
			cnt++
		}
	}
	log.Printf("Cache saved to file, %v entries\n", cnt)
}

// loadFromFile Read cache entries from the given file, merge them with current cache
func loadFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error loading cache from file, could not open file", err.Error())
		return
	}
	defer file.Close()

	var cnt int = 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
		// one entry is one line
		line := scanner.Text()
		var entry CacheEntry
		lineParsed := entryLoadFromString(line, &entry)
		if !lineParsed { continue; }
		addToCacheInternal(entry)
		cnt++
	}

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}

	log.Printf("Cache loaded from file, %v entries read, %v entries\n", cnt, StatusCacheSize())
}
