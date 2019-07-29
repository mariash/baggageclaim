// +build !linux

package volume

import (
	"io"
	"os"
	"path/filepath"

	"github.com/DataDog/zstd"
	"github.com/concourse/go-archive/tarfs"
	"github.com/concourse/go-archive/tgzfs"
)

func (streamer *tarGzipStreamer) In(stream io.Reader, dest string, privileged bool) (bool, error) {
	err := tgzfs.Extract(stream, dest)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (streamer *tarGzipStreamer) Out(w io.Writer, src string, privileged bool) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	var tarDir, tarPath string

	if fileInfo.IsDir() {
		tarDir = src
		tarPath = "."
	} else {
		tarDir = filepath.Dir(src)
		tarPath = filepath.Base(src)
	}

	return tgzfs.Compress(w, tarDir, tarPath)
}

func (streamer *tarZstdStreamer) In(stream io.Reader, dest string, privileged bool) (bool, error) {
	zstdStreamReader := zstd.NewReader(stream)

	err := tarfs.Extract(zstdStreamReader, dest)
	if err != nil {
		zstdStreamReader.Close()
		return true, err
	}

	err = zstdStreamReader.Close()
	if err != nil {
		return true, err
	}

	return false, nil
}

func (streamer *tarZstdStreamer) Out(w io.Writer, src string, privileged bool) error {
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	var tarDir, tarPath string

	if fileInfo.IsDir() {
		tarDir = src
		tarPath = "."
	} else {
		tarDir = filepath.Dir(src)
		tarPath = filepath.Base(src)
	}

	zstdStreamWriter := zstd.NewWriter(w)

	err = tarfs.Compress(zstdStreamWriter, tarDir, tarPath)
	if err != nil {
		zstdStreamWriter.Close()
		return err
	}

	err = zstdStreamWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
