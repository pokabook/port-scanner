package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	ip        string
	lock      sync.Mutex
	cond      *sync.Cond
	available int64
}

type Scanner interface {
	StartScan(ip string, f, l int, timeout time.Duration)
}

func Ulimit() int64 {
	out, err := exec.Command("ulimit", "-n").Output()
	if err != nil {
		panic(err)
	}

	s := strings.TrimSpace(string(out))

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

type SemaphoreScanner struct{}
type MutexScanner struct{}
type MonitorScanner struct{}
type SequentialScanner struct{}

func (SemaphoreScanner) StartScan(ip string, f, l int, timeout time.Duration) {
	StartWithSemaphore(ip, f, l, timeout)
}

func (MutexScanner) StartScan(ip string, f, l int, timeout time.Duration) {
	StartWithMutex(ip, f, l, timeout)
}

func (MonitorScanner) StartScan(ip string, f, l int, timeout time.Duration) {
	StartWithMonitor(ip, f, l, timeout)
}

func (SequentialScanner) StartScan(ip string, f, l int, timeout time.Duration) {
	StartWithSequential(ip, f, l, timeout)
}

func TCPScanPort(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			TCPScanPort(ip, port, timeout)
		}
		return
	}

	//fmt.Println(target)

	conn.Close()
}

func StartWithSemaphore(ip string, f, l int, timeout time.Duration) {
	sem := semaphore.NewWeighted(Ulimit())

	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		sem.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer sem.Release(1)
			defer wg.Done()
			TCPScanPort(ip, port, timeout)
		}(port)
	}
}

func StartWithMutex(ip string, f, l int, timeout time.Duration) {
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		lock.Lock()
		wg.Add(1)
		go func(port int) {
			defer lock.Unlock()
			defer wg.Done()
			TCPScanPort(ip, port, timeout)
		}(port)
	}
}

func StartWithMonitor(ip string, f, l int, timeout time.Duration) {
	ps := &PortScanner{
		ip:        ip,
		available: Ulimit(),
	}
	ps.cond = sync.NewCond(&ps.lock)

	for port := f; port <= l; port++ {
		ps.lock.Lock()
		for ps.available <= 0 {
			ps.cond.Wait()
		}
		ps.available--
		ps.lock.Unlock()

		go func(port int) {
			TCPScanPort(ip, port, timeout)

			ps.lock.Lock()
			ps.available++
			ps.cond.Signal()
			ps.lock.Unlock()
		}(port)
	}
}

func StartWithSequential(ip string, f, l int, timeout time.Duration) {
	for port := f; port <= l; port++ {
		TCPScanPort(ip, port, timeout)
	}
}

func main() {
	start := time.Now()

	ip := "127.0.0.1"
	scanner := MonitorScanner{}
	scanner.StartScan(ip, 1, 65535, 500*time.Millisecond)
	fmt.Println(time.Since(start))
}
