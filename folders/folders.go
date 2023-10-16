package folders

import "os"

func PhotoFolder() string {
	fld, _ := os.Getwd()

	return fld
}

func MoveFile(src, dest string, onlyCopy bool) error {
	if err := copyFile(src, dest); err != nil {
		return err
	}

	if !onlyCopy {
		return os.Remove(src)
	}

	return nil
}
