package validators

import (
	"kessler/internal/hashes"
	"kessler/internal/objects/files"
	"kessler/internal/s3utils"
)

func ValidateExtensionFromHash(fileManager s3utils.KesslerFileManager, hash hashes.KesslerHash, extension files.KnownFileExtension) error {
	filepath, err := fileManager.DownloadFileFromS3(hash)
	if err != nil {
		return err
	}
	err = ValidateExtensionFromFilepath(filepath, extension)
	return err
}
