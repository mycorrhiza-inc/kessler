package s3utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"thaumaturgy/common/hashes"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/charmbracelet/log"
)

func downloadFile(url, dir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	randomName := fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
	filePath := filepath.Join(dir, randomName)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// Client structure for S3
type KesslerFileManager struct {
	S3Client *s3.S3
	S3Bucket string
	S3RawDir string
	RawDir   string
	TmpDir   string
}

func NewKeFileManager() *KesslerFileManager {
	CloudRegion := "sfo3" // Your region here. Change if needed
	EndpointURL := "https://sfo3.digitaloceanspaces.com"
	S3Bucket := "kesslerproddocs"
	S3AccessKey := os.Getenv("S3_ACCESS_KEY")
	S3SecretKey := os.Getenv("S3_SECRET_KEY")
	S3RawDir := "raw/"
	RawDir := "raw/"
	TmpDir := os.TempDir()
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(CloudRegion),
		Endpoint:    aws.String(EndpointURL),
		Credentials: credentials.NewStaticCredentials(S3AccessKey, S3SecretKey, ""),
	}))
	return &KesslerFileManager{
		S3Client: s3.New(sess),
		S3Bucket: S3Bucket,
		S3RawDir: S3RawDir,
		RawDir:   RawDir,
		TmpDir:   TmpDir,
	}
}

// Hash computation

func (manager *KesslerFileManager) getS3KeyFromHash(hash hashes.KesslerHash) string {
	return filepath.Join(manager.S3RawDir, hash.String())
}

func (manager *KesslerFileManager) getLocalPathFromHash(hash hashes.KesslerHash) string {
	return filepath.Join("/", manager.RawDir, hash.String())
}

// Upload file to S3
func (manager *KesslerFileManager) uploadFileToS3(filePath string) (hashes.KesslerHash, error) {
	// F_ile opened twice, potential for optimisation.
	hash_result, err := hashes.HashFromFile(filePath)
	if err != nil {
		return hashes.KesslerHash{}, fmt.Errorf("Error hashing file: %v", err)
	}
	return hash_result, manager.pushFileToS3GivenHash(filePath, hash_result)
}

func (manager *KesslerFileManager) pushFileToS3GivenHash(filePath string, hash hashes.KesslerHash) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileKey := manager.getS3KeyFromHash(hash)
	_, err = manager.S3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(manager.S3Bucket),
		Key:    aws.String(fileKey),
		Body:   file,
	})
	return err
}

func (manager *KesslerFileManager) downloadFileFromS3(hash hashes.KesslerHash) (string, error) {
	// listInput := s3.ListObjectsInput{Bucket: aws.String(manager.S3Bucket)}
	// result, err := manager.S3Client.ListObjects(&listInput)
	// log.Info(fmt.Sprintf("result: %v\nerror: %v\n", result, err))

	localFilePath := manager.getLocalPathFromHash(hash)
	if _, err := os.Stat(localFilePath); err == nil {
		return localFilePath, nil
	}
	fileKey := manager.getS3KeyFromHash(hash)
	log.Info(fmt.Sprintf("Attempting to download %s from s3 bucket %s\n", fileKey, manager.S3Bucket))
	result_s3, err := manager.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(manager.S3Bucket),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return "", err
	}
	log.Info(fmt.Sprintf("Got file from s3 successfully\n"))
	var body io.ReadCloser = result_s3.Body
	defer body.Close()

	// Create the file
	// TODO: Move this directory code so that it only runs on startup or something.
	// Ensure the directory structure exists
	if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directories for %s: %v", localFilePath, err)
	}
	out, err := os.Create(localFilePath)
	if err != nil {
		return "", err
	}
	log.Info(fmt.Sprintf("Created local file successfully\n"))
	defer out.Close()

	// Copy the body to the file
	_, err = io.Copy(out, body)
	if err != nil {
		return "", err
	}
	log.Info(fmt.Sprintf("Copied file successfully\n"))
	return localFilePath, nil
}
