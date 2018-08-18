package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

var watcher *fsnotify.Watcher

func main() {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	//line := flag.Int("lines", 50, "an int")
	dir := flag.String("dir", "", "a filename")
	f := flag.String("file", "", "a directory")
	flag.Parse()

	fmt.Println("File: ", *f)
	fmt.Println("Dir:", *dir)

	var file *os.File
	var err error
	if *f != "" {
		file, err = os.Open(*f)
		if err != nil {
			fmt.Println("Error ", errors.Wrap(err, "opening file"))
			os.Exit(0)
		}
		defer file.Close()
	}
	// offset how many bytes to move. can be positive or negative
	var offset int64 = 5

	//whence is the point of reference of the offset
	// 0 - beginning of file
	// 1 - current position of file
	// 2 - end of file
	var whence = 0

	if *dir != "" {
		filepath.Walk(*dir, watchDir)
	}

	newPos, err := file.Seek(offset, whence)
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 1024)
	n, err := file.ReadAt(buf, 0)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(buf[:n]))

	fmt.Println("JUst moved to 5: ", newPos)

	//go back 2 bytes from the current position
	newPos, err = file.Seek(-2, 1)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Just moved back two: ", newPos)
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}

func connect(ip, user, pwd string) (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tmp", ip+":22", config)
	if err != nil {
		return nil, err
	}

	cs, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	return cs, nil
}
