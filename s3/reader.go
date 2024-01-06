package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"jackpot-mab/reward-predictor/exchange"
	"log"
)

type Reader interface {
	Read(objectKey string) (*exchange.File, error)
	List() []string
}

type ReaderImpl struct {
	s3Client   *s3.S3
	bucketName string
}

func MakeReader(config *AwsConfig) (Reader, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         &config.Endpoint,
		Credentials:      credentials.NewStaticCredentials(config.AwsKey, config.AwsSecret, ""),
		S3ForcePathStyle: aws.Bool(config.S3ForcePathStyle),
	})
	if err != nil {
		return nil, err
	}
	return &ReaderImpl{s3Client: s3.New(sess), bucketName: config.S3Bucket}, nil
}

func (s *ReaderImpl) List() []string {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(""),
	}

	resp, err := s.s3Client.ListObjects(params)

	if err != nil {
		log.Print(fmt.Sprintf("Error ocurred when listing bucket files: %s", err.Error()))
		return nil
	}

	var result []string
	for _, key := range resp.Contents {
		result = append(result, *key.Key)
	}

	return result
}

func (s *ReaderImpl) Read(objectKey string) (*exchange.File, error) {
	output, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return nil, err
	}

	return &exchange.File{
		Name:     objectKey,
		Body:     output.Body,
		Checksum: output.ETag,
	}, nil
}
