//go:build windows

package bins

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

const formatName string = "zip"

func installPath() string {
	return fmt.Sprintf("%s%cSystem32", os.Getenv("SYSTEMROOT"), os.PathSeparator)
}

func decompress(resp *http.Response) (io.ReadCloser, error) {
	// Have to copy the response as we need to get an io.ReaderAt for the zip API.
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, resp.Body)
	if err != nil {
		return nil, err
	}
	z, err := zip.NewReader(bytes.NewReader(buff.Bytes()), size)
	if err != nil {
		return nil, err
	}

	if len(z.File) != 1 {
		return nil, fmt.Errorf("Expecting one file in zip download, got:%d", len(z.File))
	}

	return z.File[0].Open()
}

func exeName(f string) string {
	return fmt.Sprintf("%s/%s.exe", installPath, f)
}
