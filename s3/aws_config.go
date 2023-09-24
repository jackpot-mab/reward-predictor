package s3

type AwsConfig struct {
	Region           string
	AwsKey           string
	AwsSecret        string
	Endpoint         string
	S3ForcePathStyle bool
}
