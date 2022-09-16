package validators

import (
	"errors"
	reporter_dtos "image-reports/api-gateway/dtos/reporter"

	shared_models "image-reports/shared/models"
)

func ReportCreationValidator(form reporter_dtos.ReportCreation) error {
	if form.Type == shared_models.ReportCreationTypeUrl {
		if form.Url == "" {
			return errors.New("url is required for this report creation type")
		}
		return nil
	}
	if form.File == nil {
		return errors.New("file is required for this report creation type")
	}
	return nil
}
