// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func isPersistToFileEnabled() bool {
	filename := ConfigGetString("Main.CachePeristFileName")
	if len(filename) == 0 {
		return false
	}
	return true
}

func persistFileName() string {
	filename := ConfigGetString("Main.CachePeristFileName")
	return filename
}

// SaveCache Save the cache to file or other persistence configured
func SaveCache() {
	if !isPersistToFileEnabled() {
		return
	}
	saveToFile(persistFileName())
}

// LoadCache Load the cache from file or other persistence configured
func LoadCache() {
	if !isPersistToFileEnabled() {
		return
	}
	loadFromFile(persistFileName())
}

// saveToFile save cache to the given file
func saveToFile(filename string) {
	// first try to rename old file to .bak (ignore error if not possible / not exists)
	fileNameBak := filename + ".bak"
	_ = os.Remove(fileNameBak)
	_ = os.Rename(filename, fileNameBak)
	// open new file for writing
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("Error saving cache to file, could not open file;", err.Error())
		return
	}
	defer file.Close()

	var cnt int = 0
	for _, entry := range workCache {
		line := entryToString(entry)
		if len(line) > 0 {
			fmt.Fprintln(file, entryToString(entry))
			cnt++
		}
	}
	log.Printf("Cache saved to file, %v entries\n", cnt)
}

// loadFromFile Read cache entries from the given file, merge them with current cache
func loadFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error loading cache from file, could not open file;", err.Error())
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
		if !lineParsed {
			continue
		}
		// omit non-"valid" entries
		if entry.status != "valid" {
			continue
		}
		addToCacheInternal(entry)
		cnt++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Cache loaded from file %v, %v entries read, %v entries\n", filename, cnt, StatusCacheSize())
}
