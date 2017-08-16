package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func cmdStatus(flag CommandFlags) {

}

func cmdDoUnarchive(flag CommandFlags) error {
	Lg.Infof("start to unarchiving...")
	if !CheckFolderExist(flag.DirSrc) {
		return fmt.Errorf("src folder not exist! : [%s]", flag.DirSrc)
	}

	if !CheckFolderExist(flag.DirTo) {
		return fmt.Errorf("destination folder not exist! : [%s]", flag.DirTo)
	}

	filePattern, err := regexp.Compile(flag.FilePattern)
	Lg.Infof("search for file matching: %s", flag.FilePattern)
	if err != nil {
		return fmt.Errorf("Cannot compile regex file pattern: [%s], err: %v",
			flag.FilePattern, err)
	}

	files, err := ScanDir(flag.DirSrc)
	if err != nil {
		return fmt.Errorf("Cannot read files in src dir! :[%s]. %v", flag.DirSrc, err)
	}

	Lg.Infof("found %d file desciptors under folder %s", len(files), flag.DirSrc)
	countUnarchived := 0
	for _, file := range files {
		resCheck := filePattern.MatchString(file.Name())
		if !resCheck {
			Lg.Infof("filename not match: %s", file.Name())
			continue
		}

		filepath := filepath.Join(flag.DirSrc, file.Name())
		cnt, err := ExtractTarGz(filepath, flag.DirTo)
		if err != nil {
			return fmt.Errorf("Cannot decompress file [%s]. err: %v", filepath, err)
		}
		Lg.Infof("decompressed %d files in %s", cnt, file.Name())
		countUnarchived += cnt
	}

	Lg.Infof("End of Unarchiving. Total archived %d files.", countUnarchived)
	return nil
}

func cmdDoArchive(flag CommandFlags) error {
	if flag.Prefix == "" {
		return fmt.Errorf("must set flag prefix")
	}

	Lg.Infof("start to archiving...")
	if !CheckFolderExist(flag.DirSrc) {
		return fmt.Errorf("src folder not exist! : [%s]", flag.DirSrc)
	}

	if !CheckFolderExist(flag.DirTo) {
		return fmt.Errorf("destination folder not exist! : [%s]", flag.DirTo)
	}

	filePattern, err := regexp.Compile(flag.FilePattern)
	if err != nil {
		return fmt.Errorf("Cannot compile regex file pattern: [%s], err: %v",
			flag.FilePattern, err)
	}

	files, err := ScanDir(flag.DirSrc)
	if err != nil {
		return fmt.Errorf("Cannot read files in src dir! :[%s]. %v", flag.DirSrc, err)
	}

	Lg.Infof("found %d file desciptors under folder %s", len(files), flag.DirSrc)
	countArchived := 0
	timeLimit := time.Now().Add(-time.Second * time.Duration(flag.RetainSeconds))
	Lg.Infof("archiving matching file with mod time before: %s", timeLimit.Format("06-01-02 15:04:05"))
	fisToArchive := make([]string, 0, flag.PackNumber)
	for _, file := range files {
		resCheck := checkFileOkayToArchive(file, timeLimit, filePattern)
		// Lg.Infof("testing result %d, file: %s", resCheck, file.Name())
		if resCheck != 0 {
			continue
		}
		filepath := filepath.Join(flag.DirSrc, file.Name())
		fisToArchive = append(fisToArchive, filepath)

		if len(fisToArchive) >= (int)(flag.PackNumber) {
			err = makeArchive(flag, fisToArchive)
			if err != nil {
				return err
			}
			countArchived += len(fisToArchive)
			fisToArchive = make([]string, 0, flag.PackNumber)
		}
	}

	if len(fisToArchive) > 0 {
		err = makeArchive(flag, fisToArchive)
		if err != nil {
			return err
		}
		countArchived += len(fisToArchive)
	}
	Lg.Infof("End. Total decompressed %d files.", countArchived)
	return nil
}

func makeArchive(flag CommandFlags, fps []string) error {
	now := time.Now()
	tgzFilename := fmt.Sprintf("%s_%s_%03d.tgz", flag.Prefix, now.Format("060102_150405"), now.Nanosecond()/int(1.e+6))
	tgzPath := filepath.Join(flag.DirTo, tgzFilename)
	Lg.Infof("compress to file: %s", tgzPath)
	err := CreateTgz(tgzPath, fps)
	if err != nil {
		return err
	}
	Lg.Infof("Success archived %d files in %s", len(fps), tgzFilename)
	for _, fpath := range fps {
		err = os.Remove(fpath)
		if err != nil {
			return err
		}
	}
	Lg.Infof("Success deleted %d archived files.", len(fps))
	return nil
}

// CreateTgz compress files into a tgz file.
func CreateTgz(tgzPath string, files []string) error {
	fTgz, err := os.Create(tgzPath)
	if err != nil {
		return fmt.Errorf("Cannot create tgz file: %s, err: %v", tgzPath, err)
	}
	wGz, _ := gzip.NewWriterLevel(fTgz, gzip.BestSpeed)
	defer wGz.Close()
	wTar := tar.NewWriter(wGz)
	defer wTar.Close()

	for _, fp := range files {
		file, err := os.Open(fp)
		if err != nil {
			return fmt.Errorf("Cannot open file to compress: %s, err: %v", fp, err)
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("Cannot read file status before compress: %s, err: %v", fp, err)
		}
		header := new(tar.Header)
		header.Name = filepath.Base(file.Name())
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		header.Typeflag = tar.TypeReg
		err = wTar.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("Cannot write file header before compress: %s, err: %v", fp, err)
		}
		_, err = io.Copy(wTar, file)
		if err != nil {
			return fmt.Errorf("Cannot copy file content during compression: %s, err: %v", fp, err)
		}
		time.Sleep(time.Millisecond * 2)
	}
	return nil
}

func checkFileOkayToArchive(f os.FileInfo, timeLimit time.Time,
	pattern *regexp.Regexp) int32 {
	if f.IsDir() {
		return 9
	}

	if timeLimit.Before(f.ModTime()) {
		return 1
	}

	if !pattern.MatchString(f.Name()) {
		return 2
	}

	return 0
}
