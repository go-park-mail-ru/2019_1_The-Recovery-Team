package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"
	"time"
)

// SaveFile saves file to dir, if dir doesn't exist, creates it
func SaveFile(file io.Reader, dir, filename string) error {
	fileCopy, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fileCopy.Close()

	_, err = io.Copy(fileCopy, file)
	if err != nil {
		return err
	}
	return nil
}

// HashFilename returns hashed filename
func HashFileName(filename string, id uint64) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(filename + string(id) + time.Now().String()))
	if err != nil {
		return "", err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash[:16] + path.Ext(filename), nil
}
