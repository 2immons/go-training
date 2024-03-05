// go run . --src="../src.txt" --dst="../test"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	// считывание файла в отдельную функцию
	srcPath := flag.String("src", "", "Source file path")
	inputPath := flag.String("dst", "", "Destination dir")
	flag.Parse()

	fileCounter := 0

	if err := os.Mkdir(*inputPath, os.ModePerm); err != nil {
		fmt.Println(err)
	}

	srcFile, err := os.Open(*srcPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer srcFile.Close()

	var urls []string
	fileScanner := bufio.NewScanner(srcFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		urls = append(urls, fileScanner.Text())
	}

	// используем waitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup
	// используем канал для передачи результатов обработки URL в горутине

	for i := 0; i < len(urls); i++ {
		// с каждым новым urls[i] увеличиваем счетчик wg горутин и запускаем отдельную горутину для обработки urls[i]
		wg.Add(1)
		go processURL(urls[i], *inputPath, &fileCounter, &wg)
	}

	// ждем завершения всех горутин обработки URL в цикле for (т.е. пока wg счетчик не станет равным 0)
	wg.Wait()
}

// функция обработки URL
func processURL(url, inputPath string, counter *int, wg *sync.WaitGroup) {
	// уменьшаем счетчик горутин на -1 после отработки очередной
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Response wasn't grabbed: Wrong URL scheme of %s\n", url)
		return
	}
	defer response.Body.Close()

	if response.Status == "200 OK" {
		*counter++
		path := inputPath + "/" + fmt.Sprint(*counter) + ".txt"
		outputFile, err := os.Create(path)

		if err != nil {
			fmt.Printf("Unable to create file: %v\n", err)
			return
		}
		defer outputFile.Close()

		scanner := bufio.NewScanner(response.Body)
		for i := 0; scanner.Scan(); i++ {
			outputFile.WriteString(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		fmt.Printf("Response was grabbed: %s\n", url)
	} else {
		fmt.Printf("Response wasn't grabbed: Wrong response status of %s\n", url)
	}
}
