package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type servers map[string]Server
type Server struct {
	Ip          string `json:"ip"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Pins        []struct {
		Name      string `json:"name"`
		Posistion int    `json:"position"`
	}
}
type Body struct {
	Server string
	Pin    string
}

func getBody(req *http.Request) Body {
	decoder := json.NewDecoder(req.Body)
	var reqBody Body
	err := decoder.Decode(&reqBody)
	if err != nil {
		panic(err)
	}
	return reqBody
}

func makeRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
	}
}

func on(w http.ResponseWriter, req *http.Request) {
	b := getBody(req)
	sconfig := getConfig(b.Server)
	s := fmt.Sprintf("http://%s/%s/ON", sconfig.Ip, b.Pin)
	go makeRequest(s)
	fmt.Fprintf(w, "on\n")
}

func off(w http.ResponseWriter, req *http.Request) {
	b := getBody(req)
	sconfig := getConfig(b.Server)
	s := fmt.Sprintf("http://%s/%s/OFF", sconfig.Ip, b.Pin)
	go makeRequest(s)
	fmt.Fprintf(w, "off\n")
}

func config(w http.ResponseWriter, req *http.Request) {
	sconfig := readConfig()
	r, err := json.Marshal(sconfig)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "there was an error\n")
	}
	fmt.Fprint(w, string(r))
}

func readConfig() servers {
	yfile, err := os.ReadFile("servers.json")

	if err != nil {
		panic(err)
	}
	config := servers{}
	err2 := json.Unmarshal(yfile, &config)

	if err2 != nil {
		panic(err2)
	}
	return config
}

func getConfig(key string) Server {
	config := readConfig()
	return config[key]
}

func main() {
	http.HandleFunc("/on", on)
	http.HandleFunc("/off", off)
	http.HandleFunc("/config", config)

	http.ListenAndServe(":8090", nil)
}
