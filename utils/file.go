package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

func ListFiles(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return []string{"error reading path"}
	}

	files := []string{"../"}
	for _, entry := range entries {
		files = append(files, entry.Name())
	}

	return files
}

func CopyFile(src, dst string) error {
	if src == dst {
		return fmt.Errorf("source and destination are the same")
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func FormatDate(t time.Time) string {
	return t.Format("Jan 02 2006 15:04")
}

func CreateNewFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

func CreateNewFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

func RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
