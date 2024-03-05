package main

// go run . --src="../src.txt" --dst="../test"

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	srcPath := flag.String("src", "", "Soure file path")
	inputPath := flag.String("dst", "", "Destination dir")
	fileCounter := 0

	flag.Parse()

	if err := os.Mkdir(*inputPath, os.ModePerm); err != nil {
		fmt.Println(err)
	}

	srcFile, err := os.Open(*srcPath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(srcFile)
	fileScanner.Split(bufio.ScanLines)
	var urls []string

	for fileScanner.Scan() {
		urls = append(urls, fileScanner.Text())
	}
	srcFile.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer srcFile.Close()

	for i := 0; i < len(urls); i++ {
		response, err := http.Get(urls[i])
		if err != nil {
			fmt.Println("Response wasn't grabbed: Wrong URl scheme of", urls[i])
			continue
		}
		defer response.Body.Close()

		if response.Status == "200 OK" {
			fileCounter++
			path := *inputPath + "/" + fmt.Sprint(fileCounter) + ".txt"
			outputFile, err := os.Create(path)

			if err != nil {
				fmt.Println("Unable to create file:", err)
				os.Exit(1)
			}
			defer outputFile.Close()

			scanner := bufio.NewScanner(response.Body)
			for i := 0; scanner.Scan(); i++ {
				outputFile.WriteString(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}
			fmt.Println("Response was grabbed:", urls[i])
		} else {
			fmt.Println("Response wasn't grabbed: Wrong response status of", urls[i])
		}
	}
}
