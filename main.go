package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	_url    = flag.String("url", "https://example.com", "Endpoint URL to request")
	_method = flag.String("method", "GET", "HTTP request method")
	_data   = flag.String("data", "", "Raw body data as string")
	_token  = flag.String("token", "", "Authorization bearer token")
	_user   = flag.String("user", "", "Basic Auth username")
	_pass   = flag.String("pass", "", "Basic Auth password")
	_h      = flag.String("h", "", "List of headers to send in the in the following format: Header1:Value1;Header2:Value2;HeaderN:ValueN")
	_n      = flag.Int("n", 1, "Amount of iterations")
	_c      = flag.Int("c", 1, "Concurrent workers")
	_v      = flag.Int("v", 1, "Verbosity level [0,1,2,3]")
)

func main() {
	// Parse config params
	flag.Parse()

	print(1, `Starting goma with the following configuration:
- HTTP method: %s
- URL endpoint: %s
- Data: %s
- Bearer token: %s
- BasicAuth: %s:%s
- Headers: %s
- Amount of requests to send: %d
- Concurrent request workers: %d
- Verbosity: %d
`, *_method, *_url, *_data, *_token, *_user, *_pass, *_h, *_n, *_c, *_v)

	// Init sync
	var wg sync.WaitGroup
	ch := make(chan int)

	// Create workers
	for i := 1; i <= *_c; i++ {
		print(2, "Starting worker #%d", i)

		wg.Add(1)
		go worker(ch, i, &wg)
	}

	// Run
	for i := 1; i <= *_n; i++ {
		print(2, "Running request #%d", i)
		ch <- i
	}

	// End process
	close(ch)

	// Wait for last requests to finish
	wg.Wait()
}

func worker(ch chan int, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	print(3, "Init Worker #%d", workerID)
	for i := range ch {
		print(2, "Sending request #%d from Worker %d", i, workerID)

		// Buffer data
		var data *bytes.Buffer
		if _data != nil {
			data = bytes.NewBufferString(*_data)
		}

		// Create request
		req, err := http.NewRequest(*_method, *_url, data)
		if err != nil {
			print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
		}

		// Add token/Auth
		if *_token != "" {
			req.Header.Add("Authorization", "bearer "+*_token)
		}
		if *_user != "" && *_pass != "" {
			req.SetBasicAuth(*_user, *_pass)
		}

		// Add headers
		if *_h != "" {
			headers := strings.Split(*_h, ";")
			for _, h := range headers {
				parts := strings.Split(h, ":")
				if len(parts) > 1 {
					req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[0]))
				}
			}
		}

		// Request
		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		end := time.Now()
		if err != nil {
			print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
		}
		print(3, "R#%d W#%d RESPONSE: %v \n%s", i, workerID, resp.StatusCode, string(body))

		delta := end.Sub(start)
		print(1, "Request #%d took %+v and returned %d", i, delta, resp.StatusCode)
	}
}

func print(level int, line string, vars ...interface{}) {
	if *_v >= level {
		log.Printf(line, vars...)
	}
}
