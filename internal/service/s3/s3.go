package s3

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nrakhay/ONEsports/internal/config"
)

func startSession() *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})

	if err != nil {
		fmt.Printf("Error creating session: %s\n", err)
		return nil
	}

	svc := s3.New(sess)

	fmt.Println("S3 service initialized:", svc)
	return svc
}

func UploadFileToS3(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	svc := startSession()

	// get file size and read content to buffer
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// upload to s3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(config.BucketName),
		Key:                  aws.String(filepath.Base(filePath)),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String("audio/ogg"),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
