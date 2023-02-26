package main

import (
	"fmt"
	"image/png"
	"io"
	"os"
	"path"

	"github.com/wisdman/vsv-decoder/internal/bar"
	"github.com/wisdman/vsv-decoder/internal/vsv"
)

type File struct {
	*vsv.File
	targetPath string
}

func NewFile(sourcePath, targetPath string) (*File, error) {
	vsvFile, err := vsv.New(sourcePath)
	if err != nil {
		return nil, err
	}

	return &File{vsvFile, targetPath}, nil
}

func (f *File) TargetPath() string {
	return f.targetPath
}

func (f *File) Process() error {
	if err := f.createTargetPath(); err != nil {
		return err
	}

	if err := f.checkTargetPath(); err != nil {
		return err
	}

	if err := f.Open(); err != nil {
		return fmt.Errorf("file '%s' open error: %w", f.Path(), err)
	}

	fmt.Printf("Processsing '%s'\nHeader:%+v\n", f.Path(), f.Header())
	var progressBar = bar.New(f.FrameCount())

	for i := uint64(0); true; i += 1 {
		progressBar.Play(uint64(i))
		
		frame, err := f.Frame(i)
		
		if err == io.EOF {
			break
		}

		if err != nil {
			progressBar.Finish()
			return fmt.Errorf("file '%s' frame error: %w", f.Path(), err)
		}		

		targetFile := path.Join(f.TargetPath(), fmt.Sprintf("%06d%s", i, TARGET_EXT))
		pngFile, err := os.Create(targetFile)
		if err != nil {
			progressBar.Finish()
			return fmt.Errorf("file '%s' create PNG error: %w", f.Path(), err)
		}

		if err := png.Encode(pngFile, frame); err != nil {
			progressBar.Finish()
			return fmt.Errorf("file '%s' PNG encode error: %w\n", f.Path(), err)
		}
	}

	progressBar.Finish()
	return nil
}

func (f *File) createTargetPath() error {
	return os.MkdirAll(f.targetPath, os.ModePerm)
}

func (f *File) checkTargetPath() error {
	if stat, err := os.Stat(f.targetPath); err != nil || stat.IsDir() != true {
		if err != nil {
			return fmt.Errorf("file '%s' Target path '%s' error %w", f.Path(), f.targetPath, err)
		}
		return fmt.Errorf("file '%s' Target path '%s' is not a directory", f.Path(), f.targetPath)
	}

	if entry, err := os.ReadDir(f.targetPath); err != nil || len(entry) > 0 {
		if err != nil {
			return fmt.Errorf("file '%s' Target path '%s' error %w", f.Path(), f.targetPath, err)
		}
		return fmt.Errorf("file '%s' Target path '%s' is not empty directory", f.Path(), f.targetPath)
	}

	return nil
}