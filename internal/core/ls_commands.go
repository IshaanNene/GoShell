package core

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"
)

type FileInfoStruct struct {
	Fname  string
	Ftime  time.Time
	Fsize  int64
	Fgroup string
}

func checkError(err error, context string) {
	if err != nil {
		log.Fatalf("Error %s: %v", context, err)
	}
}

func listFiles(dir string, showHidden, appendSlashToDir, sortByTime, reverseOrder, sortBySize, recursive, listInode, showGroup, humanReadable, listDir bool) {
	var files []os.DirEntry
	var err error

	if listDir {
		files, err = os.ReadDir(filepath.Dir(dir))
		checkError(err, "reading directory")
	} else {
		files, err = os.ReadDir(dir)
		checkError(err, "reading directory")
	}

	var fileInfos []FileInfoStruct
	for _, file := range files {
		if !showHidden && file.Name()[0] == '.' {
			continue
		}

		info, err := file.Info()
		checkError(err, "getting file info")

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			log.Fatalf("Error retrieving file system information for %s", file.Name())
		}
		group := fmt.Sprintf("%d", stat.Gid)

		fileInfos = append(fileInfos, FileInfoStruct{
			Fname:  file.Name(),
			Ftime:  info.ModTime(),
			Fsize:  info.Size(),
			Fgroup: group,
		})
	}

	if sortByTime {
		sort.Slice(fileInfos, func(i, j int) bool {
			if reverseOrder {
				return fileInfos[i].Ftime.After(fileInfos[j].Ftime)
			}
			return fileInfos[i].Ftime.Before(fileInfos[j].Ftime)
		})
	} else if sortBySize {
		sort.Slice(fileInfos, func(i, j int) bool {
			if reverseOrder {
				return fileInfos[i].Fsize > fileInfos[j].Fsize
			}
			return fileInfos[i].Fsize < fileInfos[j].Fsize
		})
	} else if reverseOrder {
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].Fname > fileInfos[j].Fname
		})
	}

	for _, fileInfo := range fileInfos {
		name := fileInfo.Fname

		if listInode {
			info, err := os.Stat(filepath.Join(dir, fileInfo.Fname))
			checkError(err, "getting file stat")
			stat := info.Sys().(*syscall.Stat_t)
			fmt.Printf("%d ", stat.Ino)
		}

		if showGroup {
			fmt.Printf("%s ", fileInfo.Fgroup)
		}

		if appendSlashToDir && fileInfo.Fname[len(fileInfo.Fname)-1] != '/' {
			name += "/"
		}

		if humanReadable {
			fmt.Printf("%s %s\n", humanize.Bytes(uint64(fileInfo.Fsize)), name)
		} else {
			fmt.Println(name)
		}

		if recursive && fileInfo.Fname != "." && fileInfo.Fname != ".." {
			subDir := filepath.Join(dir, fileInfo.Fname)
			if fileInfo.Fname[len(fileInfo.Fname)-1] == '/' {
				listFiles(subDir, showHidden, appendSlashToDir, sortByTime, reverseOrder, sortBySize, recursive, listInode, showGroup, humanReadable, false)
			}
		}
	}
}

var LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List directory contents",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("directory")
		showHidden, _ := cmd.Flags().GetBool("a")
		appendSlashToDir, _ := cmd.Flags().GetBool("F")
		sortByTime, _ := cmd.Flags().GetBool("t")
		reverseOrder, _ := cmd.Flags().GetBool("r")
		sortBySize, _ := cmd.Flags().GetBool("S")
		recursive, _ := cmd.Flags().GetBool("R")
		listInode, _ := cmd.Flags().GetBool("i")
		showGroup, _ := cmd.Flags().GetBool("g")
		humanReadable, _ := cmd.Flags().GetBool("human-readable")
		listDir, _ := cmd.Flags().GetBool("list-dir")
		if dir == "" {
			dir = "."
		}

		listFiles(dir, showHidden, appendSlashToDir, sortByTime, reverseOrder, sortBySize, recursive, listInode, showGroup, humanReadable, listDir)
	},
}
func init() {
	LsCmd.Flags().StringP("directory", "D", ".", "Directory to list") 
	LsCmd.Flags().BoolP("a", "a", false, "Include hidden files") 
	LsCmd.Flags().BoolP("F", "F", false, "Append indicator (one of */=>@|) to entries") 
	LsCmd.Flags().BoolP("t", "t", false, "Sort by modification time, newest first") 
	LsCmd.Flags().BoolP("r", "r", false, "Reverse order while sorting") 
	LsCmd.Flags().BoolP("S", "S", false, "Sort by file size, largest first") 
	LsCmd.Flags().BoolP("R", "R", false, "List subdirectories recursively") 
	LsCmd.Flags().BoolP("i", "i", false, "Print the index number of each file") 
	LsCmd.Flags().BoolP("g", "g", false, "Display group ownership") 
	LsCmd.Flags().Bool("human-readable", false, "Print sizes in human-readable format") 
	LsCmd.Flags().BoolP("list-dir", "d", false, "List directories themselves, rather than their contents") 
}
