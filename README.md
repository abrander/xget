# xget

Have you ever found a file on your Linux filesystem and wondered where
it came from? `xget` is a simple download tool, that will download a
file from a URL and save it to current directory. It will furthermore add
an extended file attribute `user.url` to the file, containing the URL
from which the file was downloaded.

For good measure the sha1sum and sha256sum of the file is also added as
extended file attributes.

## Installation

```bash
$ go install github.com/abrander/xget@latest
```

## Usage

```bash
# Downlad a file.
$ xget https://example.com/somefile.txt

# Show all attributes for the downloaded file.
$ xattr -l somefile.txt

# Show the URL attribute.
$ xattr -p user.url somefile.txt
```
