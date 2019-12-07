package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// main
func main() {

	directoriesUsage := "Mandatory list of directories to watch separated by ','"
	excludesUsage := "List of directories to exclude separated by ','\nExclusion is executed when directory to be added starts with the excluded name, it doesnt have to match exactly."

	var directories string
	flag.StringVar(&directories, "directories", "", directoriesUsage)
	flag.StringVar(&directories, "d", "", directoriesUsage+" (shorthand)")

	var excludes string
	flag.StringVar(&excludes, "excludes", "", excludesUsage)
	flag.StringVar(&excludes, "e", "", excludesUsage+" (shorthand)")

	flag.Parse()

	if directories == "" {
		log.Fatalf("The list of directories is mandatory.")
	}

	directoryList := strings.Split(directories, ",")
	// checks if directories are valid
	for _, dir := range directoryList {
		fileInfo, err := os.Lstat(dir)
		if err != nil {
			log.Fatalf("Error checking initial directories\n%s\n", err.Error())
		}
		if !fileInfo.IsDir() {
			log.Fatalf("The specified input %s is not a directory!\n", dir)
		}
	}

	excludesList := strings.Split(excludes, ",")

	// creates a new file watcher
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	// Adds directory list and subdirectories to the Watcher excluding excludes
	initializeWatcher(watcher, directoryList, excludesList)

	done := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Create:
					fmt.Printf("FIM CREATE %s\n", event.Name)
					if !isDirectoryExcluded(event.Name, excludesList) {
						addToWatchers(event, watcher)
					}
				case fsnotify.Write:
					fmt.Printf("FIM WRITE %s\n", event.Name)
				case fsnotify.Remove:
					fmt.Printf("FIM REMOVE %s\n", event.Name)
					removeFromWatchers(event, watcher)
				case fsnotify.Rename:
					fmt.Printf("FIM RENAME %s\n", event.Name)
				case fsnotify.Chmod:
					fmt.Printf("FIM CHMOD %s\n", event.Name)
				}

			// watch for errors
			case err := <-watcher.Errors:
				fmt.Printf("ERROR received from Watcher\n%s\n", err.Error())
			}
		}
	}()

	<-done
}

func addToWatchers(event fsnotify.Event, watcher *fsnotify.Watcher) {
	fileInfo, err := os.Lstat(event.Name)
	if err != nil {
		log.Fatalf("ERROR at lstat of file:\n%s\n%v\n", event.Name, err.Error())
	}

	if fileInfo.IsDir() {
		err = watcher.Add(event.Name)
		if err != nil {
			log.Fatalf("ERROR adding new directory to watch list\n%v\n", err.Error())
		}
	}
}

func removeFromWatchers(event fsnotify.Event, watcher *fsnotify.Watcher) {
	fileInfo, err := os.Lstat(event.Name)
	if err != nil {
		log.Fatalf("ERROR at lstat of file:\n%s\n%v\n", event.Name, err.Error())
	}

	if fileInfo.IsDir() {
		err = watcher.Remove(event.Name)
		if err != nil {
			log.Fatalf("ERROR removing directory from watchers list\n%v\n", err.Error())
		}
	}
}

func isDirectoryExcluded(absolutePath string, excludedList []string) bool {
	isDirectoryExcluded := false
	for _, excludedDirectory := range excludedList {
		if strings.HasPrefix(absolutePath, excludedDirectory) {
			isDirectoryExcluded = true
			break
		}
	}
	return isDirectoryExcluded
}

func initializeWatcher(watcher *fsnotify.Watcher, directoryList, excludesList []string) {

	var watchDirFunc = func(path string, fi os.FileInfo, err error) error {
		// since fsnotify can watch all the files in a directory, watchers only need
		// to be added to each nested directory
		if fi.Mode().IsDir() && !isDirectoryExcluded(path, excludesList) {
			return watcher.Add(path)
		}
		return nil
	}

	for _, dir := range directoryList {
		if err := filepath.Walk(dir, watchDirFunc); err != nil {
			log.Fatalf("ERROR Walking directories:\n%s\n", err.Error())
		}
	}
}
