package Helpers

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

func TailFile(filepath string, lines int) string {
	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var linesInFile int
	for scanner.Scan() {
		linesInFile++
	}

	if linesInFile < lines {
		lines = linesInFile
	}

	f.Seek(0, 0)
	scanner = bufio.NewScanner(f)
	var lastLines []string
	for i := 0; scanner.Scan(); i++ {
		if i >= linesInFile-lines {
			lastLines = append(lastLines, scanner.Text())
		}
	}
	return strings.Join(lastLines, "")
}
