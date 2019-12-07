# fsnotifier
Process that keeps a look at the inotify events and prints when a file system changes. Uses the github.com/fsnotify/fsnotify library.

```bash
Usage of fsnotifier:
  -d string
    	Mandatory list of directories to watch separated by ',' (shorthand)
  -directories string
    	Mandatory list of directories to watch separated by ','
  -e string
    	List of directories to exclude separated by ','
    	Exclusion is executed when directory to be added starts with the excluded name, it doesnt have to match exactly. (shorthand)
  -excludes string
    	List of directories to exclude separated by ','
    	Exclusion is executed when directory to be added starts with the excluded name, it doesnt have to match exactly.
```
