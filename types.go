package main

type LocustStats struct {
	Errors    []LocustError `json:"errors"`
	Stats     []LocustStat  `json:"stats"`
	State     string        `json:"state"`
	TotalRps  float32       `json:"total_rps"`
	FailRatio float32       `json:"fail_ratio"`
	UserCount int           `json:"user_count"`
}

type LocustError struct {
	Error      string `json:"error"`
	Method     string `json:"method"`
	Occurences int    `json:"occurences"`
	Name       string `json:"name"`
}

type LocustStat struct {
	MedianResponseTime int     `json:"median_response_time"`
	MinResponseTime    int     `json:"min_response_time"`
	CurrentRps         float32 `json:"current_rps"`
	Name               string  `json:"name"`
	NumFailures        int     `json:"num_failures"`
	MaxResponseTime    int     `json:"max_response_time"`
	AvgContentLength   int     `json:"avg_content_length"`
	AvgResponseTime    float32 `json:"avg_response_time"`
	Method             string  `json:"method"`
	NumRequests        int     `json:"num_requests"`
}
