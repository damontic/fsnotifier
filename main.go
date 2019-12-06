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

//
var watcher *fsnotify.Watcher

// main
func main() {

	directoriesUsage := "Mandatory list of directories to watch separated by ','"

	var directories string
	flag.StringVar(&directories, "directories", "", directoriesUsage)
	flag.StringVar(&directories, "d", "", directoriesUsage+" (shorthand)")
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

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// Adds directory list and subdirectories to the Watcher
	for _, dir := range directoryList {
		if err := filepath.Walk(dir, watchDir); err != nil {
			log.Fatalf("ERROR Walking directories:\n%s\n", err.Error())
		}
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Create:
					fmt.Printf("FIM CREATE %s\n", event.Name)
					addToWatchers(event)
				case fsnotify.Write:
					fmt.Printf("FIM WRITE %s\n", event.Name)
				case fsnotify.Remove:
					fmt.Printf("FIM REMOVE %s\n", event.Name)
					removeFromWatchers(event)
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

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func addToWatchers(event fsnotify.Event) {
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

func removeFromWatchers(event fsnotify.Event) {
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
