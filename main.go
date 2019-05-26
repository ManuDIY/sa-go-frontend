package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func main() {

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {

		url := "http://192.168.137.204/sentiment"
		var jsonStr = []byte(`{"sentence":"Golang is awesome."}`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("User-Agent", "User-Agent: Go client")
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

		defer resp.Body.Close()
		w.Write(body)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Healthy"))
	})
	http.ListenAndServe(":8090", nil)
}
