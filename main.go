package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//Balancers - taken from -b input
type Balancers struct {
	serverNames     string
	chosenBalancers []string
}

//Alive - available load balancers after health check
type Alive struct {
	areAlive []string
}

func randomHost() {
	var a Alive
	for i, j := range a.areAlive {
		fmt.Println("randomHost")
		fmt.Println(i)
		_ = j
	}
}

func reverseProxy() {

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "zamunda.net",
	})
	proxy.Transport = &http.Transport{DialTLS: dialTLS}

	// Change req.Host so badssl.com host check is passed
	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	log.Fatal(http.ListenAndServe("127.0.0.1:3000", proxy))
}

func dialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}

	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]

	// Verify here
	cert.VerifyHostname(host)
	log.Println(cert.Subject)

	return tlsConn, nil
}

func healthCheck(b Balancers) {
	var a Alive
	doExist := map[int]bool{}
	result := []int{}
	for {
		for _, s := range b.chosenBalancers {
			_, err := http.Get(s)
			if err == nil {
				a.areAlive = append(a.areAlive, s)
			}
		}
		fmt.Println("")
		fmt.Println("Are Alive:")
		for i, j := range a.areAlive {
			fmt.Println(i)
			fmt.Println(j)
		}
	}
}

func main() {
	// var serverNames string
	var b Balancers
	flag.StringVar(&b.serverNames, "b", "", "Input balancers")
	flag.Parse()

	serverNameWithPrefix := strings.Split(b.serverNames, " ")
	for _, i := range serverNameWithPrefix {
		i = "http://" + i
		b.chosenBalancers = append(b.chosenBalancers, i)
	}

	fmt.Println("Chosen Balancers")
	for _, server := range b.chosenBalancers {
		fmt.Println(server)
	}

	go healthCheck(b)

	reverseProxy()
}
