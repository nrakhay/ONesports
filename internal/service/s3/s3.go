package s3

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nrakhay/ONEsports/internal/config"
)

var s3Service *s3.S3

func StartS3Session() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})

	if err != nil {
		slog.Error("Error creating session:", "error", err)
		return
	}

	s3Service = s3.New(sess)

	slog.Info("S3 service initialized")
}

func UploadBufferToS3(buffer *bytes.Buffer, key string) (string, error) {
	if buffer == nil {
		return "", errors.New("buffer is nil")
	}

	size := int64(buffer.Len())

	if size == 0 {
		return "", errors.New("buffer is empty")
	}

	// Upload to S3
	_, err := s3Service.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(config.BucketName),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer.Bytes()),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String("audio/ogg"),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("s3://%s/%s", config.BucketName, key)

	return fileURL, nil
}

func RetrieveFileFromS3(key string) ([]byte, error) {
	result, err := s3Service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	defer result.Body.Close()

	buffer, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
