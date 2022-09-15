package mappers

import (
	"image-reports/processing/dtos"
	"image-reports/processing/models"

	"image-reports/helpers/services/kafka"
)

func MapToEvaluationDTO(evaluation *models.Evaluations) dtos.Evaluation {
	return dtos.Evaluation{
		ReportID:   evaluation.ReportID,
		ImageID:    evaluation.ImageID,
		Grade:      evaluation.Grade,
		Categories: MapToCategoriesDTO(evaluation.Categories),
	}
}

func MapToEvaluations(evaluation dtos.Evaluation) *models.Evaluations {
	return &models.Evaluations{
		ReportID:   evaluation.ReportID,
		ImageID:    evaluation.ImageID,
		Grade:      evaluation.Grade,
		Categories: MapToCategories(evaluation.Categories),
	}
}

func MapToCategoriesDTO(categories []models.Categories) []string {
	res := make([]string, len(categories))

	for i, c := range categories {
		res[i] = c.Name
	}

	return res
}

func MapToCategories(categories []string) []models.Categories {
	res := make([]models.Categories, len(categories))

	for i, c := range categories {
		res[i] = models.Categories{Name: c}
	}

	return res
}

func MapProcessedMessageToEvaluationDTO(message *kafka.ImageProcessedMessage) dtos.Evaluation {
	return dtos.Evaluation{
		ReportID:   message.ReportId,
		ImageID:    message.ImageId,
		Grade:      message.Grade,
		Categories: message.Categories,
	}
}
