package main

import (
	"fmt"
	"time"
)

import (
	"context"
	"golang.org/x/sync/semaphore"
	"net"
	"strings"
	"sync"
)

type PortScanner struct {
	ip   string
	lock *semaphore.Weighted
}

func ScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		}
		return
	}

	conn.Close()
}

func (ps *PortScanner) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			ScanPort(ps.ip, port, timeout)
		}(port)
	}
}

func main() {
	start := time.Now()

	ps := &PortScanner{
		ip:   "127.0.0.1",
		lock: semaphore.NewWeighted(256),
	}
	ps.Start(1, 65535, 500*time.Millisecond)
	fmt.Println(time.Since(start))
}
