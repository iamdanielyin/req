package req

import (
	"errors"
	"net/http"
	"testing"
)

func Test_urlCaller_URL(t *testing.T) {
	if URL("a=%d", 1).URL() != "a=1" {
		t.Fatal("test failed")
	}
}

func Test_urlCaller_Headers(t *testing.T) {
	c := &urlCaller{}
	c.AddHeader("a", "1")
	c.AddHeader("b", "2")
	c.DelHeader("b")
	if c.Headers()["a"] != "1" || c.Headers()["b"] != "" {
		t.Fatal("test failed")
	}
	c.SetHeaders(nil)
	if c.Headers() != nil {
		t.Fatal("test failed")
	}
}

func TestGET(t *testing.T) {
	var posts []struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := GET("https://jsonplaceholder.typicode.com/posts", &posts); err != nil {
		t.Fatal(err)
	}
	t.Log(len(posts))
}

func TestPOST(t *testing.T) {
	var post struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := POST("https://jsonplaceholder.typicode.com/posts", map[string]any{"Title": "Hello"}, &post, map[string]string{"UID": "1"}); err != nil {
		t.Fatal(err)
	}
	if post.Id < 0 || post.Title != "Hello" {
		t.Fatal(errors.New("request failed"))
	}
	t.Log(post)
}

func TestPATCH(t *testing.T) {
	var post struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := PATCH("https://jsonplaceholder.typicode.com/posts/1", map[string]any{"Title": "Hello"}, &post, map[string]string{"UID": "1"}); err != nil {
		t.Fatal(err)
	}
	if post.Id != 1 || post.Title != "Hello" {
		t.Fatal(errors.New("request failed"))
	}
	t.Log(post)
}

func TestPUT(t *testing.T) {
	var post struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := PUT("https://jsonplaceholder.typicode.com/posts/1", map[string]any{"Title": "Hello"}, &post, map[string]string{"UID": "1"}); err != nil {
		t.Fatal(err)
	}
	if post.Id != 1 || post.Title != "Hello" {
		t.Fatal(errors.New("request failed"))
	}
	t.Log(post)
}

func TestDELETE(t *testing.T) {
	var ret map[string]string
	if err := DELETE("https://jsonplaceholder.typicode.com/posts/1", &ret, map[string]string{"UID": "1"}); err != nil {
		t.Fatal(err)
	}
	t.Log(ret)

	if err := DELETE("https://reqres.in/api/users/2", &ret); err == nil {
		t.Fatal("test failed")
	}
}

func TestDELETEWithBody(t *testing.T) {
	var ret map[string]string
	if err := DELETEWithBody("https://jsonplaceholder.typicode.com/posts/1", map[string]any{"Title": "Hello"}, &ret, map[string]string{"UID": "1"}); err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
}

func TestCALL(t *testing.T) {
	var post struct {
		UserId int    `json:"userId"`
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := CALL(http.MethodPost, "https://jsonplaceholder.typicode.com/posts", map[string]any{"Title": "Hello"}, &post, map[string]string{"UID": "1"}); err != nil {
		t.Fatal(err)
	}
	if post.Id < 0 || post.Title != "Hello" {
		t.Fatal(errors.New("request failed"))
	}
	t.Log(post)

	if err := CALL(" ", "https://54d483ff-bb2b-a4f7-26e4-6b8b11c09ed1.com/api/abc", nil, &post); err == nil {
		t.Fatal("test failed")
	}
}
