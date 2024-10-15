package crud

import (
	"crypto/blake2b"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Configuration constants
const (
	CloudRegion   = "us-west-1" // Your region here. Change if needed
	EndpointURL   = "https://sfo3.digitaloceanspaces.com"
	S3Bucket      = "kesslerproddocs"
	S3AccessKey   = "your-access-key"
	S3SecretKey   = "your-secret-key"
	RawDir        = "raw/"
	LocalCacheDir = "/path/to/cache"
	TmpDir        = os.TempDir()
)

// Client structure for S3
type S3FileManager struct {
	s3Client *s3.S3
}

func NewS3FileManager() *S3FileManager {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(CloudRegion),
		Endpoint:    aws.String(EndpointURL),
		Credentials: credentials.NewStaticCredentials(S3AccessKey, S3SecretKey, ""),
	}))
	return &S3FileManager{
		s3Client: s3.New(sess),
	}
}

// Hash computation
func calculateBlake2bHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := blake2b.New512()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Upload file to S3
func (manager *S3FileManager) pushFileToS3(filePath, hash string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	filename := RawDir + hash // change the filename accordingly
	_, err = manager.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	return err
}

func (manager *S3FileManager) downloadFileFromS3(hash string) (string, error) {
	localFilePath := filepath.Join(LocalCacheDir, hash)
	if _, err := os.Stat(localFilePath); err == nil {
		return "", errors.New("file already exists at the path")
	}

	filename := RawDir + hash
	buffer := aws.NewWriteAtBuffer([]byte{})
	_, err := manager.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(localFilePath, buffer.Bytes(), 0644)
	return localFilePath, err
}
