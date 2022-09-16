package http_client

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"image-reports/api-gateway/dtos"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	tools "image-reports/api-gateway/helpers/http-client"

	"image-reports/helpers/configs"
)

type HttpClient interface {
	GetFile(ctx context.Context, imageId string) (*http.Response, *dtos.ErrorOutbound)
	SaveFile(ctx context.Context, data reporter_dtos.SaveImage) (*http.Response, *dtos.ErrorOutbound)
}

type httpClient struct {
	config configs.GlobalConfig
}

func NewHttpClient(config configs.GlobalConfig) HttpClient {
	return &httpClient{config: config}
}

func (client *httpClient) GetFile(ctx context.Context, imageId string) (*http.Response, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("GET", tools.Url(client.config.Services.Storage, "/v1/storage/%s", imageId), nil)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[any](ctx, req).Execute(nil)
}

func (client *httpClient) SaveFile(ctx context.Context, data reporter_dtos.SaveImage) (*http.Response, *dtos.ErrorOutbound) {
	filename := data.ImageID
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("image", filename)
	if err != nil {
		bodyWriter.Close()
		return nil, dtos.NewInternalError(err.Error())
	}

	// open file handle
	fh, err := data.Image.Open()
	if err != nil {
		bodyWriter.Close()
		return nil, dtos.NewInternalError(err.Error())
	}
	defer fh.Close()

	//Copy content
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}

	bodyWriter.Boundary()

	bodyWriter.WriteField("imageId", filename)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", tools.Url(client.config.Services.Storage, "/v1/storage"), bodyBuf)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[any](ctx, req).ContentType(contentType).Execute(nil)
}
