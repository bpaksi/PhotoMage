package findgroups

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bpaksi/PhotoMage/folders"
	"github.com/bpaksi/PhotoMage/photos"
)

type track struct {
	Date time.Time
	Long float64
	Lat  float64

	Files []string
}

func Work(files []photos.Photo, targetFolder string, onlyCopy bool, days int, miles int) {
	dayThreshold := 24 * time.Hour * time.Duration(days)
	groups := make([]*track, 0)

	for _, file := range files {
		long1 := convertDegreeAngle(file.Longitude)
		lat1 := convertDegreeAngle(file.Latitude)

		foundMatch := false
		for _, grp := range groups {
			match := true
			if match && days > 0 {
				if file.Taken.IsZero() {
					match = false
				} else {
					match = absDiff(file.Taken, grp.Date) <= dayThreshold
				}
			}
			if match && miles > 0 {
				dist := -1.0
				if long1 == 0 || lat1 == 0 {
					match = false
				} else {
					dist = distance(lat1, long1, grp.Lat, grp.Long)

					match = dist <= float64(miles)
				}
			}

			if match {
				grp.Files = append(grp.Files, file.FilePath)
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			groups = append(groups, &track{
				Date:  file.Taken,
				Long:  long1,
				Lat:   lat1,
				Files: []string{file.FilePath},
			})
		}
	}

	count := 0
	for _, grp := range groups {
		if len(grp.Files) < 2 {
			continue
		}

		destFolder := filepath.Join(targetFolder, fmt.Sprint(count+1))
		_ = os.MkdirAll(destFolder, 0700)

		for _, f := range grp.Files {
			moveFile(f, destFolder, onlyCopy)
		}

		count++
	}
	fmt.Println("Found ", count, " potential groups")
}

func absDiff(t1, t2 time.Time) time.Duration {
	if t1.After(t2) {
		return t1.Sub(t2)
	}

	return t2.Sub(t1)
}

func moveFile(srcFile string, targetFolder string, onlyCopy bool) {
	dest := filepath.Join(targetFolder, filepath.Base(srcFile))

	if err := folders.MoveFile(srcFile, dest, onlyCopy); err != nil {
		fmt.Printf("Error moving file %s: %s\n", filepath.Base(srcFile), err.Error())
	}
}

func convertDegreeAngle(s string) float64 {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return 0
	}

	degrees, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	minutes, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	seconds, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)

	//Decimal degrees =
	//   whole number of degrees,
	//   plus minutes divided by 60,
	//   plus seconds divided by 3600

	return degrees + (minutes / 60) + (seconds / 3600)
}
