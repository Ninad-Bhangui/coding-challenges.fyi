package lb

import (
	"io"
	"log/slog"
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
	client              *http.Client
}

func NewRoundRobinLb(urls []string, healthCheckInterval time.Duration) *RoundRobinLB {
	slog.Info("Initializing Round Robin Load Balancer", "server_count", len(urls), "urls", urls)
	serverUrls := make([]ServerDetail, 0, len(urls))
	for _, url := range urls {
		serverUrls = append(serverUrls, ServerDetail{url: url, healthy: true})
	}
	lb := &RoundRobinLB{
		serverUrls:          serverUrls,
		healthCheckInterval: healthCheckInterval,
		currentServerIndex:  0,
	}
	lb.client = &http.Client{}
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
	res, err := lb.client.Get(serverUrl)

	lb.mu.Lock()
	defer lb.mu.Unlock()

	if err != nil {
		lb.serverUrls[index].healthy = false
		slog.Debug("Server health check failed", "url", serverUrl, "error", err)
		return
	}

	defer res.Body.Close()
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
		slog.Error("No healthy servers available")
		w.WriteHeader(http.StatusServiceUnavailable)
		io.WriteString(w, "Sorry, no healthy servers available")
		return
	}
	if len(healthyServers) <= lb.currentServerIndex {
		lb.currentServerIndex = 0
	}
	serverUrl := healthyServers[lb.currentServerIndex]
	lb.currentServerIndex = (lb.currentServerIndex + 1) % len(healthyServers)
	slog.Debug("Routing request", "method", req.Method, "path", req.URL.Path, "server", serverUrl.url, "index", lb.currentServerIndex)
	lb.mu.Unlock()
	err := lb.serve(w, req, serverUrl.url)
	if err != nil {
		slog.Error("Failed to serve request", "server", serverUrl.url, "error", err)
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Something went wrong")
	}
}

func (lb *RoundRobinLB) serve(w http.ResponseWriter, req *http.Request, url string) error {
	slog.Debug("Creating proxy request", "url", url)
	lbReq, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		slog.Error("Failed to create proxy request", "error", err)
		return err
	}
	lbReq.Header = req.Header
	slog.Debug("Sending request to backend server", "url", url)
	res, err := lb.client.Do(lbReq)
	if err != nil {
		slog.Error("Failed to connect to backend server", "url", url, "error", err)
		http.Error(w, "Could not connect", http.StatusBadGateway)
		return nil
	}
	defer res.Body.Close()
	slog.Debug("Received response from backend", "url", url, "status", res.StatusCode)
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)

	return nil

}
