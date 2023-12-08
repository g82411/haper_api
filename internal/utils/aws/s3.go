package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

type S3Client struct {
	client *s3.S3
}

type PutObjectInput struct {
	Bucket      string
	Key         string
	Body        []byte
	ContentType string
}

type PutObjectOutput struct {
	Success bool
	Bucket  string
	Key     string
	Error   error
}

func NewS3Client(sess *session.Session) *S3Client {
	return &S3Client{
		client: s3.New(sess),
	}
}

func (s *S3Client) PutObject(req PutObjectInput) PutObjectOutput {
	contentType := req.ContentType
	if contentType == "" {
		contentType = "binary/octet-stream"
	}
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(req.Bucket),
		Key:         aws.String(req.Key),
		Body:        bytes.NewReader(req.Body),
		ContentType: aws.String(contentType),
	})
	return PutObjectOutput{
		Success: err == nil,
		Bucket:  req.Bucket,
		Key:     req.Key,
		Error:   err,
	}
}

func (s *S3Client) PutObjects(req []PutObjectInput) ([]PutObjectOutput, error) {
	n := len(req)
	ret := make([]PutObjectOutput, n)
	var wg sync.WaitGroup

	for idx, v := range req {
		wg.Add(1)
		go func(obj PutObjectInput, i int) {
			defer wg.Done()
			ret[i] = s.PutObject(obj)
		}(v, idx)
	}
	wg.Wait()
	return ret, nil
}
