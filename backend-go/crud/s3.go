package crud

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Configuration constants
const ()

// Client structure for S3
type KesslerFileManager struct {
	s3Client *s3.S3
	S3Bucket string
	RawDir   string
	TmpDir   string
}

func NewS3FileManager() *KesslerFileManager {
	CloudRegion := "us-west-1" // Your region here. Change if needed
	EndpointURL := "https://sfo3.digitaloceanspaces.com"
	S3Bucket := "kesslerproddocs"
	S3AccessKey := "your-access-key"
	S3SecretKey := "your-secret-key"
	RawDir := "raw/"
	TmpDir := os.TempDir()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(CloudRegion),
		Endpoint:    aws.String(EndpointURL),
		Credentials: credentials.NewStaticCredentials(S3AccessKey, S3SecretKey, ""),
	}))
	return &KesslerFileManager{
		s3Client: s3.New(sess),
		S3Bucket: S3Bucket,
		RawDir:   RawDir,
		TmpDir:   TmpDir,
	}
}

// Hash computation
func calculateBlake2bHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash, err := blake2b.New256(file)
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Upload file to S3
func (manager *KesslerFileManager) pushFileToS3(filePath, hash string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	filename := manager.RawDir + hash // change the filename accordingly
	_, err = manager.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(manager.S3Bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	return err
}

func (manager *KesslerFileManager) downloadFileFromS3(hash string) (string, error) {
	localFilePath := filepath.Join(manager.RawDir, hash)
	if _, err := os.Stat(localFilePath); err == nil {
		return "", errors.New("file already exists at the path")
	}

	filename := manager.RawDir + hash
	buffer := aws.NewWriteAtBuffer([]byte{})
	_, err := manager.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(manager.S3Bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(localFilePath, buffer.Bytes(), 0644)
	return localFilePath, err
}
