package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxResponseBodySize = 1 << 20 // 最多读取 1 MiB，防止异常下游响应占用过多内存。

// Config 定义天气 HTTP 客户端所需的基础设施配置。
type Config struct {
	BaseURL string
	Timeout time.Duration
}

// CurrentWeather 是客户端从天气服务协议中提取的最小结果。
// 该类型只供 Adapter 使用，不作为项目的业务接口或 HTTP 响应结构。
type CurrentWeather struct {
	Description  string
	TemperatureC float64
}

// Client 封装天气服务的地址、HTTP 连接复用和协议解析。
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient 创建可复用的天气客户端。地址或超时不合法时立即返回错误，由装配层决定退出启动。
func NewClient(cfg Config) (*Client, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid weather base URL %q", cfg.BaseURL)
	}
	if cfg.Timeout <= 0 {
		return nil, errors.New("weather client timeout must be greater than zero")
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}, nil
}

// GetCurrent 查询城市当前天气，并将 wttr.in 的响应转换为稳定的客户端结果。
func (c *Client) GetCurrent(ctx context.Context, city string) (*CurrentWeather, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return nil, errors.New("city is required")
	}

	requestURL := c.baseURL + "/" + url.PathEscape(city) + "?format=j1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create weather request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request weather service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("weather service returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodySize+1))
	if err != nil {
		return nil, fmt.Errorf("read weather response: %w", err)
	}
	if len(body) > maxResponseBodySize {
		return nil, errors.New("weather response exceeds size limit")
	}
	var payload currentWeatherResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode weather response: %w", err)
	}
	if len(payload.CurrentCondition) == 0 || len(payload.CurrentCondition[0].WeatherDescription) == 0 {
		return nil, errors.New("weather response does not contain current condition")
	}
	temperature, err := payload.CurrentCondition[0].temperatureCelsius()
	if err != nil {
		return nil, err
	}
	return &CurrentWeather{
		Description:  payload.CurrentCondition[0].WeatherDescription[0].Value,
		TemperatureC: temperature,
	}, nil
}
