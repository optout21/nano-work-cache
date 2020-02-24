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
				log.Println("action", action)
				handleJson(action.Action, body, w)
			}
		}
		/*
        // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
        if err := r.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
        fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
        name := r.FormValue("name")
        address := r.FormValue("address")
        fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
		*/
    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }
}

type workGenerateJson struct {
	Action string
	Hash string
	Difficulty string
}

type workPregenerateByHashJson struct {
	Action string
	Hash string
}

type workPregenerateByAccountJson struct {
	Action string
	Account string
}

type accountBalanceJson struct {
	Action string
	Account string
}

type accountsBalancesJson struct {
	Action string
	Accounts []string
}

/// Not the normal Json Encode way, due to the difficult hex formatting.  Using simple string concatenation.
func workResponseToJson(resp workcache.WorkResponse) string {
	return fmt.Sprintf(`{"hash":"%v","work":"%v","difficulty":"%x","multiplier":"%v","source":"%v"}`,
		resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, resp.Source)
}

var nanoNodeUrl string

func Start(nanoNodeUrl1 string, listenIpPort string) {
	nanoNodeUrl = nanoNodeUrl1
    http.HandleFunc("/", handleRequest)

	log.Println("Starting listening on", listenIpPort, "...")
    http.ListenAndServe(listenIpPort, nil)
}
