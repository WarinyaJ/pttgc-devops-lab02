package blog

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateBlogEndpoint  endpoint.Endpoint
	GetBlogEndpoint     endpoint.Endpoint
	ListBlogEndpoint    endpoint.Endpoint
	PublishBlogEndpoint endpoint.Endpoint
}

func MakeServerEndpoint(s Service) Endpoints {
	return Endpoints{
		CreateBlogEndpoint:  makeCreateBlogEndpoint(s),
		GetBlogEndpoint:     makeGetBlogEndpoint(s),
		ListBlogEndpoint:    makeListBlogsEndpoint(s),
		PublishBlogEndpoint: makePublishBlogEndpoint(s),
	}
}

func makeCreateBlogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createBlogRequest)
		b := Blog{
			Topic:   req.Topic,
			Content: req.Content,
			Author:  req.Author,
			Status:  Draft,
		}
		id, err := s.CreateBlog(ctx, b)
		return createBlogResponse{ID: id, Err: err}, err
	}
}

func makeGetBlogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getBlogRequest)
		b, err := s.GetBlog(ctx, req.ID)
		return getBlogResponse{Blog: *b, Err: err}, err
	}
}

func makeListBlogsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listBlogRequest)
		list, err := s.ListBlogs(ctx)
		return listBlogResponse{Blogs: list, Err: err}, err
	}
}

func makePublishBlogEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(publishBlogRequest)
		b, err := s.PublishBlog(ctx, req.ID)
		if b == nil {
			return publishBlogResponse{Blog: Blog{}, Err: err}, err
		}

		return publishBlogResponse{Blog: *b, Err: err}, err
	}
}

type createBlogRequest struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type createBlogResponse struct {
	ID  string `json:"id,omitempty"`
	Err error  `json:"error,omitemty"`
}

func (r createBlogResponse) error() error      { return r.Err }
func (r createBlogResponse) data() interface{} { return r }

type getBlogRequest struct {
	ID string
}

type getBlogResponse struct {
	Blog `json:"blog,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r getBlogResponse) error() error      { return r.Err }
func (r getBlogResponse) data() interface{} { return r.Blog }

type listBlogRequest struct {
}

type listBlogResponse struct {
	Blogs []Blog `json:"blogs,omitempty"`
	Err   error  `json:"error,omitempty"`
}

func (r listBlogResponse) error() error      { return r.Err }
func (r listBlogResponse) data() interface{} { return r.Blogs }

type publishBlogRequest struct {
	ID string
}

type publishBlogResponse struct {
	Blog `json:"blog,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r publishBlogResponse) error() error      { return r.Err }
func (r publishBlogResponse) data() interface{} { return r.Blog }
