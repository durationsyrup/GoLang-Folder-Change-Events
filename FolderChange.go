package main

import (
	"fmt"
	"os"
	"time"
)

var (
	watchDir = "/home/seraph/GoLang"
)

func main() {
	//func handleChanges() handles the changes!
	if len(os.Args) > 1 {
		watchDir = os.Args[1]
	}
	watchChanges(watchDir, len(watchDir)+1)
}
func handleChanges(obj [][]string, areFiles bool) {
	Type := "File"
	if !areFiles {
		Type = "Folder"
	}
	for _, items := range obj {
		switch items[0] {
		case "add":
			fmt.Println(Type, "has been added (", items[1], ")")
		case "remove":
			fmt.Println(Type, "has been removed (", items[1], ")")
		case "update":
			fmt.Println(Type, "has been Updated (", items[1], ")")
		}
	}
}
func watchChanges(d string, cut int) {
	files, dirs := crawlDirectory(d, cut)
	if len(files) == 0 && len(dirs) == 0 {
		fmt.Println("Folder is empty or does not exist.")
	}
	for {
		f, dir := crawlDirectory(d, cut)

		anyFolderChanges, folderChanges := checkChange(dirs, dir, false)
		anyFileChanges, fileChanges := checkChange(files, f, true)

		if anyFolderChanges {
			dirs = dir
			handleChanges(folderChanges, false)
		}

		if anyFileChanges {
			files = f
			handleChanges(fileChanges, true)
		}

		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
}
func checkChange(a, b map[string]os.FileInfo, areFiles bool) (bool, [][]string) {
	// a should be old map
	// b should be new map
	difference := func(one, two map[string]os.FileInfo, c bool) [][]string {
		change := make([][]string, 0)
		for val := range one {
			if _, ok := two[val]; !ok {
				if c {
					change = append(change, []string{"remove", val})
				} else {
					change = append(change, []string{"add", val})
				}
			}
		}
		return change
	}

	if areFiles {
		if len(a) == len(b) {
			updates := make([][]string, 0)
			for name, val := range a {
				if b[name].ModTime() != val.ModTime() {
					updates = append(updates, []string{"update", name})
				}
			}
			if len(updates) != 0 {
				return true, updates
			}
		} else {
			if len(a) > len(b) || len(a) < len(b) {
				if len(a) > len(b) {
					// Files removed
					return true, difference(a, b, true)
				}
				//Files added
				return true, difference(b, a, false)
			}
		}
	} else {
		if len(a) > len(b) || len(a) < len(b) {
			if len(a) > len(b) {
				// Files removed
				return true, difference(a, b, true)
			}
			//Files added
			return true, difference(b, a, false)
		}
	}
	return false, [][]string{}
}

func crawlDirectory(directory string, cut int) (map[string]os.FileInfo, map[string]os.FileInfo) {
	dir, _ := os.Open(directory)
	defer dir.Close()
	files, _ := dir.Readdir(0)
	filenames := make(map[string]os.FileInfo, 0)
	dirs := make(map[string]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			directory = directory + "/" + file.Name()
			filen, filed := crawlDirectory(directory, cut)
			directory = directory[:len(directory)-(len("/"+file.Name()))]
			for d, found := range filed {
				dirs[d] = found
			}
			for d, found := range filen {
				filenames[file.Name()+"/"+d] = found
			}
			dirs[string(directory + "/" + file.Name())[cut:]] = file
		} else {
			filenames[file.Name()] = file
		}
	}
	return filenames, dirs
}
