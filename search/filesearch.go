package search

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type searchResult struct {
	filepath string
	matches  *[]string
}

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

	resultChan := make(chan searchResult)
	resultCount := 0

	// Map each file to a slice of matches.
	for _, entry := range filelist {
		// Avoid searching subdirectories.
		// First-level search only.
		if !entry.IsDir() {
			filepath := filepath.Join(dirpath, entry.Name())

			go searchFile(needle, filepath, threshold, resultChan)
			resultCount++
		}
	}

	results := make(map[string]*[]string)
	// Receive result for each file that is searched.
	for i := 0; i < resultCount; i++ {
		result := <-resultChan
		results[result.filepath] = result.matches
	}

	close(resultChan)
	return &results, nil
}

// Searches file for approximate matches above the given threshold.
func searchFile(needle string, filepath string, threshold float32, c chan searchResult) {
	contents, err := readFile(filepath)

	if err != nil {
		c <- searchResult{filepath: filepath, matches: nil}
	} else {
		c <- searchResult{
			filepath: filepath,
			matches:  match(needle, contents, threshold),
		}
	}
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
		return "", fmt.Errorf("file cannot be a directory")
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
