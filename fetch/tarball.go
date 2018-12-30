package fetch

import (
	"archive/tar"
	"compress/gzip"
	"os"

	"github.com/atecce/canon/fs"
)

// TarballFetcher hits https://gutenberg.org and writes the text directly into a tarball
//
// tarball produces a single clean artifact which is easily moved around with file
// operations
//
// however, it is not parallelizable or resilient to failure. when it exits you will
// mostly likely need to start it from scratch. in addition, because you need to
// write file sizes in tar headers, it can create a considerable footprint counting
// bytes in memory
type TarballFetcher struct {
	Root string

	tw *tar.Writer
}

func (tf TarballFetcher) MkRoot() error {

	f, err := os.Create(tf.Root)
	if err != nil {
		return err
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	tf.tw = tw

	return nil
}

func (tf TarballFetcher) MkAuthorDir(name string) error {
	return nil
}

func (tf TarballFetcher) Fetch(url, path string) error {
	if err := fs.GetTarFile(url, path, tf.tw); err != nil {
		return err
	}
	return nil
}
