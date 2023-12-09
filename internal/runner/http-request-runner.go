package runner

import (
	"blastoise/internal/structs"
	"blastoise/internal/utils"
	"fmt"
	"net/http"
	"time"
)

type HttpRequestRunner struct {
	ctx *structs.Ctx
}

func NewHttpRequestRunner(ctx *structs.Ctx) *HttpRequestRunner {
	return &HttpRequestRunner{
		ctx: ctx,
	}
}

func (h *HttpRequestRunner) Run() {

	pool := utils.NewWorkerPool(h.ctx.Rps)

	totalRequests := 0
	results := make([]*structs.RequestResult, 0)

	defer pool.Wait()

	ticker := time.NewTicker(time.Second / time.Duration(h.ctx.Rps))
	defer ticker.Stop()

	end := time.Now().Add(time.Duration(h.ctx.Duration) * time.Second)

	for {

		if time.Now().After(end) {
			break
		}

		pool.Add(func() error {

			totalRequests++
			result := request(h.ctx.Url, h.ctx.Method)

			results = append(results, &result)

			return nil
		})

		<-ticker.C
	}

	pool.Wait()

	h.ctx.ResultChan <- results
}

func request(url string, method string) structs.RequestResult {

	start := time.Now()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {

		fmt.Println(err)

		return structs.RequestResult{
			StatusCode: 500,
			Duration:   int(time.Since(start).Milliseconds()),
		}

	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return structs.RequestResult{
			StatusCode: 500,
			Duration:   int(time.Since(start).Milliseconds()),
		}

	}

	result := structs.RequestResult{
		StatusCode: resp.StatusCode,
		Duration:   int(time.Since(start).Milliseconds()),
	}

	return result
}
