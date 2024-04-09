package tools

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func GetFileExtension(file *multipart.File) (string, error) {
	mtype, err := mimetype.DetectReader(*file)
	if err != nil {
		return "", err
	}

	//возвращаем reader в начало к нулевому байту
	(*file).Seek(0, io.SeekStart)

	return mtype.Extension(), nil
}

func VerifyMIMEType(file *multipart.File, verifyingType string) error {
	mtype, err := mimetype.DetectReader(*file)
	if err != nil {
		return err
	}

	//возвращаем reader в начало к нулевому байту
	(*file).Seek(0, io.SeekStart)

	filetype := strings.Split(mtype.String(), "/")[0]
	if filetype != verifyingType {
		return errors.New(fmt.Sprintf("it is not %s file", verifyingType))
	}
	return nil
}

func SaveFile(file io.Reader, filepath string, byteStart int64) error {
	filepathArray := strings.Split(filepath, string(os.PathSeparator))
	var folderPath string

	if len(filepathArray) > 1 {
		folderPath = strings.Join(filepathArray[0:len(filepathArray)-1], string(os.PathSeparator))
	}

	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	openedFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer openedFile.Close()

	_, err = openedFile.Seek(byteStart, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(openedFile, file)
	if err != nil {
		return err
	}
	return nil
}
