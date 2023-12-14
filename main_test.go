package main

import (
	"testing"
	"time"
)

func BenchmarkStartWithMonitor(b *testing.B) {
	scanner := MonitorScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
}

func BenchmarkStartWithSemaphore(b *testing.B) {
	scanner := SemaphoreScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
}

func BenchmarkStartWithMutex(b *testing.B) {
	scanner := MutexScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
}

func BenchmarkStartWithSequential(b *testing.B) {
	scanner := SequentialScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
}
