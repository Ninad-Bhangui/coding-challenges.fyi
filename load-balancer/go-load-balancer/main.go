package main

import (
	"flag"
	"fmt"
	"load-balancer/lb"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
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
	// Configure structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	
	slog.Info("Starting Go Load Balancer")
	var serverUrls arrayFlags
	flag.Var(&serverUrls, "url", "URLs of server")
	flag.Parse()
	if len(serverUrls) == 0 {
		slog.Error("No server URLs provided")
		flag.Usage()
		return
	}
	
	// Validate URLs
	for _, urlStr := range serverUrls {
		if _, err := url.Parse(urlStr); err != nil {
			slog.Error("Invalid URL", "url", urlStr, "error", err)
			return
		}
		if parsedURL, _ := url.Parse(urlStr); parsedURL.Scheme == "" || parsedURL.Host == "" {
			slog.Error("URL must include scheme (http/https) and host", "url", urlStr)
			return
		}
	}
	
	slog.Info("Parsed server URLs from command line", "count", len(serverUrls), "urls", serverUrls)
	lb := lb.NewRoundRobinLb(serverUrls, 10*time.Second)
	s := &http.Server{
		Addr:    ":8080",
		Handler: lb,
	}
	slog.Info("Starting HTTP server", "port", 8080)
	fmt.Printf("Listening on port 8080")
	if err := s.ListenAndServe(); err != nil {
		slog.Error("Failed to start HTTP server", "error", err)
	}
}
