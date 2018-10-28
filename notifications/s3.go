package notifications

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/s3"
	"bytes"
)

var s3Session *session.Session

const (
	S3_REGION = "eu-central-1"
)

type S3Conf interface {
	GetId() string
	GetKey() string
	GetToken() string
	GetRegion() string
	GetBucket() string
}

type S3Sevice struct {
	conf      S3Conf
	s3Session *s3.S3
}

func InitS3(conf S3Conf) *S3Sevice {

	aws_access_key_id := conf.GetId()
	aws_secret_access_key := conf.GetKey()
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, conf.GetToken())
	_, err := creds.Get()
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
		}).Error("Initialization creds on s3")
	}
	cfg := aws.NewConfig().WithRegion(conf.GetRegion()).WithCredentials(creds)
	s3Session := s3.New(session.New(), cfg)

	return &S3Sevice{
		conf:      conf,
		s3Session: s3Session,
	}
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (s3s S3Sevice) AddObjectToS3(path string, objBytes []byte) {

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	params := &s3.PutObjectInput{
		Bucket: aws.String(s3s.conf.GetBucket()),
		Key: aws.String(path),
		Body: bytes.NewReader(objBytes),
		ContentLength: aws.Int64(int64(len(objBytes))),
	}
	resp, err := s3s.s3Session.PutObject(params)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
		}).Error("Error while puting object s3")
	}


	log.WithFields(log.Fields{
		"s3_response":  awsutil.StringValue(resp),
	}).Info("Response for s3")

}
