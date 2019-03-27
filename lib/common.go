package lib

import (
	"os"
	"path/filepath"
	"strings"
)

func SplitAuthorWork(info os.FileInfo, path string) (string, string) {
	return filepath.Base(filepath.Dir(path)), strings.TrimSuffix(info.Name(), ".txt")
}
