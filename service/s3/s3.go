package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nrakhay/ONEsports/config"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})

	if err != nil {
		fmt.Printf("Error creating session: %s\n", err)
		return
	}

	svc := s3.New(sess)
	fmt.Println("S3 service initialized:", svc)
}
