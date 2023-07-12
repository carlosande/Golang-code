package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
)

func main() {

	path, err := os.Getwd()
	if err == nil {
		dirEntries, err := os.ReadDir(path)
		done := make(map[os.DirEntry]bool, len(dirEntries))

		if err == nil {
			for i := 0; i < runtime.NumCPU(); i++ {
				fmt.Printf("running CPU %d\n", i+1)
				for _, dentry := range dirEntries {
					if !done[dentry] {
						done[dentry] = true
						go ff(os.Args[1], dentry.Name())
						//time.Sleep(time.Microsecond * 500)
					}

				}
			}
		}

	} else {
		fmt.Printf("Error %v", err.Error())
	}

}

func ff(str, filepath string) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	found, err := regexp.Match(str, b)
	if found {
		fmt.Printf("%s has been found in %s\n", str, filepath)
	} else {
		fmt.Printf("%s has NOT been found in %s\n", str, filepath)
	}
	if err != nil {
		panic(err)
	}
}
