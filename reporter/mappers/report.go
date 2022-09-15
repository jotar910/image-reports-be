package mappers

import (
	"image-reports/reporter/dtos"
	"image-reports/reporter/models"
	"math"

	shared_models "image-reports/shared/models"
)

func MapToReportsListDTO(reports []*models.Reports, total int64, page, limit int64) dtos.PageableList[dtos.Report] {
	reportsLen := len(reports)
	list := dtos.PageableList[dtos.Report]{
		Content:          make([]dtos.Report, reportsLen),
		Page:             page,
		TotalPages:       int64(math.Max(math.Ceil(float64(total)/float64(limit)), 1)),
		TotalElements:    total,
		NumberOfElements: int64(reportsLen),
	}
	for i, report := range reports {
		list.Content[i] = MapToReportDTO(report)
	}
	return list
}

func MapToReportDTO(report *models.Reports) dtos.Report {
	res := dtos.Report{
		ID:       report.ID,
		Name:     report.Name,
		UserID:   report.UserID,
		ImageID:  report.ImageID,
		Callback: report.Callback,
		Status:   report.Status,
		Date:     report.CreatedAt.String(),
	}
	if report.Status == shared_models.ReportStatusPublished {
		approval := MapToApprovalDTO(&report.Approval)
		res.Approval = &approval
	}
	return res
}

func MapToApprovalDTO(approval *models.Approvals) dtos.Approval {
	return dtos.Approval{
		UserID: approval.UserID,
		Status: approval.Status,
		Date:   approval.CreatedAt.String(),
	}
}
