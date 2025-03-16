package internal

type Protocol string

const (
	HTTPProtocol  Protocol = "http"
	HTTPSProtocol Protocol = "https"
)

type LoadBalancingStrategy string

const (
	StrategyRoundRobin         LoadBalancingStrategy = "round_robin"
	StrategyRandom             LoadBalancingStrategy = "random"
	StrategyWeightedRoundRobin LoadBalancingStrategy = "weighted_round_robin"
)

type LogOutputConfig struct {
}

type LoggingConfig struct {
	Enabled bool     `json:"enabled"`
	Level   LogLevel `json:"level"`
	File    string   `json:"file"`
}

type HostConfig struct {
	Address        string `json:"address"`
	Port           int    `json:"port"`
	Weight         int    `json:"weight"`
	MaxConnections int    `json:"max_connections"`
}

type HealthCheckConfig struct {
	Endpoint           string `json:"endpoint"`
	Interval           int    `json:"interval"`
	Timeout            int    `json:"timeout"`
	UnhealthyThreshold int    `json:"unhealthy_threshold"`
	HealthyThreshold   int    `json:"healthy_threshold"`
}

type ListenerConfig struct {
	Port     int      `json:"port"`
	Protocol Protocol `json:"protocol"`
}

type ServiceConfig struct {
	Name     string         `json:"name"`
	Listener ListenerConfig `json:"listen"`
}

type LoadBalancing struct {
	Algorithm         LoadBalancingStrategy `json:"algorithm"`
	ConnectionTimeout int                   `json:"connection_timeout"`
}

type LoadBalancerConfig struct {
	Service           ServiceConfig     `json:"service"`
	LoadBalancing     LoadBalancing     `json:"load_balancing"`
	ConnectionTimeout int               `json:"connection_timeout"`
	HealthCheck       HealthCheckConfig `json:"health_checks"`
	Hosts             []HostConfig      `json:"hosts"`
	Logging           LoggingConfig     `json:"logging"`
}
