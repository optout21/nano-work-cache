// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
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

func backupFileName(filename string) string {
	return filename + ".bak"
}

// LoadCache Load the cache from file or other persistence configured
func LoadCache() {
	if !isPersistToFileEnabled() {
		return
	}
	filename := persistFileName()
	err := loadFromFile(filename)
	if err == nil {
		return
	}
	// try bak file
	_ = loadFromFile(backupFileName(filename))
}

// saveToFile save cache to the given file
func saveToFile(filename string) {
	startTime := time.Now()

	// try to rename old file to .bak (ignore error if not possible / not exists)
	fileNameBak := backupFileName(filename)
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
	elapsed2 := time.Now().Sub(startTime)
	log.Printf("Cache saved to file, %v entries, dur %v ms", cnt, elapsed2.Milliseconds())
}

// loadFromFile Read cache entries from the given file, merge them with current cache
func loadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error loading cache from file, could not open file;", err.Error())
		return err
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
		return err
	}

	log.Printf("Cache loaded from file %v, %v entries read, %v entries\n", filename, cnt, StatusCacheSize())
	return nil
}
