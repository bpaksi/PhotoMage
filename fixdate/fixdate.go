package fixdate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bpaksi/PhotoMage/folders"
	"github.com/bpaksi/PhotoMage/photos"
)

func Work(files []photos.Photo, targetFolder string, onlyCopy bool) {
	count := 0
	for _, f := range files {
		if f.Taken.IsZero() {
			continue
		}

		filename := filepath.Base(f.FilePath)

		fi, err := os.Stat(f.FilePath)
		if err != nil {
			fmt.Println("\tError getting info", filename, err, f.FilePath)
			continue
		}

		mtime := fi.ModTime()

		fmt.Println("fixdate", filename, f.Taken)

		if mtime != f.Taken {
			dest := filepath.Join(targetFolder, filename)
			if err := folders.MoveFile(f.FilePath, dest, onlyCopy); err != nil {
				fmt.Printf("Error moving file %s: %s\n", filename, err)

				continue
			}

			if err := os.Chtimes(dest, f.Taken, f.Taken); err != nil {
				fmt.Printf("Error changing time %s: %s\n", filename, err)

				continue
			}

			count++
		}
	}
	fmt.Println("Fixed ", count, " photos")
}

// func statTimes(name string) (atime, mtime, ctime time.Time, err error) {
// 	fi, err := os.Stat(name)
// 	if err != nil {
// 		return
// 	}
// 	mtime = fi.ModTime()
// 	stat := fi.Sys().(*syscall.Stat_t)
// 	atime = time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
// 	ctime = time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)
// 	return
// }
