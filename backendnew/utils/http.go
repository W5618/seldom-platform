package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient HTTP客户端配置
type HTTPClient struct {
	Client  *http.Client
	BaseURL string
	Headers map[string]string
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
		BaseURL: baseURL,
		Headers: make(map[string]string),
	}
}

// SetHeader 设置请求头
func (c *HTTPClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// SetAuthToken 设置认证令牌
func (c *HTTPClient) SetAuthToken(token string) {
	c.SetHeader("Authorization", "Bearer "+token)
}

// Get 发送GET请求
func (c *HTTPClient) Get(endpoint string, params map[string]string) (*http.Response, error) {
	fullURL := c.buildURL(endpoint, params)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.Client.Do(req)
}

// Post 发送POST请求
func (c *HTTPClient) Post(endpoint string, data interface{}) (*http.Response, error) {
	fullURL := c.buildURL(endpoint, nil)
	
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	c.setHeaders(req)
	return c.Client.Do(req)
}

// Put 发送PUT请求
func (c *HTTPClient) Put(endpoint string, data interface{}) (*http.Response, error) {
	fullURL := c.buildURL(endpoint, nil)
	
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("PUT", fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	c.setHeaders(req)
	return c.Client.Do(req)
}

// Delete 发送DELETE请求
func (c *HTTPClient) Delete(endpoint string) (*http.Response, error) {
	fullURL := c.buildURL(endpoint, nil)
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.Client.Do(req)
}

// PostForm 发送表单POST请求
func (c *HTTPClient) PostForm(endpoint string, data url.Values) (*http.Response, error) {
	fullURL := c.buildURL(endpoint, nil)
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.setHeaders(req)
	return c.Client.Do(req)
}

// buildURL 构建完整URL
func (c *HTTPClient) buildURL(endpoint string, params map[string]string) string {
	fullURL := c.BaseURL + endpoint
	
	if len(params) > 0 {
		values := url.Values{}
		for key, value := range params {
			values.Add(key, value)
		}
		fullURL += "?" + values.Encode()
	}
	
	return fullURL
}

// setHeaders 设置请求头
func (c *HTTPClient) setHeaders(req *http.Request) {
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
}

// ParseJSONResponse 解析JSON响应
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, v)
}

// GetResponseBody 获取响应体内容
func GetResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// IsSuccessStatusCode 检查是否为成功状态码
func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// GetClientIP 获取客户端IP地址
func GetClientIP(req *http.Request) string {
	// 检查X-Forwarded-For头
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// 检查X-Real-IP头
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 使用RemoteAddr
	ip := req.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	
	return ip
}

// GetUserAgent 获取用户代理
func GetUserAgent(req *http.Request) string {
	return req.Header.Get("User-Agent")
}

// SetCORSHeaders 设置CORS头
func SetCORSHeaders(w http.ResponseWriter, allowedOrigins []string) {
	origin := "*"
	if len(allowedOrigins) > 0 {
		origin = strings.Join(allowedOrigins, ",")
	}
	
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}