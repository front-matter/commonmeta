// Package fileutils provides utility functions for commonmeta.
package fileutils

import (
	"archive/zip"
	"bytes"
	"embed"
	"io"
	"os"
	"path"
	"path/filepath"
)

//go:embed *.zip
var Files embed.FS

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	output := make([]byte, info.Size())
	_, err = file.Read(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// ReadZIPFile opens a zip archive for reading
func ReadZIPFile(filename string) ([]byte, error) {
	var zipfile *zip.Reader
	if filename == "affiliations_ror.yaml.zip" {
		file, err := Files.ReadFile(filepath.Join("affiliations_ror.yaml.zip"))
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(file)
		len := len(file)
		zipfile, err = zip.NewReader(reader, int64(len))
		if err != nil {
			return nil, err
		}
	} else {
		zipfile, err := zip.OpenReader(filename)
		if err != nil {
			return nil, err
		}
		defer zipfile.Close()
	}
	var output []byte

	// Iterate through the files in the archive,
	for _, file := range zipfile.File {
		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		out, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		output = append(output, out...)
	}
	return output, nil
}

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

// GetExtension extracts the file extension and checks if the output file should be zipped.
func GetExtension(filename string, ext string) (string, string, bool) {
	var extension string
	var compress bool

	if filename != "" {
		extension = path.Ext(filename)
		if extension == ".zip" {
			compress = true

			// Remove the ".zip" extension from the filename
			filename = filename[:len(filename)-4]
			extension = path.Ext(filename)
		} else {
			compress = false
		}
		return filename, extension, compress
	} 
	
	if ext == "" {
		ext = ".json"
	}
	extension = ext
	compress = false
	return filename, extension, compress
}
