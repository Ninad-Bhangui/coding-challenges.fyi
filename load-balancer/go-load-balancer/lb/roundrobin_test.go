package lb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRoundRobinLB(t *testing.T) {
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "hello from backend1")
	}))
	defer backend1.Close()
	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "hello from backend2")
	}))
	defer backend2.Close()

	lb := NewRoundRobinLb([]string{backend1.URL, backend2.URL}, 10*time.Second)

	rr1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/", nil)

	lb.ServeHTTP(rr1, req1)

	if rr1.Code != 200 {
		t.Errorf("Expected 200, got %d", rr1.Code)
	}
	if rr1.Body.String() != "hello from backend1" {
		t.Errorf("Expected backend1, got %s", rr1.Body.String())
	}
	rr2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/", nil)

	lb.ServeHTTP(rr2, req2)

	if rr2.Code != 200 {
		t.Errorf("Expected 200, got %d", rr2.Code)
	}
	if rr2.Body.String() != "hello from backend2" {
		t.Errorf("Expected backend2, got %s", rr2.Body.String())
	}

}

func TestRoundRobinRace(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer backend.Close()

	lb := NewRoundRobinLb([]string{backend.URL, backend.URL, backend.URL}, 10*time.Second)

	var wg sync.WaitGroup
	numGoroutines := 2

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			lb.ServeHTTP(rr, req)
		}()
	}

	wg.Wait()
}

func TestGetOnlyHealthyRoundRobinLB(t *testing.T) {
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "hello from backend1")
	}))
	defer backend1.Close()
	unhealthyBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error")
	}))
	defer unhealthyBackend.Close()

	lb := NewRoundRobinLb([]string{backend1.URL, unhealthyBackend.URL}, 2*time.Second)
	time.Sleep(4 * time.Second)
	healthyServers := lb.getHealthyServers()
	if len(healthyServers) == 2 {
		t.Errorf("Error: Expected there to be only 1 healhyServer")
	}

}

// func TestUnhealthyRoundRobinLB(t *testing.T) {
// 	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(200)
// 		fmt.Fprintf(w, "hello from backend1")
// 	}))
// 	defer backend1.Close()
// 	unhealthyBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(500)
// 		fmt.Fprintf(w, "Error")
// 	}))
// 	defer unhealthyBackend.Close()
//
// 	lb := NewRoundRobinLb([]string{backend1.URL, unhealthyBackend.URL}, 10*time.Second)
//
// 	//TODO: Not working
// 	for i := 0; i < 10; i++ {
// 		rr := httptest.NewRecorder()
// 		req1, err := http.NewRequest("GET", "/", nil)
// 		lb.ServeHTTP(rr, req1)
// 		if err != nil {
// 			t.Errorf("ERROR: should not have errored")
// 		}
// 		if rr.Code == 500 {
// 			t.Errorf("Should not have gotten 500 as lb should have redirected to healthy URLs only")
// 		}
//
// 	}
//
// }
