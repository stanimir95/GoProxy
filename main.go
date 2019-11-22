package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
)

//FixedHTTPPrefix please stop
var FixedHTTPPrefix []string

func getServers() []string {

	var serverList string
	// var fixedHTTPPrefix []string

	flag.StringVar(&serverList, "b", "", "Load balancer goes here")
	flag.Parse()

	servers := strings.Split(serverList, " ")
	for _, server := range servers {
		server = "http://" + server
		FixedHTTPPrefix = append(FixedHTTPPrefix, server)
	}
	return FixedHTTPPrefix
}

func checkServerStatus(fixedHTTPPrefix []string) []string {

	var liveServers = fixedHTTPPrefix
	for i, server := range liveServers {
		req, err := http.Get(server)
		if err == nil {
			liveServers = append(liveServers, server)
		} else {
			liveServers[i] = liveServers[len(liveServers)-1]
			return liveServers[:len(liveServers)-1]
		}

		defer req.Body.Close()
	}
	return liveServers

}

func main() {
	fmt.Println("	fmt.Println(getServers())")
	fmt.Println(getServers())

	fmt.Println("	fmt.Println(FixedHTTPPrefix)")
	fmt.Println(FixedHTTPPrefix)

	fmt.Println("	VAJNOOOO")
	for {
		fmt.Println(checkServerStatus(FixedHTTPPrefix))
	}
}
