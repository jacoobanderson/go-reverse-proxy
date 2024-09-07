package loadbalancer

import "sync"

type LoadBalancer struct {
	servers []string
	index   int
	mu      sync.Mutex
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{servers: servers}
}

func (lb *LoadBalancer) NextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}
