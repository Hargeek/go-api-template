package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	errort "go-api-template/common/error"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 使用 testify/mock 定义 mock
// 只在测试文件中定义

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) QueryWeather(ctx context.Context, city string) (string, *errort.ApiError) {
	args := m.Called(ctx, city)
	var apiErr *errort.ApiError
	if e := args.Get(1); e != nil {
		apiErr = e.(*errort.ApiError)
	}
	return args.String(0), apiErr
}

func TestWeatherController_QueryWeather(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockWeatherService)
	weatherController := NewWeatherController(mockService)
	router := gin.New()
	router.GET("/weather", weatherController.QueryWeather)

	// 1. 正常返回
	mockService.On("QueryWeather", mock.Anything, "Beijing").Return("Beijing: 晴，25°C", nil)
	req1, _ := http.NewRequest("GET", "/weather?city=Beijing", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Contains(t, w1.Body.String(), "Beijing: 晴")

	// 2. 缺少city参数
	req2, _ := http.NewRequest("GET", "/weather", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	// 3. 查询失败
	mockService.On("QueryWeather", mock.Anything, "Unknown").
		Return("", errort.NewApiError(errort.GeneralError, fmt.Errorf(errort.MsgWeatherQueryFailed, errors.New("not found"))))
	req3, _ := http.NewRequest("GET", "/weather?city=Unknown", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusInternalServerError, w3.Code)
	assert.Contains(t, w3.Body.String(), "查询天气失败")

	mockService.AssertExpectations(t)
}
