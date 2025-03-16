package internal

type LoadBalancer struct {
	numHosts  int
	algorithm Algorithm
}

func NewLoadBalancer(hosts []*Host, strategy LoadBalancingStrategy, logger Logger) *LoadBalancer {
	algorithms := map[LoadBalancingStrategy]func([]*Host, Logger) Algorithm{
		StrategyRandom:             NewRandomAlgorithm,
		StrategyRoundRobin:         NewRoundRobinAlgorithm,
		StrategyWeightedRoundRobin: NewWeightedRoundRobinAlgorithm,
	}

	algorithm, exists := algorithms[strategy]
	if !exists {
		logger.Errorf("Unknown strategy: %s", strategy)
		return nil
	}

	return &LoadBalancer{
		algorithm: algorithm(hosts, logger),
		numHosts:  len(hosts),
	}
}

func (l *LoadBalancer) GetNextHost() *Host {
	for i := 0; i < l.numHosts; i++ {
		server := l.algorithm.GetNextHost()
		if server.Status == HEALTHY {
			return server
		}
	}
	return nil
}
