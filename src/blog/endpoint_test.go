package blog

import (
	"errors"
	"testing"
)

func TestPublishBlogResponseError(t *testing.T) {
	expected := errors.New("Test")
	p := publishBlogResponse{
		Err: expected,
	}

	result := p.error()
	if result != expected {
		t.Errorf("got %v, expected %v", result, expected)
	}
}

func TestListBlogResponseError(t *testing.T) {
	expected := errors.New("Test")
	l := listBlogResponse{
		Err: expected,
	}

	result := l.error()
	if result != expected {
		t.Errorf("got %v, expected %v", result, expected)
	}
}

func TestGetBlogResponse(t *testing.T) {
	expected := errors.New("Test")
	l := getBlogResponse{
		Err: expected,
	}

	result := l.error()
	if result != expected {
		t.Errorf("got %v, expected %v", result, expected)
	}
}
