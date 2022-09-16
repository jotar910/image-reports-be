package http_client

import (
	"context"
	"net/http"

	"image-reports/api-gateway/dtos"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	tools "image-reports/api-gateway/helpers/http-client"

	"image-reports/helpers/configs"
)

type HttpClient interface {
	List(ctx context.Context, filters reporter_dtos.ListFilters) (*dtos.PageableList[reporter_dtos.Report], *dtos.ErrorOutbound)
	Get(ctx context.Context, id uint) (*reporter_dtos.Report, *dtos.ErrorOutbound)
	Create(ctx context.Context, data reporter_dtos.ReportCreationData) (*reporter_dtos.Report, *dtos.ErrorOutbound)
	Patch(ctx context.Context, id uint, patch reporter_dtos.ReportPatch) (*reporter_dtos.Report, *dtos.ErrorOutbound)
	Approval() (*dtos.ErrorOutbound, *dtos.ErrorOutbound)
}

type httpClient struct {
	config configs.GlobalConfig
}

func NewHttpClient(config configs.GlobalConfig) HttpClient {
	return &httpClient{config: config}
}

func (client *httpClient) List(ctx context.Context, filters reporter_dtos.ListFilters) (*dtos.PageableList[reporter_dtos.Report], *dtos.ErrorOutbound) {
	req, err := http.NewRequest("GET", tools.Url(client.config.Services.Reporter, "/v1/reports"), nil)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[dtos.PageableList[reporter_dtos.Report]](ctx, req).
		ContentType(tools.MimeJSON).
		AddQuery("page", "%d", filters.Page).
		AddQuery("count", "%d", filters.Count).
		Do(nil)
}

func (client *httpClient) Get(ctx context.Context, id uint) (*reporter_dtos.Report, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("GET", tools.Url(client.config.Services.Reporter, "/v1/reports/%d", id), nil)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[reporter_dtos.Report](ctx, req).
		ContentType(tools.MimeJSON).
		Do(nil)
}

func (client *httpClient) Create(ctx context.Context, data reporter_dtos.ReportCreationData) (*reporter_dtos.Report, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("POST", tools.Url(client.config.Services.Reporter, "/v1/reports"), tools.ToJSON(data))
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[reporter_dtos.Report](ctx, req).
		ContentType(tools.MimeJSON).
		Do(nil)
}

func (client *httpClient) Patch(ctx context.Context, id uint, patch reporter_dtos.ReportPatch) (*reporter_dtos.Report, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("PATCH", tools.Url(client.config.Services.Reporter, "/v1/reports/%d", id), tools.ToJSON(patch))
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[reporter_dtos.Report](ctx, req).
		ContentType(tools.MimeJSON).
		Do(nil)
}

func (client *httpClient) Approval() (*dtos.ErrorOutbound, *dtos.ErrorOutbound) {

	return nil, nil
}
