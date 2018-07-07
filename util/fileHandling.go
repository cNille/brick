package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Ensure that a directory exists
func EnsureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

// Copy source file/directory to destination location.
func Copy(source, destination string) error {
	// Ensure source file exists
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	return copy(source, destination, info)
}

func copy(source, destination string, info os.FileInfo) error {
	if info.IsDir() {
		fmt.Printf("Copying dir: %s to %s... \n", source, destination)
		return directoryCopy(source, destination, info)
	}
	return fileCopy(source, destination, info)
}

// Copy one file from source to destination.
func fileCopy(source, destination string, info os.FileInfo) error {

	// Create the new file
	dFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dFile.Close()

	// Set same mode as the source file
	if err = os.Chmod(dFile.Name(), info.Mode()); err != nil {
		return err
	}

	// Open the source fiel
	sFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sFile.Close()

	// Copy content from source to destination.
	_, err = io.Copy(dFile, sFile)
	return err
}

// Recursively copy a whole directory to a new location
func directoryCopy(source, destination string, info os.FileInfo) error {

	// Create destination directory, create all directories in between
	if err := os.MkdirAll(destination, info.Mode()); err != nil {
		return err
	}

	// Get info from all files
	infos, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	// Iterate file info and do a copy command.
	for _, info := range infos {
		if info.Name() != "/srv" {
			if err := copy(
				filepath.Join(source, info.Name()),
				filepath.Join(destination, info.Name()),
				info,
			); err != nil {
				return err
			}
		} else {
			fmt.Println("Skipping /srv")
		}
	}
	return nil
}
