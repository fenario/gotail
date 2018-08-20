package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	//"time"

	"golang.org/x/crypto/ssh"

	"github.com/fsnotify/fsnotify"
)

var (
	// offset how many bytes to move. can be positive or negative
	offset int64 = -1024
	//whence is the point of reference of the offset
	// 0 - beginning of file
	// 1 - current position of file
	// 2 - end of file
	whence  = 2
	watcher *fsnotify.Watcher
)

func main() {
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	//line := flag.Int("lines", 50, "an int")
	dir := flag.String("dir", "", "A directory that contains files")
	f := flag.String("file", "", "a file for reading")
	flag.Parse()

	fmt.Println("File: ", *f)
	fmt.Println("Dir:", *dir)

	done := make(chan bool)
	go func() {
		for {
			//time.Sleep(2 * time.Second)
			select {
			case event := <-watcher.Events:
				fmt.Println(event.Name)
				b, _ := readFile(event.Name)
				fmt.Println(string(b))

			case err := <-watcher.Errors:
				fmt.Println(err)
			}
		}
	}()

	if *dir != "" {
		if err := filepath.Walk(*dir, watchDir); err != nil {
			fmt.Println("Walk dir error: ", err)
		}
	}

	if *f != "" {
		err := watcher.Add(*f)
		if err != nil {
			fmt.Println(err)
		}
	}
	<-done
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}

func readFile(f string) ([]byte, error) {
	var (
		file     *os.File
		err      error
		stat     os.FileInfo
		numBytes int
	)
	buf := make([]byte, 1022)

	file, err = os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err = file.Stat()
	if err != nil {
		return nil, err
	}

	newPos, err := file.Seek(offset, whence)
	fmt.Println("Len:", stat.Size())
	fmt.Println("NewPos:", newPos)
	if err != nil {
		return nil, err
	}

	numBytes, err = file.ReadAt(buf, newPos)
	if err != nil {
		return nil, err
	}
	return buf[:numBytes], nil
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
