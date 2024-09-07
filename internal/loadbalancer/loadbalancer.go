package loadbalancer

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type LoadBalancer struct {
	servers             []string
	index               int
	mu                  sync.Mutex
	healthy             []bool
	checkHealthInterval time.Duration
}

func NewLoadBalancer(servers []string, checkInterval time.Duration) *LoadBalancer {
	lb := &LoadBalancer{
		servers:             servers,
		checkHealthInterval: checkInterval,
		healthy:             make([]bool, len(servers)),
	}

	go startHealthChecks(lb)

	return lb
}

func startHealthChecks(lb *LoadBalancer) {
	for {
		time.Sleep(lb.checkHealthInterval)
		for i, server := range lb.servers {
			if err := lb.checkHealth(server); err != nil {
				lb.mu.Lock()
				lb.healthy[i] = false
				lb.mu.Unlock()
			} else {
				lb.mu.Lock()
				lb.healthy[i] = true
				lb.mu.Unlock()
			}
			log.Println("Server " + strconv.FormatInt(int64(i), 10) + " " + strconv.FormatBool(lb.healthy[i]))
		}
	}
}

func (lb *LoadBalancer) checkHealth(server string) error {
	response, err := http.Get(server + "/health")

	if err != nil || response.StatusCode != http.StatusOK {
		return err
	}
	return nil

}

func (lb *LoadBalancer) NextServer() string {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for {
		if lb.healthy[lb.index] {
			server := lb.servers[lb.index]
			lb.index = (lb.index + 1) % len(lb.servers)
			return server
		} else {
			lb.index = (lb.index + 1) % len(lb.servers)
		}
	}
}
