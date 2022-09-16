package service

import (
	"context"
	"fmt"
	"image-reports/api-gateway/dtos"
	processing_dtos "image-reports/api-gateway/dtos/processing"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	"image-reports/api-gateway/mappers"
	"image-reports/helpers/utils"

	"github.com/google/uuid"
)

func (svc service) ListReports(ctx context.Context, filters reporter_dtos.ListFilters) (*dtos.PageableList[reporter_dtos.ReportOutbound], *dtos.ErrorOutbound) {
	reportsList, oerr := svc.reporterClient.List(ctx, filters)
	if oerr != nil {
		return nil, oerr
	}
	reportsIdsQuery := mappers.MapReportsListToIds(*reportsList)
	users, oerr := svc.userClient.Search(ctx, reportsIdsQuery)
	if oerr != nil {
		return nil, oerr
	}
	evaluations, oerr := svc.processingClient.Search(ctx, reportsIdsQuery)
	if oerr != nil {
		return nil, oerr
	}
	return utils.Pointer(mappers.MapReportsList(*reportsList, users, evaluations)), nil
}

func (svc service) GetReport(ctx context.Context, id uint) (*reporter_dtos.ReportOutbound, *dtos.ErrorOutbound) {
	report, oerr := svc.reporterClient.Get(ctx, id)
	if oerr != nil {
		return nil, oerr
	}
	user, oerr := svc.userClient.GetUserById(ctx, id)
	if oerr != nil {
		return nil, oerr
	}
	evaluations, oerr := svc.processingClient.GetEvaluation(ctx, id)
	if oerr != nil {
		return nil, oerr
	}
	return utils.Pointer(mappers.MapReport(*report, *user, *evaluations)), nil
}

func (svc service) CreateReports(ctx context.Context, form reporter_dtos.ReportCreation) (*reporter_dtos.ReportOutbound, *dtos.ErrorOutbound) {
	imageId := uuid.NewString()
	report, oerr := svc.reporterClient.Create(ctx, mappers.MapReportCreationData(form, imageId))
	if oerr != nil {
		return nil, oerr
	}
	if _, oerr := svc.storageClient.SaveFile(ctx, mappers.MapReportSaveImage(form, imageId)); oerr != nil {
		return nil, oerr
	}
	if _, oerr := svc.processingClient.ProcessImage(ctx, mappers.MapReportProcessImage(form, imageId, report.ID)); oerr != nil {
		return nil, oerr
	}
	user, oerr := svc.userClient.GetUserById(ctx, report.UserID)
	if oerr != nil {
		return nil, oerr
	}
	return utils.Pointer(mappers.MapReport(*report, *user, processing_dtos.Evaluation{})), nil
}
