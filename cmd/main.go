package main

import (
	"encoding/json"
	"fmt"
	"load-balancer/internal"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}

	logger := internal.NewLogger(
		config.Logging.Enabled,
		config.Logging.Level,
		config.Logging.File,
	)

	hosts, err := createHosts(config, logger)
	if err != nil {
		logger.Errorf("Error creating hosts: %s", err)
		os.Exit(1)
	}

	healthCheck := internal.NewHealthCheck(
		hosts,
		time.Duration(config.HealthCheck.Interval),
		time.Duration(config.HealthCheck.Timeout),
		config.HealthCheck.Endpoint,
		config.HealthCheck.UnhealthyThreshold,
		config.HealthCheck.HealthyThreshold,
		logger,
	)
	healthCheck.Start()

	loadBalancer := internal.NewLoadBalancer(
		hosts,
		config.LoadBalancing.Algorithm,
		logger,
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := loadBalancer.GetNextHost()
		if host == nil {
			logger.Errorf("No healthy hosts available")
			http.Error(w, "No healthy hosts available", http.StatusServiceUnavailable)
			return
		}

		host.ReverseProxy().ServeHTTP(w, r)
	})

	port := fmt.Sprintf(":%d", config.Service.Listener.Port)
	logger.Infof("Server is running on port %s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		logger.Errorf("Error starting server: %s", err)
		os.Exit(1)
	}
}

func loadConfig(path string) (*internal.LoadBalancerConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config internal.LoadBalancerConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func createHosts(config *internal.LoadBalancerConfig, logger internal.Logger) ([]*internal.Host, error) {
	var hosts []*internal.Host
	for _, hostConfig := range config.Hosts {
		hostURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", config.Service.Listener.Protocol, hostConfig.Address, hostConfig.Port))
		if err != nil {
			logger.Errorf("Error parsing host URL: %s", err)
			return nil, err
		}
		host := internal.NewHost(
			hostURL,
			hostConfig.Weight,
			time.Duration(config.LoadBalancing.ConnectionTimeout),
		)
		hosts = append(hosts, host)
	}
	return hosts, nil
}
