package validators

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func ImageValidator(field string, file *multipart.FileHeader, maxSize int64, availableExtensions string) error {
	if !imageSizeValidator(file, maxSize) {
		return fmt.Errorf("%s must be under %d bytes", field, maxSize)
	}
	if !imageExtensionValidator(file, availableExtensions) {
		return fmt.Errorf("%s must be an image file", field)
	}
	return nil
}

func imageSizeValidator(fileHeader *multipart.FileHeader, maxSize int64) bool {
	return fileHeader.Size <= maxSize
}

func imageExtensionValidator(fileHeader *multipart.FileHeader, extensionsList string) bool {
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	mimetype.SetLimit(512) // Bytes
	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return false
	}

	fmime := strings.ToLower(mtype.String())
	extensions := strings.Split(extensionsList, ",")
	for _, ext := range extensions {
		if fmt.Sprintf("image/%s", strings.ToLower(strings.TrimSpace(ext))) == fmime {
			return true
		}
	}
	return false
}
