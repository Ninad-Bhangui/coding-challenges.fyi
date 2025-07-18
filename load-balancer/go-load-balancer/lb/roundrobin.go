package lb

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServerDetail struct {
	url     string
	healthy bool
}

type RoundRobinLB struct {
	serverUrls          []ServerDetail
	currentServerIndex  int
	healthCheckInterval time.Duration
	mu                  sync.Mutex
}

func NewRoundRobinLb(urls []string, healthCheckInterval time.Duration) *RoundRobinLB {
	log.Printf("INFO: Initializing Round Robin Load Balancer with %d servers: %v", len(urls), urls)
	serverUrls := make([]ServerDetail, 0, len(urls))
	for _, url := range urls {
		serverUrls = append(serverUrls, ServerDetail{url: url, healthy: true})
	}
	lb := &RoundRobinLB{
		serverUrls:          serverUrls,
		healthCheckInterval: healthCheckInterval,
		currentServerIndex:  0,
	}
	go lb.healthCheck()
	return lb
}

func (lb *RoundRobinLB) getHealthyServers() []ServerDetail {
	//This is for test cases only because healthcheck loop may be modifying resource when we read this. ServeHTTP which is caller has it's own mutex lock
	lb.mu.Lock()
	defer lb.mu.Unlock()
	return lb.getHealthyServersUnlocked()
}

func (lb *RoundRobinLB) getHealthyServersUnlocked() []ServerDetail {
	healthyServers := make([]ServerDetail, 0, len(lb.serverUrls))
	for _, server := range lb.serverUrls {
		if server.healthy {
			healthyServers = append(healthyServers, server)
		}
	}
	return healthyServers
}
func (lb *RoundRobinLB) healthCheck() {
	for {
		var wg sync.WaitGroup
		for i := range lb.serverUrls {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				lb.healthCheckUrl(index)
			}(i)
		}
		wg.Wait()
		time.Sleep(lb.healthCheckInterval)
	}
}

func (lb *RoundRobinLB) healthCheckUrl(index int) {
	serverUrl := lb.serverUrls[index].url
	res, err := http.Get(serverUrl)

	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err != nil {
		lb.serverUrls[index].healthy = false
		log.Printf("DEBUG: %s is unhealthy due to error: %s", serverUrl, err)
		return
	}
	if res.StatusCode != http.StatusOK {
		lb.serverUrls[index].healthy = false
		return
	}
	lb.serverUrls[index].healthy = true
}

func (lb *RoundRobinLB) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lb.mu.Lock()
	healthyServers := lb.getHealthyServersUnlocked()
	if len(healthyServers) == 0 {
		lb.mu.Unlock()
		log.Printf("ERROR: no healthy servers available")
		w.WriteHeader(http.StatusBadGateway) //TODO: Check if right status code
		io.WriteString(w, "Sorry, no healthy servers available")
		return
	}
	if len(healthyServers) <= lb.currentServerIndex {
		lb.currentServerIndex = 0 //TODO: Currently marking to 0 if the next relevant server no longer fits in healthy server list. What to do here?
	}
	serverUrl := healthyServers[lb.currentServerIndex]
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
