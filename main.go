package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/xattr"
	"github.com/schollz/progressbar/v3"
)

func fileNameFromUrl(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
}

func main() {
	url := os.Args[1]
	resp, err := http.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

		os.Exit(1)
	}
	defer resp.Body.Close()

	fname := fileNameFromUrl(url)

	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create file: %v\n", err)

		os.Exit(1)
	}

	hash1 := sha1.New()
	hash256 := sha256.New()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	w := io.MultiWriter(f, bar, hash1, hash256)

	_, err = io.Copy(w, resp.Body)

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %T %v\n", err, err)

		os.Exit(1)
	}

	_ = xattr.FSet(f, "user.sha1sum", []byte(fmt.Sprintf("%x", hash1.Sum(nil))))
	_ = xattr.FSet(f, "user.sha256sum", []byte(fmt.Sprintf("%x", hash256.Sum(nil))))
	err = xattr.FSet(f, "user.url", []byte(url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting file attribute: %T %v\n", err, err)
	}

	f.Close()
}
