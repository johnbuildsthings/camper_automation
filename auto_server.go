package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

type servers map[string]Server
type Server struct {
	Ip          string `yaml:"ip"`
	Description string `yaml:"description"`
	Enabled     bool   `yaml:"enabled"`
	Pins        []struct {
		Name      string `yaml:"name"`
		Posistion int    `yaml:"position"`
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

func getConfig(server string) Server {
	yfile, err := ioutil.ReadFile("servers.yaml")

	if err != nil {
		panic(err)
	}
	config := servers{}
	err2 := yaml.Unmarshal(yfile, &config)

	if err2 != nil {
		panic(err2)
	}
	return config[server]
}

func main() {
	http.HandleFunc("/on", on)
	http.HandleFunc("/off", off)

	http.ListenAndServe(":8090", nil)
}
