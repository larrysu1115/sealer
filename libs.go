package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CheckFolderExist check folder exist.
func CheckFolderExist(path string) bool {
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || !fi.IsDir() {
		return false
	}
	return true
}

// ScanDir return lists of files
func ScanDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

//ExtractTarGz extracts tgz file.
func ExtractTarGz(tgzFile string, dirTo string) (int, error) {
	countSucc := 0
	gzipStream, err := os.Open(tgzFile)
	if err != nil {
		return countSucc, fmt.Errorf("Cannot open tgz file %s, err: %v", tgzFile, err)
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)

	if err != nil {
		return countSucc, fmt.Errorf("ExtractTarGz: NewReader failed, err: %v", err)
	}

	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return countSucc, fmt.Errorf("ExtractTarGz: Next() failed: %s", err.Error())
		}
		switch header.Typeflag {
		case tar.TypeDir:
			toName := filepath.Join(dirTo, header.Name)
			if err := os.Mkdir(toName, 0755); err != nil {
				return countSucc, fmt.Errorf("ExtractTarGz: Mkdir() failed: %s, file: %s", err.Error(), toName)
			}
		case tar.TypeReg:
			toName := filepath.Join(dirTo, header.Name)
			outFile, err := os.Create(toName)
			if err != nil {
				return countSucc, fmt.Errorf("ExtractTarGz: Create() failed: %s, file: %s", err.Error(), toName)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return countSucc, fmt.Errorf("ExtractTarGz: Copy() failed: %s, file: %s", err.Error(), toName)
			}
			if err := os.Chtimes(toName, header.ModTime, header.ModTime); err != nil {
				return countSucc, fmt.Errorf("ExtractTarGz: Chtimes() failed: %s, file: %s", err.Error(), toName)
			}
			countSucc++
		default:
			return countSucc, fmt.Errorf(
				"ExtractTarGz: uknown type: %x in %s",
				header.Typeflag,
				header.Name)
		}
	}

	return countSucc, nil
}
