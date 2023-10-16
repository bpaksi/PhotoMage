package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bpaksi/PhotoMage/finddups"
	"github.com/bpaksi/PhotoMage/findgroups"
	"github.com/bpaksi/PhotoMage/fixdate"
	"github.com/bpaksi/PhotoMage/photos"
)

type Loc struct {
	Long string
	Lat  string
}

var (
	cmd          string
	onlyCopy     bool
	targetFolder string

	grpDays  int
	grpMiles int

	dupWhenTaken      bool
	dupOnLocation     bool
	dupOnAltitude     bool
	dupOnImgDirection bool
	dupFileSize       bool
)

func main() {
	srcFolder, ok := args()
	if !ok {
		os.Exit(1)
	}

	if srcFolder == "" {
		srcFolder, _ = os.Getwd()
	}

	fmt.Println("Source Folder", srcFolder)

	cwd, _ := os.Getwd()
	targetFolder = filepath.Clean(filepath.Join(cwd, targetFolder))
	if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
		if err = os.MkdirAll(targetFolder, 0700); err != nil {
			log.Fatalln(err)
		}
	}

	files := photos.Scan(srcFolder)
	fmt.Printf("Found %d number of photos\n", len(files))

	switch cmd {
	case "dups":
		finddups.Work(files, targetFolder, onlyCopy,
			dupWhenTaken, dupOnLocation, dupOnAltitude, dupOnImgDirection, dupFileSize)

	case "groups":
		findgroups.Work(files, targetFolder, onlyCopy, grpDays, grpMiles)

	case "fixdate":
		fixdate.Work(files, targetFolder, onlyCopy)

	}
}

func args() (string, bool) {
	dupCmd := flag.NewFlagSet("dups", flag.ExitOnError)
	dupCmd.BoolVar(&onlyCopy, "copyonly", false, "Don't move source file only copy changes")
	dupCmd.StringVar(&targetFolder, "dest", "./out", "Folder where changed filed reside")

	dupCmd.BoolVar(&dupWhenTaken, "date", true, "Use when taken when comparing")
	dupCmd.BoolVar(&dupOnLocation, "location", true, "Use location when comparing")
	dupCmd.BoolVar(&dupOnAltitude, "altitude", false, "Use atltitude when comparing")
	dupCmd.BoolVar(&dupOnImgDirection, "direction", false, "Use camera facing direction when comparing")
	dupCmd.BoolVar(&dupFileSize, "size", true, "Use file size when comparing")

	grpCmd := flag.NewFlagSet("groups", flag.ExitOnError)
	grpCmd.BoolVar(&onlyCopy, "copyonly", false, "Don't move source file only copy changes")
	grpCmd.StringVar(&targetFolder, "dest", "./out", "Folder where changed filed reside")
	grpCmd.IntVar(&grpDays, "days", 3, "# of day represent a group")
	grpCmd.IntVar(&grpMiles, "miles", 25, "# of miles represent a group (<= 0 to disable)")

	dateCmd := flag.NewFlagSet("fixdate", flag.ExitOnError)
	dateCmd.BoolVar(&onlyCopy, "copyonly", false, "Don't move source file only copy changes")
	dateCmd.StringVar(&targetFolder, "dest", "./out", "Folder where changed filed reside")

	if len(os.Args) < 2 {
		fmt.Println("Usage")
		fmt.Println("\tphotomage <cmd> [<switches>] [source folder]")

		os.Exit(1)
	}

	cmd = os.Args[1]
	switch cmd {

	case "dups":
		if err := dupCmd.Parse(os.Args[2:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				dupCmd.Usage()
				return "", false
			}

			fmt.Println("Error occured:", err)
			return "", false
		}
		return dupCmd.Arg(0), true

	case "groups":
		if err := grpCmd.Parse(os.Args[2:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				grpCmd.Usage()
				return "", false
			}

			fmt.Println("Error occured:", err)
			return "", false
		}
		return grpCmd.Arg(0), true

	case "fixdate":
		if err := dateCmd.Parse(os.Args[2:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				dateCmd.Usage()
				return "", false
			}

			fmt.Println("Error occured:", err)
			return "", false
		}
		return dateCmd.Arg(0), true

	default:
		fmt.Println("unexpected command")
		return "", false
	}
}
