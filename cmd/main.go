package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DataStub struct {
	Decription string          `json:"decription"`
	Request    RequestRawData  `json:"request"`
	Response   ResponseRawData `json:"response"`
}

type RequestRawData struct {
	Path   string            `json:"path"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Query  map[string]string `json:"query"`
}

type ResponseRawData struct {
	Status int               `json:"status"`
	Header map[string]string `json:"header"`
	Body   interface{}       `json:"body"`
}

func main() {
	var port, pathfile string
	flag.StringVar(&port, "port", ":3000", "port service on")
	flag.StringVar(&pathfile, "pathfile", "", "path file data stub")
	flag.Parse()

	rawFileData, err := ioutil.ReadFile(pathfile)
	if err != nil {
		log.Fatal(err)
	}

	var rawData DataStub
	if err = json.Unmarshal(rawFileData, &rawData); err != nil {
		log.Panicln(err)
	}

	http.HandleFunc(rawData.Request.Path, func(w http.ResponseWriter, r *http.Request) {
		for kay, value := range rawData.Request.Header {
			w.Header().Set(kay, value)
		}
		responseBody, err := json.Marshal(rawData.Response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", responseBody)
	})
	log.Println("server start on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
