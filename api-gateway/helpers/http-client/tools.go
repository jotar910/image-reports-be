package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"image-reports/api-gateway/dtos"

	"image-reports/helpers/configs"
	log "image-reports/helpers/services/logger"
)

var MimeJSON = "text/json"

func Url(config configs.ServiceConfig, path string, args ...any) string {
	return fmt.Sprintf("http://%s:%d%s", config.Host, config.Port, fmt.Sprintf(path, args...))
}

func ToJSON[TValue any](value TValue) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(value); err != nil {
		log.Fatalf("encoding to json: %v", err)
	}
	return buf
}

func ReadJSON[TValue any, TError fmt.Stringer](resp *http.Response) (*TValue, *dtos.ErrorOutbound) {
	if resp.StatusCode == http.StatusOK {
		return readJSON[TValue](resp)
	}
	return nil, ExtractError[TError](resp)
}

func ExtractError[TError fmt.Stringer](resp *http.Response) *dtos.ErrorOutbound {
	if errRes, err := readJSON[TError](resp); err == nil {
		return dtos.NewError(resp.StatusCode, (*errRes).String())
	}
	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		return dtos.NewInternalError(err.Error())
	}
	if n == 0 {
		return dtos.NewError(resp.StatusCode, resp.Status)
	}
	return dtos.NewError(resp.StatusCode, buf.String())
}

func readJSON[T any](resp *http.Response) (*T, *dtos.ErrorOutbound) {
	res := new(T)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(res); err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return res, nil
}

type HttpRequest[T any] struct {
	req   *http.Request
	query url.Values
}

func NewHttpRequest[T any](ctx context.Context, req *http.Request) *HttpRequest[T] {
	hr := &HttpRequest[T]{req, req.URL.Query()}
	return hr.init(ctx)
}

func (hr *HttpRequest[T]) Execute(body any) (*http.Response, *dtos.ErrorOutbound) {
	c := http.Client{}
	hr.req.URL.RawQuery = hr.query.Encode()
	resp, err := c.Do(hr.req)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ExtractError[dtos.ErrorOutbound](resp)
	}
	return resp, nil
}

func (hr *HttpRequest[T]) Do(body any) (*T, *dtos.ErrorOutbound) {
	resp, err := hr.Execute(hr.req)
	if err != nil {
		return nil, err
	}
	result, oerr := ReadJSON[T, dtos.ErrorResponse](resp)
	if oerr != nil {
		return nil, oerr
	}
	return result, nil
}

func (hr *HttpRequest[T]) AddQuery(key string, value string, args ...any) *HttpRequest[T] {
	hr.query.Add(key, fmt.Sprintf(value, args...))
	return hr
}

func (hr *HttpRequest[T]) AddHeader(key string, value string, args ...any) *HttpRequest[T] {
	hr.req.Header.Add(key, fmt.Sprintf(value, args...))
	return hr
}

func (hr *HttpRequest[T]) ContentType(value string) *HttpRequest[T] {
	return hr.AddHeader("Content-type", value)
}

func (hr *HttpRequest[T]) init(ctx context.Context) *HttpRequest[T] {
	if token, ok := ctx.Value("token").(string); ok {
		hr.AddHeader("Authorization", token)
	}
	return hr
}
