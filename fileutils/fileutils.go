// Package fileutils provides utility functions for commonmeta.
package fileutils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

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

// UnzipContent extracts the content from a zip archive,
// optionally only extract the content with filename
func UnzipContent(input []byte, filename string) ([]byte, error) {
	var output []byte

	reader := bytes.NewReader(input)
	len := len(input)
	zipfile, err := zip.NewReader(reader, int64(len))
	if err != nil {
		return nil, err
	}

	// extract the files from the zip archive
	// optionally only extract the content with filename
	for _, file := range zipfile.File {
		if filename != "" && file.Name != filename {
			continue
		}
		// Open the zip file for reading
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

// ReadZIPFile opens a zip archive for reading
func ReadZIPFile(filename string, name string) ([]byte, error) {
	var input, output []byte
	var err error

	input, err = ReadFile(filename)
	if err != nil {
		return nil, err
	}

	output, err = UnzipContent(input, name)
	return output, err
}

// DownloadFile downloads content from the given URL and saves it as a file.
func DownloadFile(url string, filename string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = WriteFile(filename, body)
	return err
}

// WriteFile saves content as a file.
func WriteFile(filename string, output []byte) error {
	file, err := os.Create(path.Base(filename))
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

// WriteZIPFile saves content as a zip file.
func WriteZIPFile(filename string, output []byte) error {
	zipfile, err := os.Create(filename + ".zip")
	if err != nil {
		panic(err)
	}
	defer zipfile.Close()

	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()

	// create zip header
	header := &zip.FileHeader{
		Name:     path.Base(filename),
		Method:   zip.Deflate,
		Modified: time.Now(),
	}

	// Add the file header to the zip archive.
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write the output to the zip archive.
	_, err = writer.Write(output)
	if err != nil {
		return err
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
