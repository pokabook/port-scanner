package main

import (
	"testing"
	"time"
)

func BenchmarkStartWithMonitor(b *testing.B) {
	ip := "127.0.0.1"
	f := 1
	l := 65535
	timeout := 500 * time.Millisecond

	for n := 0; n < b.N; n++ {
		StartWithMonitor(ip, f, l, timeout)
	}
}

func BenchmarkStartWithSemaphore(b *testing.B) {
	ip := "127.0.0.1"
	f := 1
	l := 65535
	timeout := 500 * time.Millisecond

	for n := 0; n < b.N; n++ {
		StartWithSemaphore(ip, f, l, timeout)
	}
}

func BenchmarkStartWithMutex(b *testing.B) {
	ip := "127.0.0.1"
	f := 1
	l := 65535
	timeout := 500 * time.Millisecond

	for n := 0; n < b.N; n++ {
		StartWithMutex(ip, f, l, timeout)
	}
}
