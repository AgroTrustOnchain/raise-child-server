package cloudinary

import (
	"fmt"
	"net/url"
	"os"
	"raise-child/constants/env"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	api "github.com/cloudinary/cloudinary-go/v2/api"
)

var (
	apiKey    string = os.Getenv("CLOUDINARY_API_KEY")
	apiSecret string = os.Getenv("CLOUDINARY_API_SECRET")
	cloudName string = os.Getenv("CLOUDINARY_CLOUD_NAME")
)

func GenerateSignature(timestamp int64, folder string) string {
	apiSecret := os.Getenv(env.CLOUDINARY_API_SECRET)

	paramsToSign := url.Values{
		"timestamp": {strconv.FormatInt(timestamp, 10)},
		"folder":    {folder},
	}
	signature, err := api.SignParameters(paramsToSign, apiSecret)
	if err != nil {
		return ""
	}

	return signature
}

func GetUploadUrl() string {
	cloudinaryUrl := os.Getenv(env.CLOUDINARY_UPLOAD_URL)
	uploadUrl := cloudinaryUrl + fmt.Sprintf("/%s/image/upload", cloudName)
	return uploadUrl
}

// agrotrust-uploads/pqbxghtmovrjfjgx0ait

func GetImageUrl(publicId string) string {
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	myImage, _ := cld.Image(publicId)
	imageUrl, _ := myImage.String()

	return imageUrl
}
