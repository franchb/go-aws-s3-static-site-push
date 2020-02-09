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

// GetListOfChangedFiles walks the file tree rooted at the specified directory, compares
// MD5 checksums and returns a list of absolute paths of files which checksum
// is different of absent with the checksums map.
func GetListOfChangedFiles(localPath string, checksums map[string]string) []string {
	collector := make ([]string, 128)

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
	collector []string) func(osPathname string, de *godirwalk.Dirent) error {

	return func(osPathname string, de *godirwalk.Dirent) error {
		if !de.IsDir() {
			// TODO: some clear function for extracting base part from absolute path
			// i. e.  extract `/a/b/c/` from `/a/b/c/d/e.txt`
			relativeFile := strings.Replace(osPathname, localPath, "", 1)

			checksumRemote, _ := checksums[relativeFile]
			// TODO: benchmark with go isFileChecksumDifferFromRemote(...)
			if isFileChecksumDifferFromRemote(osPathname, checksumRemote) {
				// file walk is not concurrent, so no locks are needed
				collector = append(collector, osPathname)
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

// isFileChecksumDifferFromRemote checks file MD5 checksum and returns
// true if checksums are the same, returns false if checksums are different.
// Returns true if the remote checksum is empty (so file not exists at remote).
func isFileChecksumDifferFromRemote(pathToFile string, checksumRemote string) bool {
	if checksumRemote == "" {
		return true
	}
	sumString, err := md5File(pathToFile)
	if err != nil {
		log.Error().
			Err(err).
			Str("filename", pathToFile).
			Msg("can't open local file, will skip it")
		return false
	}
	return sumString != checksumRemote
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
