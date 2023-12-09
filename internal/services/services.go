package services

import (
	"blastoise/internal/worker"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type RequestResult struct {
	StatusCode int
	Duration   int
}

func (s *Service) Execute(url string, method string, requestsPerSecond int, duration int, channel chan []*RequestResult) {
	// var wg sync.WaitGroup

	pool := worker.NewPromisePool(requestsPerSecond)

	totalRequests := 0
	results := make([]*RequestResult, 0)

	defer pool.Wait()

	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
	defer ticker.Stop()

	end := time.Now().Add(time.Duration(duration) * time.Second)

	for {

		if time.Now().After(end) {
			break
		}

		pool.Add(func() error {

			totalRequests++
			result := request(url, method)

			results = append(results, &result)

			return nil
		})

		<-ticker.C
	}

	// Wait for all tasks to complete before exiting
	pool.Wait()

	channel <- results
}

func request(url string, method string) RequestResult {

	start := time.Now()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {

		fmt.Println(err)

		return RequestResult{
			StatusCode: 500,
			Duration:   int(time.Since(start).Milliseconds()),
		}

	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return RequestResult{
			StatusCode: 500,
			Duration:   int(time.Since(start).Milliseconds()),
		}

	}

	result := RequestResult{
		StatusCode: resp.StatusCode,
		Duration:   int(time.Since(start).Milliseconds()),
	}

	// Release the semaphore

	return result
}

func (s *Service) SummarizeStatusCodes(results []*RequestResult) map[int]int {
	statusCodes := make(map[int]int)

	for _, result := range results {
		statusCodes[result.StatusCode]++
	}

	return statusCodes
}

// P10, P20, P50, P75, P90, P95, P99
func (s *Service) SummarizeReponseTimesInPercentiles(results []*RequestResult) map[string]int {
	percentiles := make(map[string]int)

	// Step 1: Sort the results by duration
	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration < results[j].Duration
	})

	// Step 2: Calculate percentiles
	totalResults := len(results)
	percentiles["P10"] = calculatePercentile(results, 10, totalResults)
	percentiles["P20"] = calculatePercentile(results, 20, totalResults)
	percentiles["P50"] = calculatePercentile(results, 50, totalResults)
	percentiles["P75"] = calculatePercentile(results, 75, totalResults)
	percentiles["P90"] = calculatePercentile(results, 90, totalResults)
	percentiles["P95"] = calculatePercentile(results, 95, totalResults)
	percentiles["P99"] = calculatePercentile(results, 99, totalResults)

	// Step 3: Return the percentiles
	return percentiles
}

func calculatePercentile(results []*RequestResult, percentile int, totalResults int) int {
	if totalResults == 0 {
		return 0
	}

	index := int(math.Ceil(float64(percentile)/100*float64(totalResults))) - 1
	if index < 0 {
		index = 0
	} else if index >= totalResults {
		index = totalResults - 1
	}

	return results[index].Duration
}
