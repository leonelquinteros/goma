# goma

```
$ goma -h
Usage of goma:
  -c int
        Concurrent workers (default 1)
  -data string
        Raw body data as string
  -h string
        List of headers to send in the in the following format: Header1:Value1;Header2:Value2;HeaderN:ValueN
  -method string
        HTTP request method (default "GET")
  -n int
        Amount of iterations (default 1)
  -pass string
        Basic Auth password
  -token string
        Authorization bearer token
  -url string
        Endpoint URL to request (default "https://example.com")
  -user string
        Basic Auth username
  -v int
        Verbosity level [0,1,2,3] (default 1)

```

```
$ goma -url https://example.com -c 2 -n 10
2020/08/13 11:15:52 Starting goma with the following configuration:
- HTTP method: GET
- URL endpoint: https://example.com
- Data:
- Bearer token:
- Amount of requests to send: 10
- Concurrent request workers: 2
- Verbosity: 1
2020/08/13 11:15:53 Request #2 took 804.459505ms and returned 200
2020/08/13 11:15:53 Request #1 took 804.516926ms and returned 200
2020/08/13 11:15:53 Request #3 took 170.424724ms and returned 200
2020/08/13 11:15:53 Request #4 took 170.366123ms and returned 200
2020/08/13 11:15:53 Request #5 took 148.087707ms and returned 200
2020/08/13 11:15:53 Request #6 took 149.154481ms and returned 200
2020/08/13 11:15:53 Request #7 took 165.530273ms and returned 200
2020/08/13 11:15:53 Request #8 took 171.168688ms and returned 200
2020/08/13 11:15:53 Request #9 took 160.275129ms and returned 200
2020/08/13 11:15:53 Request #10 took 177.637241ms and returned 200
```