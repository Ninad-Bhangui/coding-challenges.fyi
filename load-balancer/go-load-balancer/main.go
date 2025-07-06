package main

import (
	"flag"
	"fmt"
	"load-balancer/lb"
	"log"
	"net/http"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	log.Printf("INFO: Starting Go Load Balancer")
	var serverUrls arrayFlags
	flag.Var(&serverUrls, "url", "URLs of server")
	flag.Parse()
	if len(serverUrls) == 0 {
		log.Printf("ERROR: No server URLs provided")
		flag.Usage()
		return
	}
	log.Printf("INFO: Parsed %d server URLs from command line", len(serverUrls))
	lb := lb.NewRoundRobinLb(serverUrls)
	s := &http.Server{
		Addr:    ":8080",
		Handler: &lb,
	}
	log.Printf("INFO: Starting HTTP server on port 8080")
	fmt.Printf("Listening on port 8080")
	log.Fatal(s.ListenAndServe())
}
