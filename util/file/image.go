package file

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	file_support "raise-child/constants/file/file_support"
	base64_file_fomrat "raise-child/constants/file/format/base64"

	"raise-child/constants/noti"
)

// Check data URL prefix
func IsImageByDataURL(imageUrl string) bool {
	for _, format := range []string{
		base64_file_fomrat.JPEG_FILE_FORMAT,
		base64_file_fomrat.JPG_FILE_FORMAT,
		base64_file_fomrat.PNG_FILE_FORMAT,
		base64_file_fomrat.GIF_FILE_FORMAT,
	} {
		if strings.HasPrefix(imageUrl, format) {
			return true
		}
	}

	return false
}

// Validate by decoding and checking magic bytes
func IsImageByMagicBytes(imageUrl string) bool {
	var imageComponents []string = strings.Split(imageUrl, ",")

	// Remove data URL prefix if present
	if len(imageComponents) > 1 {
		imageUrl = imageComponents[1]
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(imageUrl)
	if err != nil {
		return false
	}

	if len(data) < 4 {
		return false
	}

	// Check magic bytes for different image formats
	return detectImageType(data) != ""
}

// Check magic bytes to determine image format
func detectImageType(data []byte) string {
	if len(data) < 4 {
		return ""
	}

	// JPEG: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg"
	}

	// PNG: 89 50 4E 47 0D 0A 1A 0A
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 &&
		data[2] == 0x4E && data[3] == 0x47 {
		return "png"
	}

	// GIF: 47 49 46 38 (GIF8)
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return "gif"
	}

	// WebP: RIFF....WEBP
	if len(data) >= 12 && data[0] == 0x52 && data[1] == 0x49 &&
		data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 &&
		data[10] == 0x42 && data[11] == 0x50 {
		return "webp"
	}

	// BMP: 42 4D
	if data[0] == 0x42 && data[1] == 0x4D {
		return "bmp"
	}

	return ""
}

// Validate by attempting to decode as image using Go's image package
func IsImageByDecoding(imageUrl string) bool {
	var imageComponents []string = strings.Split(imageUrl, ",")

	// Remove data URL prefix if present
	if len(imageComponents) > 1 {
		imageUrl = imageComponents[1]
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(imageUrl)
	if err != nil {
		return false
	}

	var reader = bytes.NewReader(data)

	// Try to decode as image
	_, _, capturedErr := image.Decode(reader)

	return capturedErr == nil
}

// Validate specific image format
func ValidateSpecificFormat(imageUrl string, expectedFormat string) (bool, error) {
	// Remove data URL prefix if present
	base64Data := strings.Split(imageUrl, ",")
	if len(base64Data) > 1 {
		imageUrl = base64Data[1]
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(imageUrl)
	if err != nil {
		return false, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var reader = bytes.NewReader(data)

	switch strings.ToLower(expectedFormat) {
	case file_support.JPG_FORMAT, file_support.JPEG_FORMAT:
		_, err = jpeg.Decode(reader)
	case file_support.PNG_FORMAT:
		_, err = png.Decode(reader)
	case file_support.GIF_FORMAT:
		_, err = gif.Decode(reader)
	default:
		return false, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return err == nil, err
}

//func ValidateImage(imageUrl string)

// Helper function to extract MIME type from data URL
func extractMimeType(imageUrl string) string {
	if strings.HasPrefix(imageUrl, "data:") {
		var components []string = strings.Split(imageUrl, ";")
		if len(components) > 0 {
			return strings.TrimPrefix(components[0], "data:")
		}
	}

	return ""
}

// Helper function for magic bytes checking
func checkMagicBytes(data []byte) (bool, string) {
	var imageType string = detectImageType(data)
	return imageType != "", imageType
}
