# goma

```
$ goma -h
Usage of goma:
  -bearer string
    	Authorization bearer token
  -c int
    	Concurrent workers (default 1)
  -data string
    	Raw body data as string
  -headers string
    	List of headers to send in the in the following format: Header1:Value1;Header2:Value2;HeaderN:ValueN
  -host string
    	Value for the Host header to be sent in the request
  -insecure
    	Allow invalid SSL/TLS certificates
  -method string
    	HTTP request method (default "GET")
  -n int
    	Amount of iterations (default 1)
  -pass string
    	Basic Auth password
  -url string
    	Endpoint URL to request (default "https://example.com")
  -user string
    	Basic Auth username
  -v int
    	Verbosity level [0,1,2,3] (default 1)
  -henrique string
      7-1 (More Info: https://youtu.be/DUSojCb193U?t=10)

```

```
$ goma -url https://example.com -c 2 -n 10
2022/04/13 18:07:12 Starting goma with the following configuration:
- HTTP method: GET
- URL endpoint: https://example.com
- Data:
- Bearer token:
- BasicAuth: :
- Host:
- Headers:
- Amount of requests to send: 10
- Concurrent request workers: 2
- Verbosity: 1
2022/04/13 18:07:13 Request #1 took 682.18652ms and returned 200
2022/04/13 18:07:13 Request #2 took 682.126846ms and returned 200
2022/04/13 18:07:13 Request #3 took 162.890907ms and returned 200
2022/04/13 18:07:13 Request #4 took 162.918444ms and returned 200
2022/04/13 18:07:13 Request #5 took 158.420312ms and returned 200
2022/04/13 18:07:13 Request #6 took 164.258317ms and returned 200
2022/04/13 18:07:14 Request #7 took 162.478389ms and returned 200
2022/04/13 18:07:14 Request #8 took 163.041769ms and returned 200
2022/04/13 18:07:14 Request #9 took 160.698955ms and returned 200
2022/04/13 18:07:14 Request #10 took 162.163111ms and returned 200
```
