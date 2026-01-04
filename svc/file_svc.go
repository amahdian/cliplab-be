package svc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/amahdian/cliplab-be/global/errs"
	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type FileSvc interface {
	StoreInS3(fileURL, folder, objectKey string) (*string, error)
	GetFileDownloadUrl(objectKey string, expiry time.Duration) (link string, err error)
	DownloadFiles(urls []string) ([]string, error)
}

type fileSvc struct {
	ctx    context.Context
	config StorageConfig
	s3     *s3.S3
	Client *http.Client
}

func newFileSvc(ctx context.Context, config StorageConfig) FileSvc {
	if config.AwsConfig == nil && config.AwsConfig.Endpoint == nil {
		return &noopFileSvc{}
	}

	mySession := session.Must(session.NewSession())
	s3Service := s3.New(mySession, config.AwsConfig)

	return &fileSvc{
		ctx:    ctx,
		config: config,
		s3:     s3Service,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *fileSvc) StoreInS3(fileURL, folder, objectKey string) (*string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get file from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.Newf(errs.Internal, nil, "bad status getting file from URL: %s", resp.Status)
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: s.config.AwsConfig.Credentials,
		Region:      s.config.AwsConfig.Region,
		Endpoint:    s.config.AwsConfig.Endpoint,
	})

	uploader := s3manager.NewUploader(sess)

	fullKey := path.Join(folder, objectKey)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.config.ProductName),
		Key:    aws.String(fullKey),
		Body:   resp.Body,
	})
	if err != nil {
		return nil, err
	}

	img := fmt.Sprintf("%s/%s", s.config.ProductName, fullKey)

	return lo.ToPtr(img), nil
}

func (s *fileSvc) GetFileDownloadUrl(objectKey string, expiry time.Duration) (string, error) {
	// 1. Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Credentials: s.config.AwsConfig.Credentials,
		Region:      s.config.AwsConfig.Region,
		Endpoint:    s.config.AwsConfig.Endpoint,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	// 2. Create an S3 service client
	svc := s3.New(sess)

	// 3. Create the GetObjectRequest
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.ProductName),
		Key:    aws.String(objectKey),
	})

	// 4. Generate the pre-signed URL with the specified expiration
	presignedURL, err := req.Presign(expiry)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %w", err)
	}

	return presignedURL, nil
}

func (s *fileSvc) DownloadFiles(urls []string) ([]string, error) {
	// 1. Create a unique temporary directory for this batch of downloads.
	tempDir, err := os.MkdirTemp(os.TempDir(), "downloads_*")
	if err != nil {
		return nil, errs.Wrapf(err, "could not create temp directory")
	}

	logger.Infof("Created temporary directory: %s", tempDir)

	results := make([]string, len(urls))

	g, _ := errgroup.WithContext(s.ctx)
	for i, url := range urls {
		g.Go(func() error {
			filePath, err := downloadFile(url, tempDir)
			if err != nil {
				return err
			}
			results[i] = filePath

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		logger.Error("failed to download files:", err)
	}

	return results, nil
}

func downloadFile(fileURL, destDir string) (string, error) {
	// Make the HTTP GET request.
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the destination file path.
	safeFilename, err := generateSafeFilename(fileURL)
	if err != nil {
		return "", err
	}
	destPath := filepath.Join(destDir, safeFilename)

	// Create the destination file.
	file, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Stream the response body directly to the file. This is memory-efficient.
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func generateSafeFilename(fileURL string) (string, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("could not parse url: %w", err)
	}

	// Get the extension from the original filename in the URL path.
	ext := filepath.Ext(path.Base(u.Path))

	// Generate a new random name (UUID-like).
	randomBytes := make([]byte, 16) // 128 bits for a unique name
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("could not generate random filename: %w", err)
	}

	// Create a new filename using the hex-encoded random bytes and the original extension.
	newFilename := hex.EncodeToString(randomBytes) + ext

	return newFilename, nil
}

type noopFileSvc struct{}

func (s *noopFileSvc) StoreInS3(fileURL, folder, objectKey string) (*string, error) {
	return nil, nil
}

func (s *noopFileSvc) GetFileDownloadUrl(objectKey string, expiry time.Duration) (link string, err error) {
	return "", nil
}

func (s *noopFileSvc) DownloadFiles(urls []string) ([]string, error) {
	return nil, nil
}
