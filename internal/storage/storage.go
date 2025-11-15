package storage

import (
	"strings"
	"sync"
	"time"
)

// CapturedRequest 捕获的请求数据
type CapturedRequest struct {
	ID              int64
	Timestamp       time.Time
	Method          string
	URL             string
	Host            string
	Path            string
	Protocol        string // HTTP, HTTPS, WebSocket
	StatusCode      int
	RequestHeaders  map[string][]string
	RequestBody     []byte
	ResponseHeaders map[string][]string
	ResponseBody    []byte
	Duration        time.Duration
}

// Storage 请求存储管理器
type Storage struct {
	requests []*CapturedRequest
	mu       sync.RWMutex
	nextID   int64
}

// NewStorage 创建存储管理器
func NewStorage() *Storage {
	return &Storage{
		requests: make([]*CapturedRequest, 0, 1000),
		nextID:   1,
	}
}

// Add 添加捕获的请求
func (s *Storage) Add(req *CapturedRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req.ID = s.nextID
	s.nextID++

	s.requests = append(s.requests, req)

	// 限制内存中的请求数量（保留最新的10000条）
	if len(s.requests) > 10000 {
		s.requests = s.requests[len(s.requests)-10000:]
	}
}

// GetAll 获取所有请求
func (s *Storage) GetAll() []*CapturedRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*CapturedRequest, len(s.requests))
	copy(result, s.requests)
	return result
}

// GetByID 根据ID获取请求
func (s *Storage) GetByID(id int64) *CapturedRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, req := range s.requests {
		if req.ID == id {
			return req
		}
	}
	return nil
}

// Search 搜索请求
func (s *Storage) Search(keyword string) []*CapturedRequest {
	if keyword == "" {
		return s.GetAll()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	keyword = strings.ToLower(keyword)
	result := make([]*CapturedRequest, 0)

	for _, req := range s.requests {
		if s.matchRequest(req, keyword) {
			result = append(result, req)
		}
	}

	return result
}

// matchRequest 检查请求是否匹配关键词
func (s *Storage) matchRequest(req *CapturedRequest, keyword string) bool {
	// 检查URL
	if strings.Contains(strings.ToLower(req.URL), keyword) {
		return true
	}

	// 检查Host
	if strings.Contains(strings.ToLower(req.Host), keyword) {
		return true
	}

	// 检查Path
	if strings.Contains(strings.ToLower(req.Path), keyword) {
		return true
	}

	// 检查Method
	if strings.Contains(strings.ToLower(req.Method), keyword) {
		return true
	}

	// 检查请求头
	for key, values := range req.RequestHeaders {
		if strings.Contains(strings.ToLower(key), keyword) {
			return true
		}
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), keyword) {
				return true
			}
		}
	}

	// 检查响应头
	for key, values := range req.ResponseHeaders {
		if strings.Contains(strings.ToLower(key), keyword) {
			return true
		}
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), keyword) {
				return true
			}
		}
	}

	// 检查请求体
	if strings.Contains(strings.ToLower(string(req.RequestBody)), keyword) {
		return true
	}

	// 检查响应体
	if strings.Contains(strings.ToLower(string(req.ResponseBody)), keyword) {
		return true
	}

	return false
}

// Clear 清空所有请求
func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests = make([]*CapturedRequest, 0, 1000)
}

// Count 获取请求总数
func (s *Storage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.requests)
}
