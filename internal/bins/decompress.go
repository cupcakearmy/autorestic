//go:build !windows

package bins

const formatName string = "bz2"

func installPath() string {
	return "/usr/local/bin"
}

func decompress(resp *http.Response) (io.ReadCloser, error) {
	return bzip2.NewReader(resp.Body)
}

func exeName(f string) string {
	return f
}
