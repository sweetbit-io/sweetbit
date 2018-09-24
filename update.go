package main

import (
	"github.com/inconshreveable/go-update"
	"net/http"
	"compress/gzip"
	"archive/tar"
	"io"
	"github.com/pkg/errors"
	"path/filepath"
)

func doUpdate(url string) error {
	// TODO: Check that we're updating with the right arch
	// TODO: Check archive signature

	resp, err := http.Get(url)
	if err != nil {
		return errors.Errorf("Update: Download failed: %s", err.Error())
	}

	defer resp.Body.Close()

	decompressedStream, err := gzip.NewReader(resp.Body)
	if err != nil {
		return errors.Errorf("Update: Decompressing failed: %s", err.Error())
	}

	tarReader := tar.NewReader(decompressedStream)
	binaryFound := false

	// Search for binary file
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Errorf("Update: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeReg:
			file := filepath.Base(header.Name)
			// Naively, take the first file called "sweetd" as the new binary
			if file == "sweetd" {
				binaryFound = true
				break
			}
		}
	}

	if !binaryFound {
		return errors.New("Update: No suitable binary found.")
	}

	err = update.Apply(tarReader, update.Options{})
	if err != nil {
		if err = update.RollbackError(err); err != nil {
			return errors.Errorf("Update: Failed to rollback from bad update: %s", err.Error())
		}
	}

	return nil
}
