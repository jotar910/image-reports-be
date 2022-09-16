package mappers

import (
	"image-reports/api-gateway/dtos"
	processing_dtos "image-reports/api-gateway/dtos/processing"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"
	user_dtos "image-reports/api-gateway/dtos/user"
	"image-reports/helpers/utils"
)

func MapReportsList(
	list dtos.PageableList[reporter_dtos.Report],
	users []user_dtos.UserResponse,
	evaluations []processing_dtos.Evaluation,
) dtos.PageableList[reporter_dtos.ReportOutbound] {
	res := dtos.PageableList[reporter_dtos.ReportOutbound]{
		Content:          make([]reporter_dtos.ReportOutbound, list.NumberOfElements),
		Page:             list.Page,
		TotalPages:       list.TotalPages,
		TotalElements:    list.TotalElements,
		NumberOfElements: list.NumberOfElements,
	}
	usersMap := make(map[uint]user_dtos.UserResponse, len(users))
	for _, user := range users {
		usersMap[user.Id] = user
	}
	evaluationsMap := make(map[uint]processing_dtos.Evaluation, len(users))
	for _, evaluation := range evaluations {
		evaluationsMap[evaluation.ReportID] = evaluation
	}
	for i, item := range list.Content {
		res.Content[i] = MapPartialReport(item)
		if mUser, ok := usersMap[item.UserID]; ok {
			res.Content[i].User = mUser.Email
		}
		if mEvaluation, ok := evaluationsMap[item.ID]; ok {
			evaluation := MapEvaluation(mEvaluation)
			res.Content[i].Evaluation = &evaluation
		}
	}
	return res
}

func MapReport(
	report reporter_dtos.Report,
	user user_dtos.UserResponse,
	evaluation processing_dtos.Evaluation,
) reporter_dtos.ReportOutbound {
	res := MapPartialReport(report)
	if user.Id != 0 {
		res.User = user.Email
	}
	if evaluation.ReportID != 0 {
		res.Evaluation = utils.Pointer(MapEvaluation(evaluation))
	}
	return res
}

func MapPartialReport(report reporter_dtos.Report) reporter_dtos.ReportOutbound {
	return reporter_dtos.ReportOutbound{
		ID:       report.ID,
		Name:     report.Name,
		Image:    report.ImageID,
		Status:   report.Status,
		Approval: report.Approval,
		Date:     report.Date,
	}
}

func MapEvaluation(
	evaluation processing_dtos.Evaluation,
) processing_dtos.EvaluationOutbound {
	return processing_dtos.EvaluationOutbound{
		Grade:      evaluation.Grade,
		Categories: evaluation.Categories,
	}
}

func MapReportsListToIds(list dtos.PageableList[reporter_dtos.Report]) []uint {
	ids := make([]uint, list.NumberOfElements)
	for i, item := range list.Content {
		ids[i] = item.ID
	}
	return ids
}

func MapReportCreationData(form reporter_dtos.ReportCreation, imageId string) reporter_dtos.ReportCreationData {
	return reporter_dtos.ReportCreationData{
		Name:     form.Name,
		Callback: form.Callback,
		ImageID:  imageId,
	}
}

func MapReportSaveImage(form reporter_dtos.ReportCreation, imageId string) reporter_dtos.SaveImage {
	return reporter_dtos.SaveImage{
		Image:   form.File,
		ImageID: imageId,
	}
}

func MapReportProcessImage(form reporter_dtos.ReportCreation, imageId string, reportId uint) reporter_dtos.ProcessImage {
	return reporter_dtos.ProcessImage{
		Image:    form.File,
		ImageID:  imageId,
		ReportID: reportId,
	}
}
