package validators

import (
	shared_models "image-reports/shared/models"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Initialize() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("urlTypeRequired", urlTypeRequired)
		v.RegisterValidation("fileTypeRequired", fileTypeRequired)
	}
}

func urlTypeRequired(fl validator.FieldLevel) bool {
	creationType := fl.Parent().FieldByName("Type").Interface().(string)
	creationUrl := fl.Field().Interface().(string)
	return creationType != shared_models.ReportCreationTypeUrl || creationUrl != ""
}

func fileTypeRequired(fl validator.FieldLevel) bool {
	creationType := fl.Parent().FieldByName("Type").Interface().(string)
	creationFileField := fl.Field()
	return creationType != shared_models.ReportCreationTypeFile || !creationFileField.IsNil()
}
