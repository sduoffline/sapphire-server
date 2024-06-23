package util

import (
	"archive/zip"
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Unzip 解压zip文件
func Unzip(zipData []byte, destDir string) error {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}
	// 解压读取到的文件
	for _, f := range zipReader.File {
		if isIgnoreFile(f) {
			continue
		}
		err := writeUnzipFile(f, destDir)
		if err != nil {
			return err
		}
	}
	slog.Debug("Unzip file success")
	return nil
}

func isIgnoreFile(f *zip.File) bool {
	if f.FileInfo().IsDir() {
		return true
	}
	if f.Name == "" {
		return true
	}
	if strings.Contains(f.Name, "__MACOSX") {
		return true
	}
	return false
}

func writeUnzipFile(f *zip.File, destDir string) error {
	var err error
	fName := f.Name
	destPath := destDir + "/" + fName
	println("unzip file: ", fName, " to ", destPath)
	// 判断文件夹是否存在，主要是处理zip包含多层文件目录的情况
	//if f.FileInfo().IsDir() && !isFileExist(destPath) {
	//	err := os.MkdirAll(destPath, os.ModePerm)
	//	return err
	//}

	// 创建文件
	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("close file error: ", err)
		}
	}(file)

	// 打开文件
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {
			println("close file error: ", err)
		}
	}(rc)

	// 写入文件
	_, err = io.Copy(file, rc)
	if err != nil {
		return err
	}

	return nil
}

//// isFileExist 文件或目录是否存在
//func isFileExist(filePath string) bool {
//	_, err := os.Stat(filePath)
//	if err != nil {
//		if os.IsExist(err) {
//			return true
//		}
//		return false
//	}
//	return true
//}
