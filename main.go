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

	directories := flag.String("directories", "", "Mandatory list of directories to watch separated by ','")
	flag.Parse()

	if *directories == "" {
		log.Fatalf("The list of directories is mandatory. ")
	}

	// creates a new file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// Adds directory list and subdirectories to the Watcher
	directoryList := strings.Split(*directories, ",")
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
				fmt.Printf("FIM EVENT! %#v\n", event)
				if event.Op == fsnotify.Create {
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
