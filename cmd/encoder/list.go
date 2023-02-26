package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileList []*File

func NewFileList(sourcePath, targetPath string) (FileList, error) {
	fileList := FileList{}
	err := filepath.Walk(sourcePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}

		ext := path.Ext(filePath)
		if strings.ToLower(ext) != SOURCE_EXT {
			return nil
		}

		optputPath := path.Join(
			targetPath,
			strings.TrimSuffix(
				strings.TrimPrefix(filePath, sourcePath),
				ext,
			),
		)

		file, err := NewFile(filePath, optputPath)
		if err != nil {
			return err
		}
		fileList = append(fileList, file)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("path '%s' walk error %w", sourcePath, err)
	}
	
	return fileList, nil
}
