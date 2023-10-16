package photos

import (
	"fmt"
	"os"
	"path/filepath"
)

func Scan(folder string) (photos []Photo) {
	photos = make([]Photo, 0)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		photo, readErr := readFile(path)
		if readErr != nil {
			return nil
		}

		photos = append(photos, photo)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", folder, err)
		return
	}

	return
}
