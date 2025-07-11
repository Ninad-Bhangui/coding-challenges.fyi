package lb

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type ServerDetail struct {
	url     string
	healthy bool
}

type RoundRobinLB struct {
	serverUrls         []ServerDetail
	currentServerIndex int
	mu                 sync.Mutex
}

func NewRoundRobinLb(urls []string) RoundRobinLB {
	log.Printf("INFO: Initializing Round Robin Load Balancer with %d servers: %v", len(urls), urls)
	serverUrls := make([]ServerDetail, 0, len(urls))
	for _, url := range urls {
		serverUrls = append(serverUrls, ServerDetail{url: url, healthy: true})

	}
	return RoundRobinLB{
		serverUrls:         serverUrls,
		currentServerIndex: 0,
	}

}

// func (lb *RoundRobinLB) healthCheck() {
// 	for i, url := range lb.serverUrls {
//
// 	}
// }
//
// func (lb *RoundRobinLB) healthCheckUrl(index int) {
// 	url := lb.serverUrls[index]
// 	res, err := http.Get(url)
// 	if err != nil {
//
// 	}
//
// }

func (lb *RoundRobinLB) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lb.mu.Lock()
	serverUrl := lb.serverUrls[lb.currentServerIndex]
	lb.currentServerIndex = (lb.currentServerIndex + 1) % len(lb.serverUrls)
	log.Printf("DEBUG: Routing request %s %s to server %s (index: %d)", req.Method, req.URL.Path, serverUrl.url, lb.currentServerIndex)
	lb.mu.Unlock()
	err := lb.serve(w, req, serverUrl.url)
	if err != nil {
		log.Printf("ERROR: Failed to serve request to %s: %v", serverUrl.url, err)
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Something went wrong")
	}
}

func (lb *RoundRobinLB) serve(w http.ResponseWriter, req *http.Request, url string) error {
	log.Printf("DEBUG: Creating proxy request to %s", url)
	lbReq, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		log.Printf("ERROR: Failed to create proxy request: %v", err)
		return err
	}
	lbReq.Header = req.Header
	log.Printf("DEBUG: Sending request to backend server %s", url)
	res, err := http.DefaultClient.Do(lbReq)
	if err != nil {
		log.Printf("ERROR: Failed to connect to backend server %s: %v", url, err)
		http.Error(w, "Could not connect", http.StatusBadGateway)
		return nil
	}
	defer res.Body.Close()
	log.Printf("DEBUG: Received response from %s with status %d", url, res.StatusCode)
	w.WriteHeader(res.StatusCode)
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	io.Copy(w, res.Body)

	return nil

}
