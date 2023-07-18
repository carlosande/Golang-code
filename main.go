package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func main() {
	root := flag.String("root", ".", "Place to start searching")
	str := flag.String("search", "", "String to search")

	flag.Parse()

	if *str == "" {
		fmt.Println("Missing string to search")
		fmt.Println("Usage: ./main root_dir string_to_search")
		os.Exit(1)
	}

	run(*root, *str)

}

func run(root string, str string) {
	var wg sync.WaitGroup
	filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			ff := func(cpath <-chan string, str string, info os.FileInfo) {
				defer wg.Done()
				if !info.IsDir() {
					path := <-cpath
					b, err := ioutil.ReadFile(path)
					if err != nil {
						panic(err)
					}
					found, err := regexp.Match(str, b)
					if found {
						fmt.Printf("%s has been found in %s\n", str, path)
					}
					if err != nil {
						panic(err)
					}
				}
			}

			c1 := make(chan string, 1)
			c1 <- path

			wg.Add(1)
			go ff(c1, str, info)
			wg.Wait()

			return err
		})
}
