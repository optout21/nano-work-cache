// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"log"
	"time"
)

// Background generate jobs, with low priority.  Size is large, fixed.
var pregenerateJobsMaxSize int = 10000
var pregenerateJobs chan WorkRequest = make(chan WorkRequest, pregenerateJobsMaxSize)

func addPregenerateRequest(req WorkRequest) {
	// check in cache
	found, _, _ := getWorkFromCache(req)
	if (found) {
		// found in cache, no need to compute
		return
	}
	if len(pregenerateJobs) >= pregenerateJobsMaxSize - 2 {
		// queue is full, do not put any more (to avoid blocking)
		log.Printf("WARNING: Pregeneration queue is full, not enqueuing any more, %v\n", len(pregenerateJobs))
	}
	pregenerateJobs <- req
}

func doProcess(name int) {	
	for {
		// wait on queue, with periodical timeout
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		select {
		case preJob := <- pregenerateJobs:
			//log.Printf("Worker %v : pregenerate job", name)
			getCachedWorkByAccountOrHash(preJob)

		case <-ticker.C:
			// timeout, idle loop
		}
	}
}

func startWorkers(backgroundWorkerCount int) {
	for i := 0; i < backgroundWorkerCount; i++ {
		go doProcess(i)
	}
	log.Printf("%v pool workers started\n", backgroundWorkerCount)
}

func StatusPregenerQueueSize() int { return len(pregenerateJobs) }
