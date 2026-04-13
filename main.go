package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// Config holds the configuration for the goma tool
type Config struct {
	URL      string
	Method   string
	Data     string
	Bearer   string
	User     string
	Pass     string
	Host     string
	Headers  string
	Insecure bool
	N        int
	C        int
	V        int
}

// Result holds the information about a single HTTP request result
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

// Runner manages the execution of the benchmark
type Runner struct {
	Config  *Config
	Client  *http.Client
	results []Result
}

// buildRequest creates an http.Request based on the configuration
func (c *Config) buildRequest() (*http.Request, error) {
	var body io.Reader
	if c.Data != "" {
		body = strings.NewReader(c.Data)
	}

	req, err := http.NewRequest(c.Method, c.URL, body)
	if err != nil {
		return nil, err
	}

	// Add token/Auth
	if c.Bearer != "" {
		req.Header.Add("Authorization", "bearer "+c.Bearer)
	}
	if c.User != "" && c.Pass != "" {
		req.SetBasicAuth(c.User, c.Pass)
	}

	// Add host
	if c.Host != "" {
		req.Host = c.Host
	}

	// Add headers
	if c.Headers != "" {
		headers := strings.Split(c.Headers, ";")
		for _, h := range headers {
			parts := strings.Split(h, ":")
			if len(parts) > 1 {
				req.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}

	return req, nil
}

// Run starts the benchmark
func (r *Runner) Run() {
	// Init sync
	var wg sync.WaitGroup
	ch := make(chan int)
	resCh := make(chan Result, r.Config.N)

	// Create workers
	for i := 1; i <= r.Config.C; i++ {
		r.print(2, "Starting worker #%d", i)

		wg.Add(1)
		go r.worker(ch, resCh, i, &wg)
	}

	// Run
	start := time.Now()
	for i := 1; i <= r.Config.N; i++ {
		r.print(2, "Running request #%d", i)
		ch <- i
	}

	// End process
	close(ch)

	// Wait for last requests to finish
	wg.Wait()
	totalDuration := time.Since(start)
	close(resCh)

	// Collect results
	for res := range resCh {
		r.results = append(r.results, res)
	}

	r.printSummary(totalDuration)
}

func (r *Runner) worker(ch chan int, resCh chan Result, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	r.print(3, "Init Worker #%d", workerID)
	for i := range ch {
		r.print(2, "Sending request #%d from Worker %d", i, workerID)

		req, err := r.Config.buildRequest()
		if err != nil {
			r.print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
			resCh <- Result{Error: err}
			continue
		}

		// Request
		reqStart := time.Now()
		resp, err := r.Client.Do(req)
		if err != nil {
			r.print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
			resCh <- Result{Error: err}
			continue
		}
		reqEnd := time.Now()
		delta := reqEnd.Sub(reqStart)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			r.print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
			resCh <- Result{StatusCode: resp.StatusCode, Duration: delta, Error: err}
			continue
		}
		r.print(3, "R#%d W#%d RESPONSE: %v \n%s", i, workerID, resp.StatusCode, string(body))

		r.print(1, "Request #%d took %+v and returned %d", i, delta, resp.StatusCode)
		resCh <- Result{StatusCode: resp.StatusCode, Duration: delta}
	}
}

func (r *Runner) print(level int, line string, vars ...interface{}) {
	if r.Config.V >= level {
		log.Printf(line, vars...)
	}
}

func (r *Runner) printSummary(totalDuration time.Duration) {
	if len(r.results) == 0 {
		fmt.Println("No results to summarize.")
		return
	}

	var (
		success   int
		failed    int
		totalTime time.Duration
		minTime   = r.results[0].Duration
		maxTime   time.Duration
		durations []time.Duration
	)

	statusCodes := make(map[int]int)

	for _, res := range r.results {
		if res.Error != nil {
			failed++
		} else {
			success++
		}

		if res.StatusCode > 0 {
			statusCodes[res.StatusCode]++
		}

		totalTime += res.Duration
		if res.Duration < minTime {
			minTime = res.Duration
		}
		if res.Duration > maxTime {
			maxTime = res.Duration
		}
		durations = append(durations, res.Duration)
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	p50 := durations[len(durations)*50/100]
	p90 := durations[len(durations)*90/100]
	p95 := durations[len(durations)*95/100]
	p99 := durations[len(durations)*99/100]

	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total requests:        %d\n", len(r.results))
	fmt.Printf("  Successful requests:   %d\n", success)
	fmt.Printf("  Failed requests:       %d\n", failed)
	fmt.Printf("  Total time:            %v\n", totalDuration)
	fmt.Printf("  Average request time:  %v\n", totalTime/time.Duration(len(r.results)))
	fmt.Printf("  Min request time:      %v\n", minTime)
	fmt.Printf("  Max request time:      %v\n", maxTime)
	fmt.Printf("  Requests per second:   %.2f\n", float64(len(r.results))/totalDuration.Seconds())

	fmt.Printf("\nLatency Percentiles:\n")
	fmt.Printf("  P50: %v\n", p50)
	fmt.Printf("  P90: %v\n", p90)
	fmt.Printf("  P95: %v\n", p95)
	fmt.Printf("  P99: %v\n", p99)

	fmt.Printf("\nStatus Codes:\n")
	for code, count := range statusCodes {
		fmt.Printf("  [%d] %d responses\n", code, count)
	}
}

func main() {
	var config Config

	// Parse config params
	flag.StringVar(&config.URL, "url", "https://example.com", "Endpoint URL to request")
	flag.StringVar(&config.Method, "method", "GET", "HTTP request method")
	flag.StringVar(&config.Data, "data", "", "Raw body data as string")
	flag.StringVar(&config.Bearer, "bearer", "", "Authorization bearer token")
	flag.StringVar(&config.User, "user", "", "Basic Auth username")
	flag.StringVar(&config.Pass, "pass", "", "Basic Auth password")
	flag.StringVar(&config.Host, "host", "", "Value for the Host header to be sent in the request")
	flag.StringVar(&config.Headers, "headers", "", "List of headers to send in the in the following format: Header1:Value1;Header2:Value2;HeaderN:ValueN")
	flag.BoolVar(&config.Insecure, "insecure", false, "Allow invalid SSL/TLS certificates")
	flag.IntVar(&config.N, "n", 1, "Amount of iterations")
	flag.IntVar(&config.C, "c", 1, "Concurrent workers")
	flag.IntVar(&config.V, "v", 1, "Verbosity level [0,1,2,3]")
	flag.Parse()

	log.Printf(`Starting goma with the following configuration:
- HTTP method: %s
- URL endpoint: %s
- Data: %s
- Bearer token: %s
- BasicAuth: %s:%s
- Host: %s
- Headers: %s
- Amount of requests to send: %d
- Concurrent request workers: %d
- Verbosity: %d
`, config.Method, config.URL, config.Data, config.Bearer, config.User, config.Pass, config.Host, config.Headers, config.N, config.C, config.V)

	// HTTP client config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure},
		},
	}

	runner := &Runner{
		Config: &config,
		Client: client,
	}

	runner.Run()
}
