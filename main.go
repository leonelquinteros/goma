package main

import (
	"crypto/tls"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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

// Runner manages the execution of the benchmark
type Runner struct {
	Config *Config
	Client *http.Client
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

	// Create workers
	for i := 1; i <= r.Config.C; i++ {
		r.print(2, "Starting worker #%d", i)

		wg.Add(1)
		go r.worker(ch, i, &wg)
	}

	// Run
	for i := 1; i <= r.Config.N; i++ {
		r.print(2, "Running request #%d", i)
		ch <- i
	}

	// End process
	close(ch)

	// Wait for last requests to finish
	wg.Wait()
}

func (r *Runner) worker(ch chan int, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	r.print(3, "Init Worker #%d", workerID)
	for i := range ch {
		r.print(2, "Sending request #%d from Worker %d", i, workerID)

		req, err := r.Config.buildRequest()
		if err != nil {
			r.print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
			continue
		}

		// Request
		start := time.Now()
		resp, err := r.Client.Do(req)
		if err != nil {
			r.print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
			continue
		}
		end := time.Now()

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			r.print(0, "R#%d W#%d ERROR: %v", i, workerID, err.Error())
			continue
		}
		r.print(3, "R#%d W#%d RESPONSE: %v \n%s", i, workerID, resp.StatusCode, string(body))

		delta := end.Sub(start)
		r.print(1, "Request #%d took %+v and returned %d", i, delta, resp.StatusCode)
	}
}

func (r *Runner) print(level int, line string, vars ...interface{}) {
	if r.Config.V >= level {
		log.Printf(line, vars...)
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
