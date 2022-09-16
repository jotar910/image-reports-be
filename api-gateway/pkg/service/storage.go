package service

import (
	"context"
	"image-reports/api-gateway/dtos"
	"net/http"
)

func (svc service) GetFile(ctx context.Context, imageId string) (*http.Response, *dtos.ErrorOutbound) {
	return svc.storageClient.GetFile(ctx, imageId)
}
