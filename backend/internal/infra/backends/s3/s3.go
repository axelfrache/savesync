package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Backend implements domain.Backend for S3-compatible storage
type Backend struct {
	client *s3.Client
	bucket string
}

// New creates a new S3 backend
func New() *Backend {
	return &Backend{}
}

// Init initializes the backend with configuration
func (b *Backend) Init(cfg map[string]string) error {
	bucket, ok := cfg["bucket"]
	if !ok {
		return fmt.Errorf("bucket is required in config")
	}
	b.bucket = bucket

	region := cfg["region"]
	if region == "" {
		region = "us-east-1"
	}

	accessKey := cfg["access_key"]
	secretKey := cfg["secret_key"]
	endpoint := cfg["endpoint"] // For MinIO/Garage/Backblaze

	// DEBUG: Log configuration (without secret)
	fmt.Printf("[S3 Backend Init] bucket=%s, region=%s, endpoint=%s, access_key=%s, secret_key_length=%d\n",
		bucket, region, endpoint, accessKey, len(secretKey))

	// Create AWS config with static credentials
	var awsCfg aws.Config
	var err error

	if accessKey != "" && secretKey != "" {
		// Use LoadDefaultConfig with static credentials and EC2 IMDS DISABLED
		// This is the only way to prevent EC2 IMDS timeout errors
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
			// CRITICAL: Disable EC2 IMDS to prevent timeout errors
			config.WithEC2IMDSClientEnableState(imds.ClientDisabled),
		)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}
		fmt.Printf("[S3 Backend Init] Using static credentials\n")
	} else {
		// Use default credential chain only if no credentials provided
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(region),
		)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}
	}

	// Create S3 client
	clientOptions := []func(*s3.Options){}
	if endpoint != "" {
		clientOptions = append(clientOptions, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			// Check if path_style is explicitly set
			if pathStyle, ok := cfg["path_style"]; ok && pathStyle == "true" {
				o.UsePathStyle = true
			} else {
				o.UsePathStyle = true // Default to true for non-AWS S3
			}
		})
	}

	b.client = s3.NewFromConfig(awsCfg, clientOptions...)

	return nil
}

// StoreChunk stores a chunk in S3
func (b *Backend) StoreChunk(ctx context.Context, hash string, data []byte) error {
	key := fmt.Sprintf("chunks/%s", hash)

	// Check if chunk already exists (deduplication)
	exists, err := b.ChunkExists(ctx, hash)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

// LoadChunk loads a chunk from S3
func (b *Backend) LoadChunk(ctx context.Context, hash string) ([]byte, error) {
	key := fmt.Sprintf("chunks/%s", hash)

	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return data, nil
}

// DeleteChunk deletes a chunk from S3
func (b *Backend) DeleteChunk(ctx context.Context, hash string) error {
	key := fmt.Sprintf("chunks/%s", hash)

	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// ChunkExists checks if a chunk exists in S3
func (b *Backend) ChunkExists(ctx context.Context, hash string) (bool, error) {
	key := fmt.Sprintf("chunks/%s", hash)

	_, err := b.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a "not found" error
		if isNotFoundError(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object: %w", err)
	}

	return true, nil
}

// StoreManifest stores a snapshot manifest in S3
func (b *Backend) StoreManifest(ctx context.Context, snapshotID string, manifest []byte) error {
	key := fmt.Sprintf("manifests/%s.json", snapshotID)

	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(manifest),
	})
	if err != nil {
		return fmt.Errorf("failed to put manifest: %w", err)
	}

	return nil
}

// LoadManifest loads a snapshot manifest from S3
func (b *Backend) LoadManifest(ctx context.Context, snapshotID string) ([]byte, error) {
	key := fmt.Sprintf("manifests/%s.json", snapshotID)

	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest body: %w", err)
	}

	return data, nil
}

// DeleteManifest deletes a snapshot manifest from S3
func (b *Backend) DeleteManifest(ctx context.Context, snapshotID string) error {
	key := fmt.Sprintf("manifests/%s.json", snapshotID)

	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	return nil
}

// Close closes the backend (no-op for S3)
func (b *Backend) Close() error {
	return nil
}

// isNotFoundError checks if an error is a "not found" error
func isNotFoundError(err error) bool {
	// AWS SDK v2 error handling
	return err != nil && (err.Error() == "NotFound" || err.Error() == "NoSuchKey")
}
