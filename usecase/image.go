package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Image interface {
	CreatePresignedURL(fileName, contentType *string) (string, error)
}

type ImageUseCase struct {
}

func NewImageUseCase() *ImageUseCase {
	return &ImageUseCase{}
}

func (i *ImageUseCase) CreatePresignedURL(fileName, contentType *string) (url string, err error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))

	s3Client := s3.New(sess)
	fmt.Println(*contentType)
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(*fileName),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(*contentType),
	})
	url, err = req.Presign(time.Minute)
	fmt.Println(url)
	return
}
