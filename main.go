package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	url         string
	numReqs     int
	concurrency int
)

func init() {
	flag.StringVar(&url, "url", "", "URL a ser testada")
	flag.IntVar(&numReqs, "requests", 1, "Número total de requisições")
	flag.IntVar(&concurrency, "concurrency", 1, "Número de requisições simultâneas")
}

type result struct {
	duration   time.Duration
	statusCode int
}

func makeRequest(wg *sync.WaitGroup, results chan<- result) {
	defer wg.Done()

	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		results <- result{duration: 0, statusCode: res.StatusCode}
		return
	}

	defer res.Body.Close()
	duration := time.Since(start)
	results <- result{duration: duration, statusCode: res.StatusCode}
}

func main() {
	flag.Parse()

	if url == "" {
		fmt.Println("É mandatório informar uma URL")
		return
	}

	var wg sync.WaitGroup
	results := make(chan result, numReqs)
	startTime := time.Now()

	for i := 0; i < numReqs; i++ {
		wg.Add(1)
		go makeRequest(&wg, results)

		if (i+1)%concurrency == 0 {
			wg.Wait()
		}
	}

	wg.Wait()
	close(results)

	totalDuration := time.Since(startTime)
	successfulRequests := 0
	statusCounts := make(map[int]int)

	var totalResponseTime time.Duration

	for result := range results {
		if result.statusCode == 200 {
			successfulRequests++
		}

		if result.statusCode != 0 {
			statusCounts[result.statusCode]++
		}

		totalResponseTime += result.duration
	}

	fmt.Printf("Resultado:\n")
	fmt.Printf("Tempo total: %v\n", totalDuration)
	fmt.Printf("Requests totais: %v\n", numReqs)
	fmt.Printf("Requests bem sucedidos (200): %v\n", successfulRequests)
	fmt.Printf("Outros status: %v\n", successfulRequests)
	for statusCode, count := range statusCounts {
		if statusCode != 200 {
			fmt.Printf("\t [Status]: quantidade\n")
			fmt.Printf("\t [%d]: %d\n", statusCode, count)
		}
	}
}
