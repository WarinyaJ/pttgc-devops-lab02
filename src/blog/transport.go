package blog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	ErrBadRouting = errors.New("Bad Routing")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoint(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/blogs/v1/blogs").Handler(httptransport.NewServer(
		e.CreateBlogEndpoint,
		decodePostBlogRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/blogs/v1/blogs").Handler(httptransport.NewServer(
		e.ListBlogEndpoint,
		decodeListBlogRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/blogs/v1/blogs/{id}").Handler(httptransport.NewServer(
		e.GetBlogEndpoint,
		decodeGetBlogRequest,
		encodeResponse,
		options...,
	))

	r.Methods("PUT").Path("/blogs/v1/blogs/{id}/published").Handler(httptransport.NewServer(
		e.PublishBlogEndpoint,
		decodePublishBlogRequest,
		encodeResponse,
		options...,
	))

	return r

}

func decodePublishBlogRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return publishBlogRequest{ID: id}, nil
}

func decodeGetBlogRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return getBlogRequest{ID: id}, nil
}

func decodeListBlogRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return listBlogRequest{}, nil
}

func decodePostBlogRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createBlogRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	data := response.(accessor)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(data.data())
}

type errorer interface {
	error() error
}

type accessor interface {
	data() interface{}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrBlogNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
