package req

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Caller interface {
	URL() string
	Headers() map[string]string
	AddHeader(k, v string) Caller
	DelHeader(k string) Caller
	SetHeaders(headers map[string]string) Caller
	SetBody(body any) Caller
	SetMethod(method string) Caller
	Do() error
	GET(dst any, headers ...map[string]string) error
	POST(body, dst any, headers ...map[string]string) error
	PATCH(body, dst any, headers ...map[string]string) error
	PUT(body, dst any, headers ...map[string]string) error
	DELETE(dst any, headers ...map[string]string) error
	CALL(method string, body, dst any, headers ...map[string]string) error
}

type urlCaller struct {
	url     string
	method  string
	body    any
	dst     any
	headers map[string]string
}

func (receiver *urlCaller) URL() string {
	return receiver.url
}

func (receiver *urlCaller) SetMethod(method string) Caller {
	receiver.method = method
	return receiver
}

func (receiver *urlCaller) Do() error {
	return receiver.CALL(receiver.method, receiver.body, receiver.dst, receiver.headers)
}

func (receiver *urlCaller) Headers() map[string]string {
	return receiver.headers
}

func (receiver *urlCaller) AddHeader(k, v string) Caller {
	if receiver.headers == nil {
		receiver.headers = make(map[string]string)
	}
	receiver.headers[k] = v
	return receiver
}

func (receiver *urlCaller) DelHeader(k string) Caller {
	if receiver.headers != nil {
		delete(receiver.headers, k)
	}
	return receiver
}

func (receiver *urlCaller) SetHeaders(headers map[string]string) Caller {
	receiver.headers = headers
	return receiver
}

func (receiver *urlCaller) SetBody(body any) Caller {
	receiver.body = body
	return receiver
}

func (receiver *urlCaller) GET(dst any, headers ...map[string]string) error {
	return receiver.CALL(http.MethodGet, nil, dst, headers...)
}

func (receiver *urlCaller) POST(body, dst any, headers ...map[string]string) error {
	return receiver.CALL(http.MethodPost, body, dst, headers...)
}

func (receiver *urlCaller) PATCH(body, dst any, headers ...map[string]string) error {
	return receiver.CALL(http.MethodPatch, body, dst, headers...)
}

func (receiver *urlCaller) PUT(body, dst any, headers ...map[string]string) error {
	return receiver.CALL(http.MethodPut, body, dst, headers...)
}

func (receiver *urlCaller) DELETE(dst any, headers ...map[string]string) error {
	return receiver.CALL(http.MethodDelete, nil, dst, headers...)
}

func (receiver *urlCaller) Raw(method string, body any, headers ...map[string]string) (*http.Response, error) {
	if body != nil {
		receiver.body = body
	}
	if len(headers) > 0 {
		receiver.headers = headers[0]
	}

	var payload io.Reader
	if receiver.body != nil {
		if data, err := json.Marshal(receiver.body); err != nil {
			return nil, err
		} else {
			payload = strings.NewReader(string(data))
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, receiver.url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if len(receiver.headers) > 0 {
		for k, v := range receiver.headers {
			req.Header.Add(k, v)
		}
	}

	return client.Do(req)
}

func (receiver *urlCaller) CALL(method string, body, dst any, headers ...map[string]string) error {
	res, err := receiver.Raw(method, body, headers...)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if dst != nil {
		resp, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(resp, dst); err != nil {
			return err
		}
	}
	return nil
}

func URL(format string, values ...any) Caller {
	return &urlCaller{
		url: fmt.Sprintf(format, values...),
	}
}

func GET(url string, dst any, headers ...map[string]string) error {
	return URL(url).GET(dst, headers...)
}

func POST(url string, body, dst any, headers ...map[string]string) error {
	return URL(url).POST(body, dst, headers...)
}

func PATCH(url string, body, dst any, headers ...map[string]string) error {
	return URL(url).PATCH(body, dst, headers...)
}

func PUT(url string, body, dst any, headers ...map[string]string) error {
	return URL(url).PUT(body, dst, headers...)
}

func DELETE(url string, dst any, headers ...map[string]string) error {
	return URL(url).DELETE(dst, headers...)
}

func DELETEWithBody(url string, body, dst any, headers ...map[string]string) error {
	return URL(url).SetBody(body).DELETE(dst, headers...)
}

func CALL(method, url string, body, dst any, headers ...map[string]string) error {
	return URL(url).CALL(method, body, dst, headers...)
}

func NewGET(url string, dst any, headers ...map[string]string) Caller {
	return NewCALL(http.MethodGet, url, nil, dst, headers...)
}

func NewPOST(url string, body, dst any, headers ...map[string]string) Caller {
	return NewCALL(http.MethodPost, url, body, dst, headers...)
}

func NewPATCH(url string, body, dst any, headers ...map[string]string) Caller {
	return NewCALL(http.MethodPatch, url, body, dst, headers...)
}

func NewPUT(url string, body, dst any, headers ...map[string]string) Caller {
	return NewCALL(http.MethodPut, url, body, dst, headers...)
}

func NewDELETE(url string, dst any, headers ...map[string]string) Caller {
	return NewCALL(http.MethodDelete, url, nil, dst, headers...)
}

func NewDELETEWithBody(url string, body, dst any, headers ...map[string]string) Caller {
	return NewCALL(http.MethodDelete, url, body, dst, headers...)
}

func NewCALL(method, url string, body, dst any, headers ...map[string]string) Caller {
	return &urlCaller{
		url:     url,
		method:  method,
		body:    body,
		dst:     dst,
		headers: getHeaders(headers),
	}
}

func getHeaders(headers []map[string]string) map[string]string {
	if len(headers) > 0 {
		return headers[0]
	}
	return nil
}
