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

func (s *S3Client) PutObjects(req []PutObjectInput) ([]PutObjectOutput, error) {
	n := len(req)
	ret := make([]PutObjectOutput, n)
	var wg sync.WaitGroup

	for idx, v := range req {
		wg.Add(1)
		go func(obj PutObjectInput, i int) {
			defer wg.Done()
			contentType := obj.ContentType
			if contentType == "" {
				contentType = "binary/octet-stream"
			}
			_, err := s.client.PutObject(&s3.PutObjectInput{
				Bucket:      aws.String(obj.Bucket),
				Key:         aws.String(obj.Key),
				Body:        bytes.NewReader(obj.Body),
				ContentType: aws.String(contentType),
			})
			ret[i] = PutObjectOutput{
				Success: err == nil,
				Bucket:  obj.Bucket,
				Key:     obj.Key,
				Error:   err,
			}
		}(v, idx)
	}
	wg.Wait()
	return ret, nil
}
