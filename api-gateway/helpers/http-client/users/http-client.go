package http_client

import (
	"context"
	"net/http"

	"image-reports/api-gateway/dtos"
	user_dtos "image-reports/api-gateway/dtos/user"
	tools "image-reports/api-gateway/helpers/http-client"
	"image-reports/helpers/configs"
	"image-reports/helpers/utils"
)

type HttpClient interface {
	CheckCredentials(credentials user_dtos.UserCredentials) (*user_dtos.UserResponse, *dtos.ErrorOutbound)
	CheckUserById(ctx context.Context, id uint) *dtos.ErrorOutbound
	GetUserById(ctx context.Context, id uint) (*user_dtos.UserResponse, *dtos.ErrorOutbound)
	Search(ctx context.Context, ids []uint) ([]user_dtos.UserResponse, *dtos.ErrorOutbound)
}

type httpClient struct {
	config configs.GlobalConfig
}

func NewHttpClient(config configs.GlobalConfig) HttpClient {
	return &httpClient{config: config}
}

func (client *httpClient) CheckCredentials(credentials user_dtos.UserCredentials) (*user_dtos.UserResponse, *dtos.ErrorOutbound) {
	resp, err := http.Post(tools.Url(client.config.Services.Users, "/v1/auth"), tools.MimeJSON, tools.ToJSON(credentials))
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.ReadJSON[user_dtos.UserResponse, dtos.ErrorResponse](resp)
}

func (client *httpClient) CheckUserById(ctx context.Context, id uint) *dtos.ErrorOutbound {
	req, err := http.NewRequest("GET", tools.Url(client.config.Services.Users, "/v1/auth/%d", id), nil)
	if err != nil {
		return dtos.NewInternalError(err.Error())
	}
	resp, oerr := tools.NewHttpRequest[any](ctx, req).Execute(req)
	if oerr != nil || resp.StatusCode != http.StatusOK {
		return tools.ExtractError[dtos.ErrorResponse](resp)
	}
	return nil
}

func (client *httpClient) GetUserById(ctx context.Context, id uint) (*user_dtos.UserResponse, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("GET", tools.Url(client.config.Services.Users, "/v1/users/%d", id), nil)
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return tools.NewHttpRequest[user_dtos.UserResponse](ctx, req).
		ContentType(tools.MimeJSON).
		Do(nil)
}

func (client *httpClient) Search(ctx context.Context, ids []uint) ([]user_dtos.UserResponse, *dtos.ErrorOutbound) {
	req, err := http.NewRequest("POST", tools.Url(client.config.Services.Users, "/v1/users/search"), tools.ToJSON(dtos.QueryIds[uint]{Ids: ids}))
	if err != nil {
		return nil, dtos.NewInternalError(err.Error())
	}
	return utils.NoPointerE(
		tools.NewHttpRequest[[]user_dtos.UserResponse](ctx, req).
			ContentType(tools.MimeJSON).
			Do(nil),
	)
}
