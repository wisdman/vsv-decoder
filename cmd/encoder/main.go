package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const SOURCE_EXT = ".vsv"
const TARGET_EXT = ".png"

var (
  execName = "vsv-encoder"
  version  = "DEV"
  build 	 = "unknown"
  date     = "unknown" 
  goVersion  = runtime.Version()
  versionStr = fmt.Sprintf("%s v%s-%s %s %v", execName, version, build, date, goVersion)
)

func parseCLI() (sourcePath string, targetPath string) {
	var err error
	
	var printVersion bool
	var printHelp bool

	flag.StringVar(&sourcePath, "source", "", "Path to source folder")
	flag.StringVar(&targetPath, "target", "", "Path to target folder")
	flag.BoolVar(&printVersion, "version", false, "Print version information")
	flag.BoolVar(&printHelp, "help", false, "Print help and usage information")

	flag.Parse()

	if printVersion {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if printHelp {
		fmt.Printf("%s\nDecode VSV Video format to PNG image sequence\n", versionStr)
		fmt.Printf("usage: %s -source /path/to/media/files -target /path/to/output/folder\n\n", execName)
		fmt.Println("options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if sourcePath == "" {
		fmt.Fprintf(os.Stderr, "missing required -source argument\ntry '%s -help' for usage information\n", execName)
		os.Exit(2)
	}

	if sourcePath, err = filepath.Abs(sourcePath); err != nil {
		fmt.Fprintf(os.Stderr, "source path '%s' convert to absolute path error: %s\n", sourcePath, err)
		os.Exit(2)
	}

	if _, err = os.Stat(sourcePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "source folder or file '%s' does not exist\n", sourcePath)
		os.Exit(2)
	}

	if targetPath == "" {
		fmt.Fprintf(os.Stderr, "missing required -target argument\ntry '%s -help' for usage information\n", execName)
		os.Exit(2)
	}

	if targetPath, err = filepath.Abs(targetPath); err != nil {
		fmt.Fprintf(os.Stderr, "target path '%s' convert to absolute path error: %s\n", targetPath, err)
		os.Exit(2)
	}

	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "target folder %s does not exist\n", targetPath)
		os.Exit(2)
	}

	return
}

func main() {
	sourcePath, targetPath := parseCLI()

	fileList, err := NewFileList(sourcePath, targetPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create file list error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Files to process:\n")
	for i, file := range fileList {
		fmt.Printf("ID: %d; IN: %s; OUT: %s\n", i, file.Path(), file.TargetPath())
	}

	for _, file := range fileList {
		err := file.Process()
		if err != nil {
			fmt.Fprintf(os.Stderr, "File process error: %s\n", err)
			continue
		}
	}
}