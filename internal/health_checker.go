package internal

import (
	"net/http"
	"net/url"
	"time"
)

type HealthCheck struct {
	hosts              []*Host
	interval           time.Duration
	timeout            time.Duration
	path               string
	logger             Logger
	unhealthyThreshold int
	healthyThreshold   int
}

func NewHealthCheck(
	hosts []*Host,
	interval time.Duration,
	timeout time.Duration,
	path string,
	unhealthyThreshold int,
	healthyThreshold int,
	logger Logger,
) *HealthCheck {
	return &HealthCheck{
		hosts:              hosts,
		interval:           interval * time.Second,
		timeout:            timeout * time.Second,
		path:               path,
		unhealthyThreshold: unhealthyThreshold,
		healthyThreshold:   healthyThreshold,
		logger:             logger,
	}
}

func (h *HealthCheck) Start() {
	for _, host := range h.hosts {
		go h.monitorHost(host)
	}
}

func (h *HealthCheck) monitorHost(host *Host) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for range ticker.C {
		res := h.healthCheck(host)
		host.SetStreak(res)

		if host.FailStreak >= h.unhealthyThreshold {
			h.logger.Errorf("Host %s is unhealthy", host.URL.String())
			host.Status = UNHEALTHY
		} else if host.SuccessStreak >= h.healthyThreshold {
			host.Status = HEALTHY
		}
	}
}

func (h *HealthCheck) healthCheck(host *Host) HealthResult {
	client := &http.Client{
		Timeout: h.timeout,
	}
	url, err := url.JoinPath(host.URL.String(), h.path)
	if err != nil {
		return UNHEALTHY
	}

	res, err := client.Get(url)
	if err == nil && res.StatusCode == http.StatusOK {
		return HEALTHY
	}

	return UNHEALTHY
}
