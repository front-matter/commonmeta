// Package fileutils provides utility functions for commonmeta.
package fileutils

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
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

// UncompressContent extracts the content from a gz archive,
func UncompressContent(input []byte) ([]byte, error) {
	var output []byte

	reader := bytes.NewReader(input)
	file, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	output, err = io.ReadAll(reader)
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

// ReadGZFile opens a gz archive for reading
func ReadGZFile(filename string) ([]byte, error) {
	var input, output []byte
	var err error

	input, err = ReadFile(filename)
	if err != nil {
		return nil, err
	}

	output, err = UncompressContent(input)
	return output, err
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

// DownloadFile downloads content from the given URL.
func DownloadFile(url string, progress bool) ([]byte, error) {
	var output []byte

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Create a buffer to store the response body
	buf := new(bytes.Buffer)

	// If progress is enabled, copy response body to both the buffer and the progress bar
	if progress {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)
		_, err = io.Copy(io.MultiWriter(buf, bar), resp.Body)
	} else {
		_, err = io.Copy(buf, resp.Body)
	}
	if err != nil {
		return output, err
	}

	output = buf.Bytes()
	return output, nil
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

// WriteGZFile saves content as a gz file.
func WriteGZFile(filename string, output []byte) error {
	gzfile, err := os.Create(filename + ".gz")
	if err != nil {
		panic(err)
	}
	defer gzfile.Close()

	zw := gzip.NewWriter(gzfile)
	defer zw.Close()

	// Setting the Header fields is optional.
	zw.Name = path.Base(filename)
	zw.ModTime = time.Now()

	// Write the output to the zip archive.
	_, err = zw.Write(output)
	return err
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

// GetExtension extracts the file extension and checks if the output file should be compressed.
func GetExtension(filename string, ext string) (string, string, string) {
	var extension, compress string

	if filename != "" {
		extension = path.Ext(filename)
		switch extension {
		case ".gz":
			compress = "gz"
			// Remove the ".gz" extension from the filename
			filename = strings.TrimSuffix(filename, ".gz")
			extension = path.Ext(filename)
		case ".zip":
			compress = "zip"
			// Remove the ".zip" extension from the filename
			filename = strings.TrimSuffix(filename, ".zip")
			extension = path.Ext(filename)
		default:
			compress = ""
		}
		return filename, extension, compress
	}

	if ext == "" {
		ext = ".json"
	}
	extension = ext
	compress = ""
	return filename, extension, compress
}
