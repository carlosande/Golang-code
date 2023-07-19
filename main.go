package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

func main() {
	root := flag.String("root", ".", "Place to start searching")
	str := flag.String("search", "", "String to search")
	// *str = "regexp" // for testing
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

	ff := func(cpath <-chan string, str string) {
		defer wg.Done()
		path := <-cpath
		b, err := os.ReadFile(path)
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

	files, err := os.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	c1 := make(chan string, 1)

	for _, file := range files {
		if !file.IsDir() {
			c1 <- file.Name()
			wg.Add(1)
			go ff(c1, str)
			wg.Wait()
		}
	}
}
