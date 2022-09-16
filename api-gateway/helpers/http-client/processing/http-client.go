package http_client

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"image-reports/api-gateway/dtos"
	processing_dtos "image-reports/api-gateway/dtos/processing"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	tools "image-reports/api-gateway/helpers/http-client"

	"image-reports/helpers/configs"
	"image-reports/helpers/utils"
)

type HttpClient interface {
	GetEvaluation(ctx context.Context, id uint) (*processing_dtos.Evaluation, *dtos.ErrorOutbound)
	ProcessImage(ctx context.Context, data reporter_dtos.ProcessImage) (*http.Response, *dtos.ErrorOutbound)
	Search(ctx context.Context, ids []uint) ([]processing_dtos.Evaluation, *dtos.ErrorOutbound)
}

type httpClient struct {
	config configs.GlobalConfig
}

func NewHttpClient(config configs.GlobalConfig) HttpClient {
	return &httpClient{config: config}
}

func (client *httpClient) GetEvaluation(ctx context.Context, id uint) (*processing_dtos.Evaluation, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("GET", tools.Url(client.config.Services.Processing, "/v1/evaluations/%d", id), nil)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[processing_dtos.Evaluation](ctx, req).
		ContentType(tools.MimeJSON).
		Do(nil)

}

func (client *httpClient) Search(ctx context.Context, ids []uint) ([]processing_dtos.Evaluation, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("POST", tools.Url(client.config.Services.Processing, "/v1/evaluations/search"), tools.ToJSON(dtos.QueryIds[uint]{Ids: ids}))
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return utils.NoPointerE(
		tools.NewHttpRequest[[]processing_dtos.Evaluation](ctx, req).
			ContentType(tools.MimeJSON).
			Do(nil),
	)
}

func (client *httpClient) ProcessImage(ctx context.Context, data reporter_dtos.ProcessImage) (*http.Response, *dtos.ErrorOutbound) {
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

	bodyWriter.Boundary()

	bodyWriter.WriteField("reportId", strconv.Itoa(int(data.ReportID)))

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", tools.Url(client.config.Services.Processing, "/v1/evaluations"), bodyBuf)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[any](ctx, req).ContentType(contentType).Execute(nil)
}
