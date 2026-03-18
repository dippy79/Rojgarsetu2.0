package proxy

import (
	"math/rand"
	"sync"
	"time"
)

// Rotator manages proxy rotation
type Rotator struct {
	mu      sync.Mutex
	proxies []string
	index   int
}

// NewRotator creates a new proxy rotator
func NewRotator() *Rotator {
	return &Rotator{
		proxies: []string{}, // Load from config in production
		index:   0,
	}
}

// Next returns the next proxy URL
func (r *Rotator) Next() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.proxies) == 0 {
		return ""
	}

	proxy := r.proxies[r.index]
	r.index = (r.index + 1) % len(r.proxies)
	return proxy
}

// Add adds a proxy to the rotation
func (r *Rotator) Add(proxy string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.proxies = append(r.proxies, proxy)
}

// Remove removes a proxy from the rotation
func (r *Rotator) Remove(proxy string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.proxies {
		if p == proxy {
			r.proxies = append(r.proxies[:i], r.proxies[i+1:]...)
			break
		}
	}
}

// UserAgent returns a random user agent
func UserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
	rand.Seed(time.Now().UnixNano())
	return agents[rand.Intn(len(agents))]
}
