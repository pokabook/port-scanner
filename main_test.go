package main

import (
	"testing"
	"time"
)

func BenchmarkStartWithMonitor(b *testing.B) {
	start := time.Now()

	scanner := MonitorScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
	b.ReportMetric(float64(time.Since(start).Milliseconds()), "ms")
}

func BenchmarkStartWithSemaphore(b *testing.B) {
	start := time.Now()

	scanner := SemaphoreScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
	b.ReportMetric(float64(time.Since(start).Milliseconds()), "ms")
}

func BenchmarkStartWithMutex(b *testing.B) {
	start := time.Now()

	scanner := MutexScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
	b.ReportMetric(float64(time.Since(start).Milliseconds()), "ms")
}

func BenchmarkStartWithSequential(b *testing.B) {
	start := time.Now()

	scanner := SequentialScanner{}
	for n := 0; n < b.N; n++ {
		scanner.StartScan("127.0.0.1", 1, 65535, 500*time.Millisecond)
	}
	b.ReportMetric(float64(time.Since(start).Milliseconds()), "ms")
}
