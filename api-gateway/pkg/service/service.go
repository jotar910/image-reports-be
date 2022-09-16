package service

import (
	"context"
	"net/http"

	"image-reports/api-gateway/dtos"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	user_dtos "image-reports/api-gateway/dtos/user"
	processing_client "image-reports/api-gateway/helpers/http-client/processing"
	reporter_client "image-reports/api-gateway/helpers/http-client/reporter"
	storage_client "image-reports/api-gateway/helpers/http-client/storage"
	user_client "image-reports/api-gateway/helpers/http-client/users"
)

type Service interface {
	Login(ctx context.Context, credentials user_dtos.UserCredentials) (*dtos.LoginResponse, *dtos.ErrorOutbound)
	CheckUserById(ctx context.Context, id uint) bool
	GetUserById(ctx context.Context, id uint) (*user_dtos.UserResponse, *dtos.ErrorOutbound)
	ListReports(ctx context.Context, filters reporter_dtos.ListFilters) (*dtos.PageableList[reporter_dtos.ReportOutbound], *dtos.ErrorOutbound)
	GetReport(ctx context.Context, id uint) (*reporter_dtos.ReportOutbound, *dtos.ErrorOutbound)
	GetFile(ctx context.Context, imageId string) (*http.Response, *dtos.ErrorOutbound)
	CreateReport(ctx context.Context, form reporter_dtos.ReportCreation) (*reporter_dtos.ReportOutbound, *dtos.ErrorOutbound)
	ReportApproval(ctx context.Context, id uint, patch reporter_dtos.ReportPatch) (*reporter_dtos.ReportOutbound, *dtos.ErrorOutbound)
}

type service struct {
	userClient       user_client.HttpClient
	reporterClient   reporter_client.HttpClient
	processingClient processing_client.HttpClient
	storageClient    storage_client.HttpClient
}

func NewService(
	usersClient user_client.HttpClient,
	reporterClient reporter_client.HttpClient,
	processingClient processing_client.HttpClient,
	storageClient storage_client.HttpClient,
) Service {
	return &service{
		usersClient,
		reporterClient,
		processingClient,
		storageClient,
	}
}
