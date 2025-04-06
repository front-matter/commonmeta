// Package fileutils provides utility functions for commonmeta.
package fileutils

import (
	"archive/zip"
	"io"
	"os"
)

func WriteFile(filename string, output []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		panic(err)
	}
	return nil
}

func WriteZIPFile(filename string, output []byte) error {
	zipfile, err := os.Create(filename + ".zip")
	if err != nil {
		panic(err)
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	err = WriteFile(filename, output)
	if err != nil {
		panic(err)
	}

	fileToZip, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileToZip.Close()

	// Get the file info to create a zip header.
	fileInfo, err := fileToZip.Stat()
	if err != nil {
		panic(err)
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		panic(err)
	}

	header.Name = filename
	header.Method = zip.Deflate

	// Add the file header to the zip archive.
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		panic(err)
	}

	// Write the file contents to the zip archive.
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		panic(err)
	}

	return nil
}
