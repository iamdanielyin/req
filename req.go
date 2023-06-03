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
	GET(dst any, headers ...map[string]string) error
	POST(body, dst any, headers ...map[string]string) error
	PATCH(body, dst any, headers ...map[string]string) error
	PUT(body, dst any, headers ...map[string]string) error
	DELETE(dst any, headers ...map[string]string) error
	CALL(method string, body, dst any, headers ...map[string]string) error
}

type urlCaller struct {
	url     string
	body    any
	headers map[string]string
}

func (receiver *urlCaller) URL() string {
	return receiver.url
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

func (receiver *urlCaller) CALL(method string, body, dst any, headers ...map[string]string) error {
	if body != nil {
		receiver.body = body
	}
	if len(headers) > 0 {
		receiver.headers = headers[0]
	}

	var payload io.Reader
	if receiver.body != nil {
		if data, err := json.Marshal(receiver.body); err != nil {
			return err
		} else {
			payload = strings.NewReader(string(data))
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, receiver.url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if len(receiver.headers) > 0 {
		for k, v := range receiver.headers {
			req.Header.Add(k, v)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resp, &dst); err != nil {
		return err
	}
	return nil
}

func URL(format string, values ...any) *urlCaller {
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
