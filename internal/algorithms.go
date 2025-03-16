package internal

import (
	"math/rand/v2"
)

type Algorithm interface {
	GetNextHost() *Host
}

type RandomAlgorithm struct {
	logger Logger
	hosts  []*Host
}

func NewRandomAlgorithm(hosts []*Host, logger Logger) Algorithm {
	return &RandomAlgorithm{
		logger: logger,
		hosts:  hosts,
	}
}

func (r *RandomAlgorithm) GetNextHost() *Host {
	var index = rand.IntN(len(r.hosts))
	return r.hosts[index]
}

type RoundRobinAlgorithm struct {
	currHost int
	logger   Logger
	hosts    []*Host
}

func NewRoundRobinAlgorithm(hosts []*Host, logger Logger) Algorithm {
	return &RoundRobinAlgorithm{
		currHost: 0,
		logger:   logger,
		hosts:    hosts,
	}
}

func (r *RoundRobinAlgorithm) GetNextHost() *Host {
	r.currHost = (r.currHost + 1) % len(r.hosts)
	return r.hosts[r.currHost]
}

type WeightedRoundRobinAlgorithm struct {
	maxWeight  int
	currWeight int
	currHost   int
	gcdWeight  int
	hosts      []*Host
	logger     Logger
}

func NewWeightedRoundRobinAlgorithm(hosts []*Host, logger Logger) Algorithm {
	return &WeightedRoundRobinAlgorithm{
		maxWeight:  getMaxWeight(hosts),
		gcdWeight:  getGcd(hosts),
		currHost:   -1,
		currWeight: 0,
		hosts:      hosts,
		logger:     logger,
	}
}

func (w *WeightedRoundRobinAlgorithm) GetNextHost() *Host {
	for {
		w.currHost = (w.currHost + 1) % len(w.hosts)
		if w.currHost == 0 {
			w.currWeight = w.currWeight - w.gcdWeight
			if w.currWeight <= 0 {
				w.currWeight = w.maxWeight
			}
		}

		if w.hosts[w.currHost].Weight >= w.currWeight {
			return w.hosts[w.currHost]
		}
	}
}

func getGcd(hosts []*Host) int {
	var gcd = 0
	for _, host := range hosts {
		gcd = calculateGCD(gcd, host.Weight)
	}
	return gcd
}

func getMaxWeight(hosts []*Host) int {
	var maxWeight = 0
	for _, host := range hosts {
		if host.Weight > maxWeight {
			maxWeight = host.Weight
		}
	}
	return maxWeight
}

func calculateGCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
