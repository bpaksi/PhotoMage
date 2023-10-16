package finddups

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bpaksi/PhotoMage/folders"
	"github.com/bpaksi/PhotoMage/photos"
)

func Work(files []photos.Photo,
	targetFolder string,
	onlyCopy bool,
	useDate,
	useLocation,
	useAltitude,
	useCameraOrientation,
	useFileSize bool) {

	count := 0
	worked := make([]string, 0)
	for i, f1 := range files[:len(files)-1] {
		if alreadyWorked(worked, f1.FilePath) {
			continue
		}

		fileInfo1, _ := os.Stat(f1.FilePath)

		foundMatch := false
		for _, f2 := range files[i+1:] {
			if alreadyWorked(worked, f2.FilePath) {
				continue
			}

			fileInfo2, err := os.Stat(f2.FilePath)
			if err != nil {
				fmt.Println("Error getting size for " + f2.FilePath + ": " + err.Error())
				continue
			}

			match := true
			if match && useDate {
				match = f1.Taken == f2.Taken
			}
			if match && useLocation {
				match = f1.Longitude == f2.Longitude &&
					f1.Latitude == f2.Latitude
			}
			if match && useAltitude {
				match = f1.Altitude == f2.Altitude
			}
			if match && useCameraOrientation {
				match = f1.CameraOrientation == f2.CameraOrientation
			}
			if match && useFileSize {
				match = fileInfo1.Size() == fileInfo2.Size()
			}

			if match {
				foundMatch = true
				worked = append(worked, f2.FilePath)

				moveFile(f2.FilePath, count, targetFolder, onlyCopy)
			}
		}

		if foundMatch {
			moveFile(f1.FilePath, count, targetFolder, onlyCopy)
			count++
		}
	}

	fmt.Println("Found ", count, " possible duplicates")
}

func alreadyWorked(worked []string, name string) bool {
	for _, i := range worked {
		if i == name {
			return true
		}
	}
	return false
}

func moveFile(srcFile string, count int, targetFolder string, onlyCopy bool) {
	destFile, _ := ioutil.TempFile(targetFolder, fmt.Sprintf("dup-%d-*%s", count+1, filepath.Ext(srcFile)))
	destFile.Close()
	os.Remove(destFile.Name())

	if err := folders.MoveFile(srcFile, destFile.Name(), onlyCopy); err != nil {
		fmt.Printf("Error moving file %s: %s\n", filepath.Base(srcFile), err.Error())
	}
}
