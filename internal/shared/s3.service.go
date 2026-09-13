package shared

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type S3Service struct {
	Client *s3.Client
	Bucket string
}

type S3Config struct {
	Bucket   string
	Region   string
	Endpoint string
}

func NewS3Service(ctx context.Context, cfg S3Config) (*S3Service, error) {
	s3DefaultConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %v", err)
	}
	client := s3.NewFromConfig(s3DefaultConfig, func(opt *s3.Options) {
		opt.EndpointResolverV2 = &EndpointResolver{Endpoint: cfg.Endpoint}
		opt.UsePathStyle = true
	})

	return &S3Service{client, cfg.Bucket}, nil
}

func (s *S3Service) UploadFile(ctx context.Context, key string, body io.Reader, contentType string) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

func (s *S3Service) GetPresignedURl(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.Client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

type EndpointResolver struct {
	Endpoint string
}

func (r *EndpointResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	if r.Endpoint != "" {
		params.Endpoint = aws.String(r.Endpoint)
	}
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}
