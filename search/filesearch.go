package search

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Searches all first-level files in a directory for approximate matches above the given threshold.
func SearchDir(needle string, dirpath string, threshold float32) (*map[string]*[]string, error) {
	dir, openerr := os.Open(dirpath)
	if openerr != nil {
		return nil, openerr
	}

	dirinfo, staterr := dir.Stat()
	if staterr != nil {
		return nil, staterr
	}

	if !dirinfo.IsDir() {
		return nil, fmt.Errorf("file must be a directory")
	}

	filelist, direrr := dir.ReadDir(0)

	if direrr != nil {
		return nil, direrr
	}

	filesToMatches := make(map[string]*[]string)
	// Map each file to a slice of matches.
	for _, entry := range filelist {
		// Avoid searching subdirectories.
		// First-level search only.
		if !entry.IsDir() {
			filepath := dirpath + string(os.PathSeparator) + entry.Name()

			matches, matcherr := SearchFile(needle, filepath, threshold)
			if matcherr == nil {
				filesToMatches[filepath] = matches
			}
		}
	}

	return &filesToMatches, nil
}

// Searches file for approximate matches above the given threshold.
func SearchFile(needle string, filepath string, threshold float32) (*[]string, error) {
	contents, err := readFile(filepath)
	if err != nil {
		return nil, err
	}

	return match(needle, contents, threshold), nil
}

// Returns a string containing all text in the given file.
func readFile(filepath string) (string, error) {
	f, openerr := os.Open(filepath)
	if openerr != nil {
		return "", openerr
	}

	finfo, staterr := f.Stat()

	if staterr != nil {
		return "", staterr
	}

	if finfo.IsDir() {
		return "", fmt.Errorf("file may not be a directory")
	}

	input := bufio.NewScanner(f)
	lines := strings.Builder{}

	for input.Scan() {
		lines.WriteString(input.Text())
		// Maintain newlines to stay true to original content.
		lines.WriteRune('\n')
	}

	return lines.String(), nil
}
