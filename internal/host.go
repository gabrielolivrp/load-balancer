package internal

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type HealthResult string

const (
	UNHEALTHY HealthResult = "un_healthy"
	HEALTHY   HealthResult = "healthy"
)

type Host struct {
	URL               *url.URL
	Status            HealthResult
	ConnectionTimeout time.Duration
	Weight            int // used by RoundRobinAlgorithm
	FailStreak        int
	SuccessStreak     int
}

func NewHost(url *url.URL, weight int, connectionTimeout time.Duration) *Host {
	return &Host{
		URL:               url,
		Status:            UNHEALTHY,
		Weight:            weight,
		FailStreak:        0,
		SuccessStreak:     0,
		ConnectionTimeout: connectionTimeout * time.Second,
	}
}

func (s *Host) SetStreak(status HealthResult) {
	if status == UNHEALTHY {
		s.FailStreak++
		s.SuccessStreak = 0
	} else {
		s.SuccessStreak++
		s.FailStreak = 0
	}
}

func (s *Host) ReverseProxy() *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(s.URL)
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: s.ConnectionTimeout,
		}).DialContext,
	}
	return proxy

}
