// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/catenocrypt/nano-work-cache/workcache"
)

type actionJson struct {
	Action string
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch req.Method {
	case "GET":
		fmt.Fprintf(w, "Welcome to my website!")
		//http.ServeFile(w, r, "form.html")

	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintln(w, `{"error":"could not read post data"}`)
		} else {
			//decoder := json.NewDecoder(body)
			var action actionJson
			//err := decoder.Decode(&action)
			err := json.Unmarshal(body, &action)
			if err != nil {
				fmt.Fprintln(w, `{"error":"action parse error"}`)
			} else {
				//userAgent := req.Header["User-Agent"][0]
				//log.Println("userAgent", userAgent, "remoteAddr", req.RemoteAddr, "action", action.Action)
				handleReqWithRateLimit(action.Action, body, w)
			}
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func Start() {
	listenIpPort := workcache.ConfigListenIpPort()
	maxActiveRequests = workcache.ConfigRestMaxActiveRequests()
	enablePregeneration = workcache.ConfigEnablePregeneration()

	http.HandleFunc("/", handleRequest)

	log.Println("Starting listening on", listenIpPort, "...")
	http.ListenAndServe(listenIpPort, nil)
}
