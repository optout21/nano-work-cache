package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
    "log"
	"net/http"
	"strconv"
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

/// Not the normal Json Encode way, due to the difficult hex formatting.  Using simple string concatenation.
func workResponseToJson(resp workcache.WorkResponse) string {
	return fmt.Sprintf(`{"hash":"%v","work":"%v","difficulty":"%x","multiplier":"%v","source":"%v"}`,
		resp.Hash, resp.Work, resp.Difficulty, resp.Multiplier, resp.Source)
}

func handleJson(action string, respBody []byte, w http.ResponseWriter) {
	switch action {
	case "work_generate":
		var workGenerate workGenerateJson
		err := json.Unmarshal(respBody, &workGenerate)
		if err != nil {
			fmt.Fprintln(w, `{"error":"work_generate parse error"}`)
			return;
		}
		log.Println("work_generate req", workGenerate)
		var difficulty uint64 = workcache.GetDefaultDifficulty()
		if (len(workGenerate.Difficulty) > 0) {
			difficultyParsed, err := strconv.ParseUint(workGenerate.Difficulty, 16, 64)
			if (err != nil) {
				// diff present, but could not parse
				fmt.Fprintln(w, `{"error":"work_generate difficulty parse error"}`)
				return;
			}
			difficulty = difficultyParsed
		}
		// handle
		workResp, err := workcache.GetCachedWork(nanoNodeUrl, workGenerate.Hash, difficulty)
		if (err != nil) {
			fmt.Fprintln(w, `{"error":"` + err.Error() + `"}`)
			return
		}
		log.Println("work_generate resp", workResp)
		fmt.Fprintln(w, workResponseToJson(workResp))

	default:
		fmt.Fprintln(w, `{"error":"unknown action","action":"` + action + "}")
	}
}

var nanoNodeUrl string

func Start(nanoNodeUrl1 string) {
	nanoNodeUrl = nanoNodeUrl1
    http.HandleFunc("/", handleRequest)

    http.ListenAndServe(":8080", nil)
}
