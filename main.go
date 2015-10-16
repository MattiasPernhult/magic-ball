package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Choice struct which handle response
type Choice struct {
	Answer string `json:"answer"`
	Forced bool   `json:"forced"`
}

// Error struct
type Error struct {
	Message string `json:"message"`
}

func main() {
	r := httprouter.New()
	r.NotFound = http.FileServer(http.Dir("public"))
	r.GET("/generate", Generate)
	r.GET("/generate/:force", Generate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, r)
	http.ListenAndServe(":"+port, nil)
}

// Generate method is handling the base request
func Generate(rw http.ResponseWriter,
	r *http.Request, p httprouter.Params) {
	randomNr := rand.Intn(10)

	positive := []string{"My sources says yes", "Of course", "You should definitely do that"}
	negative := []string{"My sources says no", "I don't think so", "You should definitely not do that"}

	var ch Choice

	err, forced, positiveForced := isForced(rw, p)

	if !err {
		if forced {
			if positiveForced {
				index := rand.Intn(len(positive))
				ch = Choice{positive[index], true}
			} else {
				index := rand.Intn(len(positive))
				ch = Choice{negative[index], true}
			}
		} else {
			if randomNr%2 == 0 {
				index := rand.Intn(len(positive))
				ch = Choice{positive[index], false}
			} else {
				index := rand.Intn(len(positive))
				ch = Choice{negative[index], false}
			}
		}

		writeResponse(rw, ch)
	}
}

func isForced(rw http.ResponseWriter, p httprouter.Params) (bool, bool, bool) {
	var forced bool
	var positive bool
	var errorr bool
	force := p.ByName("force")

	if force == "" {
		return false, false, false
	}

	if force == "positive" {
		forced = true
		positive = true
		errorr = false
	} else if force == "negative" {
		forced = true
		positive = false
		errorr = false
	} else {
		err1 := Error{"Must be either positive or negative"}
		js, _ := json.Marshal(err1)
		http.Error(rw, string(js), http.StatusBadRequest)
		errorr = true
	}
	return errorr, forced, positive
}

func writeResponse(rw http.ResponseWriter, ch Choice) {
	js, err := json.Marshal(ch)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
