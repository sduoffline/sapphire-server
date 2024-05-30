package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Unzip 解压zip文件
func Unzip(zipData []byte, destDir string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}
	for _, f := range zipReader.File {
		err := writeUnzipFile(f, destDir)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func writeUnzipFile(f *zip.File, destDir string) error {
	fName := f.Name
	destPath := filepath.Join(destDir, fName)

	// 判断文件夹是否存在，主要是处理zip包含多层文件目录的情况
	if f.FileInfo().IsDir() && !isFileExist(destPath) {
		err := os.MkdirAll(destPath, os.ModePerm)
		return err
	}

	// 创建要写入的文件
	fw, err := os.Open(destPath)
	if err != nil {
		return err
	}

	defer func(fw *os.File) {
		err := fw.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(fw)

	fr, err := f.Open()
	if err != nil {
		return err
	}

	defer func(fr io.ReadCloser) {
		err := fr.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(fr)

	_, err = io.Copy(fw, fr)

	return err
}

// isFileExist 文件或目录是否存在
func isFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
