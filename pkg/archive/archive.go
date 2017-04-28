package archive

import (
	"fmt"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"
)

// cores is the number of logical CPU cores the Go runtime has available to it.
var cores = runtime.GOMAXPROCS(0)

type Result struct {
	Error error `json:"error"`
}

// Dirsf takes an array of directory paths as strings, and a formatting string for the file
// names, and produces .tar.gz archives for each of the given directories. If any of the directories
// don't exist, an error will be returned.
//
// The values in `dirnames` can be absolute, or relative paths for the directories. These are simply
// passed into stdlib functions that will resolve this for us.
//
// The `namefmt` needs to have a single `%s` and a single `%d` in it, for both the base dirname and
// the current unix timestamp, e.g. `"backup-%s-%d"`.
//
// Upon success, an array of the archive filenames will be returned.
func Dirsf(dirnames []string, namefmt string) ([]string, error) {
	archives := []string{}

	namefmt = fmt.Sprintf("%s.tar.gz", namefmt)
	timestamp := time.Now().Unix()

	limiter := make(chan bool, cores)
	errs := make(chan error, 1)

	// Prepare the limiter. We fill it with as many values as we want archives to be created in
	// concurrently. For now, this is the number of logical CPU cores available to the Go runtime.
	for i := 0; i < cores; i++ {
		limiter <- true
	}

	// Archive each directory, if any one fails, we stop then and return that first error.
	for _, dirname := range dirnames {
		basename := path.Base(dirname)

		dest := fmt.Sprintf(namefmt, basename, timestamp)
		dest = strings.Replace(dest, " ", "_", -1)

		// Add archive name to list result
		archives = append(archives, dest)

		select {
		case <-limiter:
			// Process archiving a directory asynchronously.
			go func(dirname string, dest string) {
				err := doArchive(dirname, dest)
				if err != nil {
					errs <- err
				}

				// Release use of limiter
				limiter <- true
			}(dirname, dest)
		case err := <-errs:
			if err != nil {
				return archives, err
			}
		}
	}

	// Wait for all workers to finish.
	for i := 0; i < cores; i++ {
		<-limiter
	}

	// Sort tar filenames so they still come out in alphabetical order.
	sort.Strings(archives)

	return archives, nil
}

// doArchive actually performs the archiving. Taking a path to archive, and returning an error if
// one occurred.
func doArchive(path string, dest string) error {
	time.Sleep(1 * time.Second)
	fmt.Printf("%s (%d)\n", path, time.Now().Unix())

	return nil
}

// @todo: How on earth do we test this?
//   Internal implementation maybe doesn't matter? As long we're blocking waiting for all of them to
//   finish then the tests should check the result, not what the process is actually doing.

// @todo: Can we cancel all of the processes? And if we do this, do we need to clean up?
//   ???