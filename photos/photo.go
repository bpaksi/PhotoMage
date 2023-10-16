package photos

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xiam/exif"
)

type Photo struct {
	FilePath string

	Taken     time.Time
	Altitude  float64
	Longitude string
	Latitude  string

	CameraOrientation float64
}

const takenFormat = "2006:01:02 15:04:05"

// var whitelist = []string{
// 	"Exif Version",
// 	"Date and Time",
// 	"Date and Time (Digitized)",
// 	"Date and Time (Original)",
// 	"Altitude",
// 	"Altitude Reference",

// 	"Longitude",
// 	"Latitude",
// 	"North or South Latitude",
// 	"East or West Longitude",
// 	"GPS Image Direction",
// 	"GPS Image Direction Reference", // T = True North
// 	"GPS Date",
// }

var NoEXIF = errors.New("")

func readFile(filePath string) (photo Photo, err error) {
	var data *exif.Data
	data, err = exif.Read(filePath)
	if err != nil {
		fmt.Println("Error reading photo properties: ", filepath.Base(filePath))

		return
	}
	photo.FilePath = filePath

	// // results := make(map[string]string)
	// for _, key := range whitelist {
	// 	if val, ok := data.Tags[key]; ok {
	// 		fmt.Printf("\t%s = %s\n", key, val)
	// 	}
	// }

	if taken, ok := data.Tags["Date and Time"]; ok {

		// fmt.Println(taken + ",  " + takenFormat)

		var err error
		photo.Taken, err = time.Parse(takenFormat, taken)

		if err != nil {
			fmt.Println("\t" + err.Error())
		}
	}

	photo.Longitude = data.Tags["Longitude"]
	photo.Latitude = data.Tags["Latitude"]
	if alt, ok := data.Tags["Latitude"]; ok {
		photo.Altitude, _ = strconv.ParseFloat(alt, 32)
	}

	if direct, ok := data.Tags["GPS Image Direction"]; ok {
		photo.CameraOrientation, _ = strconv.ParseFloat(direct, 32)
	}

	return
}
