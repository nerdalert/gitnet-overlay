package control

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type NodeState struct {
	GitPoints   []string // fancy play on the word endpoint
	FileRecords []string
}

func NewNodeState(endpoints, fileRecords []string) *NodeState {
	return &NodeState{
		GitPoints:   endpoints,
		FileRecords: fileRecords,
	}
}
func MakeAbs(fpath string) string {
	if !filepath.IsAbs(fpath) {
		cwd, _ := os.Getwd()
		return filepath.Join(cwd, fpath)
	}
	return fpath
}

func setSlice(strslice []string) map[string]struct{} {
	m := make(map[string]struct{}, len(strslice))
	for _, str := range strslice {
		m[str] = struct{}{}
	}
	return m
}

func DirectoryList(basedir string) ([]string, error) {
	dirs := []string{}
	return dirs, filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
		dirs = append(dirs, MakeAbs(path))
		return nil
	})
}

func DiffDirectories(previousRecords, latestRecords string) (removed, added []string, err error) {
	previousSlice, err := DirectoryList(previousRecords)
	if err != nil {
		return nil, nil, err
	}
	latestSlice, err := DirectoryList(latestRecords)
	if err != nil {
		return nil, nil, err
	}

	set_a := setSlice(previousSlice[1:])
	set_b := setSlice(latestSlice[1:])
	full_path_a := MakeAbs(previousRecords)
	full_path_b := MakeAbs(latestRecords)
	for path, _ := range set_a {
		rel_path := filepath.Join(full_path_b, path[len(full_path_a):])
		if _, ok := set_b[rel_path]; !ok {
			removed = append(removed, rel_path)
		}
	}
	for path, _ := range set_b {
		rel_path := filepath.Join(full_path_a, path[len(full_path_b):])
		if _, ok := set_a[rel_path]; !ok {
			added = append(added, rel_path)
		}
	}
	// Check for file modifications with matching files
	existingRecords := FileModifiedCompare(latestSlice, previousSlice)
	node := NewNodeState(existingRecords, stripSliceFilepath(existingRecords))
	if len(added) > 0 {
		log.Debugf("New records learned %s", stripSliceFilepath(added))
		nodeAdded(added)
	}
	if len(removed) > 0 {
		log.Debugf("Existing records removed %s", removed)
		nodeRemoved(removed)
	}
	node.GitPoints = stripSliceFilepath(existingRecords)

	var current []string
	for _, c := range existingRecords {
		_, endpoint := path.Split(c)
		current = append(current, endpoint)
	}
	node.GitPoints = current
	deduped := removeDups(node.GitPoints)
	log.Debugf("Current record list: %s", deduped)
	return removed, added, err
}

func FileModifiedCompare(slice1 []string, slice2 []string) []string {
	var diff []string
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			_, s1File := path.Split(s1)
			absolute1, err := os.Stat(s1)
			if err == nil && !absolute1.IsDir() {
				diff = append(diff, s1)
			}
			found := true
			for _, s2 := range slice2 {
				absolute2, err := os.Stat(s2)
				if err == nil && !absolute2.IsDir() {
					if absolute1.Name() == absolute2.Name() {
						if ok := deepCompare(s1, s2); !ok {
							log.Debugf("Record modification event for [%s]", s1)
							log.Debugf("Record-1 [%s] does not equal Record-2 [%s], paths are [%s] and [%s]", absolute1.Name(), absolute2.Name(), s1, s2)
							nodeModified(s1)
						}
					}
				}
				_, s2File := path.Split(s2)
				if s1File == s2File {
					found = false
					break
				}
			}
			if !found {
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}

func deepCompare(file1, file2 string) bool {
	// Check file size
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	for {
		b1 := make([]byte, chunk)
		_, err1 := f1.Read(b1)
		b2 := make([]byte, chunk)
		_, err2 := f2.Read(b2)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func stripSliceFilepath(records []string) []string {
	var endpoints []string
	for _, rawRecord := range records {
		_, endpoint := path.Split(rawRecord)
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

func stripSingleFilepath(record string) string {
	_, endpoint := path.Split(record)
	return endpoint
}
