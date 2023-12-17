package bussinessLogic

import (
	"fmt"
	"github.com/google/uuid"
	"hyper_api/internal/config"
	"hyper_api/internal/utils"
	"hyper_api/internal/utils/aws"
)

func GenerateImageByPrompt(prompt string) (string, error) {
	setting := config.GetConfig()
	ans, _ := utils.GeneratePhotoUsingDallE3(prompt, 1)
	if ans == nil || len(ans) == 0 {
		return "", fmt.Errorf("generate image failed")
	}
	url := ans[0]
	downloadResult := utils.DownloadImage(url)
	uploadReq := aws.PutObjectInput{
		Bucket:      setting.GenerateS3Bucket,
		Key:         fmt.Sprintf("%v.png", uuid.New().String()),
		Body:        downloadResult.Image,
		ContentType: "image/png",
	}
	sess, _ := aws.NewAWSSession()
	s3Client := aws.NewS3Client(sess)
	uploadedImage := s3Client.PutObject(uploadReq)
	if uploadedImage.Error != nil {
		return "", fmt.Errorf("error when upload image to s3 %v", uploadedImage.Error)
	}
	newImageUrl := fmt.Sprintf("%v/%v", setting.CDNHost, uploadedImage.Key)
	return newImageUrl, nil
}
