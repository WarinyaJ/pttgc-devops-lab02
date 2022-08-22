package blog

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	s              Service
}

func NewInstrucmentingMiddleware(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingMiddleware{
		requestCount:   counter,
		requestLatency: latency,
		s:              s,
	}
}

func (m *instrumentingMiddleware) CreateBlog(ctx context.Context, b Blog) (string, error) {
	defer func(begin time.Time) {
		m.requestCount.With("method", "create_blog").Add(1)
		m.requestLatency.With("method", "create_blog").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return m.s.CreateBlog(ctx, b)
}

func (m *instrumentingMiddleware) GetBlog(ctx context.Context, id string) (*Blog, error) {
	defer func(begin time.Time) {
		m.requestCount.With("method", "get_blog").Add(1)
		m.requestLatency.With("method", "get_blog").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return m.s.GetBlog(ctx, id)
}

func (m *instrumentingMiddleware) ListBlogs(ctx context.Context) ([]Blog, error) {
	defer func(begin time.Time) {
		m.requestCount.With("method", "list_blogs").Add(1)
		m.requestLatency.With("method", "list_blogs").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return m.s.ListBlogs(ctx)
}

func (m *instrumentingMiddleware) PublishBlog(ctx context.Context, id string) (*Blog, error) {
	defer func(begin time.Time) {
		m.requestCount.With("method", "publis_blog").Add(1)
		m.requestLatency.With("method", "publish_blog").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return m.s.PublishBlog(ctx, id)
}
