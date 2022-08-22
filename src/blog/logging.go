package blog

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	s      Service
}

func NewLoggingMiddleware(logger log.Logger, s Service) Service {
	return &loggingMiddleware{
		logger: logger,
		s:      s,
	}
}

func (l *loggingMiddleware) CreateBlog(ctx context.Context, b Blog) (id string, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"methode", "create_blog",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.CreateBlog(ctx, b)
}

func (l *loggingMiddleware) GetBlog(ctx context.Context, id string) (b *Blog, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "get_blog",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.GetBlog(ctx, id)
}

func (l *loggingMiddleware) ListBlogs(ctx context.Context) (items []Blog, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "list_blogs",
			"count", len(items),
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.ListBlogs(ctx)
}

func (l *loggingMiddleware) PublishBlog(ctx context.Context, id string) (b *Blog, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "publish_blog",
			"id", id,
			"took", time.Since(begin),
			"error", err,
		)
	}(time.Now())
	return l.s.PublishBlog(ctx, id)
}
