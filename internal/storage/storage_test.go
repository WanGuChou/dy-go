package storage

import (
	"testing"
	"time"
)

func TestStorage_Add(t *testing.T) {
	s := NewStorage()

	req := &CapturedRequest{
		Timestamp:  time.Now(),
		Method:     "GET",
		URL:        "https://example.com/api/users",
		Host:       "example.com",
		Path:       "/api/users",
		Protocol:   "HTTPS",
		StatusCode: 200,
	}

	s.Add(req)

	if s.Count() != 1 {
		t.Errorf("Expected count 1, got %d", s.Count())
	}

	if req.ID != 1 {
		t.Errorf("Expected ID 1, got %d", req.ID)
	}
}

func TestStorage_Search(t *testing.T) {
	s := NewStorage()

	// 添加测试数据
	requests := []*CapturedRequest{
		{
			Method:   "GET",
			URL:      "https://api.example.com/users",
			Host:     "api.example.com",
			Path:     "/users",
			Protocol: "HTTPS",
		},
		{
			Method:   "POST",
			URL:      "https://api.example.com/posts",
			Host:     "api.example.com",
			Path:     "/posts",
			Protocol: "HTTPS",
		},
		{
			Method:   "GET",
			URL:      "https://cdn.example.com/image.png",
			Host:     "cdn.example.com",
			Path:     "/image.png",
			Protocol: "HTTPS",
		},
	}

	for _, req := range requests {
		s.Add(req)
	}

	// 测试搜索
	tests := []struct {
		keyword  string
		expected int
	}{
		{"users", 1},
		{"api", 2},
		{"cdn", 1},
		{"GET", 2},
		{"POST", 1},
		{"example.com", 3},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		results := s.Search(tt.keyword)
		if len(results) != tt.expected {
			t.Errorf("Search(%q): expected %d results, got %d", tt.keyword, tt.expected, len(results))
		}
	}
}

func TestStorage_Clear(t *testing.T) {
	s := NewStorage()

	// 添加一些请求
	for i := 0; i < 5; i++ {
		s.Add(&CapturedRequest{
			Method: "GET",
			URL:    "https://example.com",
		})
	}

	if s.Count() != 5 {
		t.Errorf("Expected count 5, got %d", s.Count())
	}

	s.Clear()

	if s.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", s.Count())
	}
}

func TestStorage_MemoryLimit(t *testing.T) {
	s := NewStorage()

	// 添加超过10000条请求
	for i := 0; i < 11000; i++ {
		s.Add(&CapturedRequest{
			Method: "GET",
			URL:    "https://example.com",
		})
	}

	// 应该只保留最新的10000条
	if s.Count() != 10000 {
		t.Errorf("Expected count 10000, got %d", s.Count())
	}
}
