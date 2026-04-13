package main

import (
	"io/ioutil"
	"testing"
)

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		expectedMethod string
		expectedURL    string
		expectedBody   string
		expectedAuth   string
		expectedHost   string
		expectedHeader map[string]string
	}{
		{
			name: "Basic GET",
			config: Config{
				Method: "GET",
				URL:    "https://example.com",
			},
			expectedMethod: "GET",
			expectedURL:    "https://example.com",
		},
		{
			name: "POST with Data",
			config: Config{
				Method: "POST",
				URL:    "https://example.com",
				Data:   "test-data",
			},
			expectedMethod: "POST",
			expectedURL:    "https://example.com",
			expectedBody:   "test-data",
		},
		{
			name: "Bearer Auth",
			config: Config{
				Method: "GET",
				URL:    "https://example.com",
				Bearer: "token123",
			},
			expectedMethod: "GET",
			expectedURL:    "https://example.com",
			expectedAuth:   "bearer token123",
		},
		{
			name: "Basic Auth",
			config: Config{
				Method: "GET",
				URL:    "https://example.com",
				User:   "user",
				Pass:   "pass",
			},
			expectedMethod: "GET",
			expectedURL:    "https://example.com",
			expectedHeader: map[string]string{"Authorization": "Basic dXNlcjpwYXNz"},
		},
		{
			name: "Custom Host",
			config: Config{
				Method: "GET",
				URL:    "https://example.com",
				Host:   "custom.host",
			},
			expectedMethod: "GET",
			expectedURL:    "https://example.com",
			expectedHost:   "custom.host",
		},
		{
			name: "Custom Headers",
			config: Config{
				Method:  "GET",
				URL:     "https://example.com",
				Headers: "X-Header1:Value1;X-Header2:Value2",
			},
			expectedMethod: "GET",
			expectedURL:    "https://example.com",
			expectedHeader: map[string]string{
				"X-Header1": "Value1",
				"X-Header2": "Value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.config.buildRequest()
			if err != nil {
				t.Fatalf("buildRequest() error = %v", err)
			}

			if req.Method != tt.expectedMethod {
				t.Errorf("Method = %v, want %v", req.Method, tt.expectedMethod)
			}

			if req.URL.String() != tt.expectedURL {
				t.Errorf("URL = %v, want %v", req.URL.String(), tt.expectedURL)
			}

			if tt.expectedBody != "" {
				body, _ := ioutil.ReadAll(req.Body)
				if string(body) != tt.expectedBody {
					t.Errorf("Body = %v, want %v", string(body), tt.expectedBody)
				}
			}

			if tt.expectedAuth != "" {
				auth := req.Header.Get("Authorization")
				if auth != tt.expectedAuth {
					t.Errorf("Auth = %v, want %v", auth, tt.expectedAuth)
				}
			}

			if tt.expectedHost != "" {
				if req.Host != tt.expectedHost {
					t.Errorf("Host = %v, want %v", req.Host, tt.expectedHost)
				}
			}

			for k, v := range tt.expectedHeader {
				if req.Header.Get(k) != v {
					t.Errorf("Header[%s] = %v, want %v", k, req.Header.Get(k), v)
				}
			}
		})
	}
}
