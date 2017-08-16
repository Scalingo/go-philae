package s3probe

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/juju/errgo"
)

type S3Probe struct {
	readOnly  bool
	region    string
	endpoint  string
	accessKey string
	secretKey string
	bucket    string
	name      string
}

func (p *S3Probe) Client() *s3.S3 {
	client := s3.New(session.New(&aws.Config{
		Region:      aws.String(p.region),
		Endpoint:    aws.String(p.endpoint),
		DisableSSL:  aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(p.accessKey, p.secretKey, "")}))

	return client
}

func (p *S3Probe) PutObject(key string, client *s3.S3) error {
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader("philae")),
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}

	// TODO: Maybe check the PutObjectOutput
	_, err := client.PutObject(input)
	if err != nil {
		return errgo.Notef(err, "Unable to write object")
	}
	return nil
}

func (p *S3Probe) DelObject(key string, client *s3.S3) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}

	_, err := client.DeleteObject(input)
	if err != nil {
		return errgo.Notef(err, "Unable to delete object")
	}
	return nil
}

func (p *S3Probe) ListObjects(client *s3.S3) error {
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(p.bucket),
		MaxKeys: aws.Int64(2),
	}

	_, err := client.ListObjects(input)
	if err != nil {
		return errgo.Notef(err, "Unable to read objetcs")
	}

	return nil
}

func (p *S3Probe) Check() error {
	client := p.Client()

	err := p.ListObjects(client)
	if err != nil {
		return err
	}

	if !p.readOnly {
		key := "philae-probe"
		err := p.PutObject(key, client)
		if err != nil {
			return err
		}

		err = p.DelObject(key, client)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *S3Probe) Name() string {
	return p.name
}

func NewS3Probe(name, accessKey, secretKey, bucket, region, endpoint string, readOnly bool) *S3Probe {
	return &S3Probe{
		readOnly:  readOnly,
		region:    region,
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		name:      name,
	}
}
