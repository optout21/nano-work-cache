// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"fmt"
	"log"
	"net/http"
)

var activeHandlerCount int = 0
var maxActiveHandlerCount int = 200

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

func handleReqWithRateLimit(action string, respBody []byte, w http.ResponseWriter) {
	if activeHandlerCount >= maxActiveHandlerCount {
		// overload, return error right away
		log.Printf("Overload, %v active request handlers, max %v\n", activeHandlerCount, maxActiveHandlerCount)
		fmt.Fprintln(w, fmt.Sprintf(`{"error":"overload, too many concurrent active requests"}`))
		return
	}
	incActiveCount()
	defer decActiveCount()
	handleReqSync(action, respBody, w)
}
