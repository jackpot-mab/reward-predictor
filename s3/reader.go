package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type File struct {
	hash int
}

type Reader interface {
	read(bucketName string, objectKey string) (File, error)
}

type ReaderImpl struct {
	s3Session *s3.S3
}

func MakeReader(config *AwsConfig) (Reader, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Endpoint:    &config.Endpoint,
		Credentials: credentials.NewStaticCredentials(config.AwsKey, config.AwsSecret, ""),
	})
	if err != nil {
		return nil, err
	}
	return &ReaderImpl{s3Session: s3.New(sess)}, nil
}

func (s *ReaderImpl) read(bucketName string, objectKey string) (File, error) {
	return nil, nil
}
