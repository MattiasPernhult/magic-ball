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
	Type   string `json:"type"`
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
}

// Generate method is handling the base request
func Generate(rw http.ResponseWriter,
	r *http.Request, p httprouter.Params) {
	randomNr := rand.Intn(10)

	positive := []string{"It is certain", "It is decidedly so", "Without a doubt",
		"Yes, definitely", "You may rely on it", "As I see it, yes", "Most likely",
		"Outlook good", "Yes", "Signs point to yes"}
	negative := []string{"Don't count on it", "My reply is no",
		"My sources say no", "Outlook not so good", "Very doubtful"}
	neutral := []string{"Reply hazy try again", "Ask again later",
		"Better not tell you now", "Cannot predict now", "Concentrate and ask again"}

	var ch Choice

	err, forced, value := isForced(rw, p)

	if !err {
		if forced {
			if value == "positive" {
				index := rand.Intn(len(positive))
				ch = Choice{positive[index], true, value}
			} else if value == "negative" {
				index := rand.Intn(len(negative))
				ch = Choice{negative[index], true, value}
			} else {
				index := rand.Intn(len(neutral))
				ch = Choice{neutral[index], true, value}
			}
		} else {
			if randomNr%3 == 1 {
				index := rand.Intn(len(positive))
				ch = Choice{positive[index], false, "positive"}
			} else if randomNr%3 == 0 {
				index := rand.Intn(len(negative))
				ch = Choice{negative[index], false, "negative"}
			} else {
				index := rand.Intn(len(neutral))
				ch = Choice{neutral[index], false, "neutral"}
			}
		}
		writeResponse(rw, ch)
	}
}

func isForced(rw http.ResponseWriter, p httprouter.Params) (bool, bool, string) {
	var forced bool
	var value string
	var errorr bool
	force := p.ByName("force")

	if force == "" {
		return false, false, ""
	}

	if force == "positive" {
		forced = true
		value = force
		errorr = false
	} else if force == "negative" {
		forced = true
		value = force
		errorr = false
	} else if force == "neutral" {
		forced = true
		value = force
		errorr = false
	} else {
		err1 := Error{"Must be either positive or negative"}
		js, _ := json.Marshal(err1)
		http.Error(rw, string(js), http.StatusBadRequest)
		errorr = true
	}
	return errorr, forced, value
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
