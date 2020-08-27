// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"fmt"
	"log"
	"net/http"
)

var activeHandlerCount int = 0
var maxActiveRequests int = 500

// ActiveHandlerCount Return current number of conccurrent active handlers
func ActiveHandlerCount() int { return activeHandlerCount }

func incActiveCount() int {
	activeHandlerCount++
	return activeHandlerCount
}
func decActiveCount() int {
	activeHandlerCount--
	return activeHandlerCount
}

// Handle incoming calls with rate limiting; if max is reached Overload error is returned
func handleReqWithRateLimit(action string, respBody []byte, w http.ResponseWriter) {
	defer decActiveCount()
	incActiveCount()
	if activeHandlerCount >= maxActiveRequests {
		// overload, return error right away
		log.Printf("Overload, %v active request handlers, max %v\n", activeHandlerCount, maxActiveRequests)
		fmt.Fprintln(w, fmt.Sprintf(`{"error":"overload, too many concurrent active requests"}`))
		return
	}

	handleReqSync(action, respBody, w)
}
