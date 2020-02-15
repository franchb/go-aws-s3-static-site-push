package filechecksum

import (
	"crypto/md5"
	"fmt"
	"github.com/karrick/godirwalk"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

// GetListOfChangedFilesChan walks the file tree rooted at the specified directory, compares
// MD5 checksums and returns a list of absolute paths of files which checksum
// is different of absent with the checksums map.
func GetListOfChangedFilesChan(localPath string, checksums map[string]string) chan [2]string {
	collector := make(chan [2]string)

	opts := godirwalk.Options{
		Callback:            withFileHandler(localPath, checksums, collector),
		FollowSymbolicLinks: false,
		ErrorCallback:       withErrorsCallback(),
		Unsorted:            true, // set faster yet non-deterministic enumeration
	}

	if err := godirwalk.Walk(localPath, &opts); err != nil {
		log.Error().
			Err(err).
			Msg("there was error while walking the local path")
	}
	return collector
}

// withFileHandler returns a godirwalk.Walk file handler
func withFileHandler(localPath string, checksums map[string]string,
	collector chan [2]string) func(osPathname string, de *godirwalk.Dirent) error {

	return func(osPathname string, de *godirwalk.Dirent) error {
		if !de.IsDir() {
			// TODO: some clear function for extracting base part from absolute path
			// i. e.  extract `/a/b/c/` from `/a/b/c/d/e.txt`
			relativeFile := strings.Replace(osPathname, localPath, "", 1)

			remote, _ := checksums[relativeFile]
			local, err := md5File(relativeFile)
			if err != nil {
				return fmt.Errorf("failed to extract MD5 sum from local file: %w", err)
			}
			if remote == "" || local != remote {
				collector <- [2]string{osPathname, local}
			}
		}
		return nil
	}
}

// withErrorsCallback returns a godirwalk.Walk error callback handler
func withErrorsCallback() func(f string, err error) godirwalk.ErrorAction {
	return func(f string, err error) godirwalk.ErrorAction {
		// TODO: on which error cases should it Halt and
		//  which could SkipNode and log a warning?
		log.Warn().
			Err(err).
			Str("filename", f).
			Msg("there was error while walking the local path")

		return godirwalk.Halt
	}
}

// md5File opens file and calculates MD5 hash for it
func md5File(pathToFile string) (string, error) {
	f, err := os.Open(pathToFile)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", pathToFile, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().
				Err(err).
				Str("filename", pathToFile).
				Msg("failed to close file after calculating MD5 checksum")
		}
	}()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf(
			"can't calculate MS5 checksum for local file %s, will skip: %w",
			pathToFile, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
