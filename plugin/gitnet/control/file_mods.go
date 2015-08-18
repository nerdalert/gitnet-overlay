package control

import (
	"errors"
	"io"
	//	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

//func DeleteDir(dir string) error {
//	files, e := ioutil.ReadDir(dir)
//	if e != nil {
//		return e
//	}
//	for _, file := range files {
//		e = deleteFilePath(filepath.Join(dir, file.Name()), true)
//		if e != nil {
//			return e
//		}
//	}
//	return nil
//}

func DeleteDir(path string) error {
	destination_list, err := RecursiveDirectoryList(path)
	if err != nil {
		return err
	}
	for x := len(destination_list) - 1; x != -1; x-- {
		err = os.Remove(destination_list[x])
		if err != nil {
			return err
		}
	}
	return nil
}

func RecursiveDirectoryList(basedir string) ([]string, error) {
	dirs := []string{}
	return dirs, filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		dirs = append(dirs, MakeAbs(path))
		return nil
	})
}

//func deleteFilePath(file string, recursive bool) error {
//	_, e := os.Stat(file)
//	if e != nil {
//		return e
//	}
//	e = DeleteDir(file)
//	if e != nil {
//		return e
//	}
//	return os.Remove(file)
//}

//PathExists verifies that the given path points to a valid file or directory
func pathExists(path string) bool {
	if path == "" {
		return false
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Errorf("Failed to read file: %v", err)
		return false
	}
	return true
}

// copyFile copies a single file to a new destination
func copyFile(src string, dst string) error {
	sourcefile, err := os.Open(src)
	if err != nil {
		log.Warnf("Error locating the source config file [ %s ]: %s", src, err)

		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dst)
	if err != nil {
		log.Warnf("Error creating the config file [ %s ] to the git sync dir [ %s ] with an error of: %s", src, dst, err)
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(src)
		if err != nil {
			err = os.Chmod(dst, sourceinfo.Mode())
		}
	}
	return nil
}

// copyDir copies the contents of the directory at source to dest
func copyDir(srcDir, dstDir string) error {
	var err error
	if !pathExists(srcDir) {
		log.Errorf("Cannot copy directory at %s: file does not exist", srcDir)
		return errors.New("Invalid directory")
	}
	// Verify the source to be copied exists
	info, err := os.Stat(srcDir)
	if err != nil {
		return err
	}
	// Create the destination dir
	err = os.MkdirAll(dstDir, info.Mode())
	if err != nil {
		return err
	}
	// Read cwd and subpaths
	directory, _ := os.Open(srcDir)
	defer directory.Close()
	subdirectory, err := directory.Readdir(-1)
	for _, sub := range subdirectory {
		sourcePath := filepath.Join(srcDir, sub.Name())
		destPath := filepath.Join(dstDir, sub.Name())
		if sub.IsDir() {
			copyDir(sourcePath, destPath)
			continue
		}
		copyFile(sourcePath, destPath)
		if err != nil {
			log.Debugf("Failed to copy directory from %s to %s", sourcePath, destPath)
		}
	}
	return nil
}
