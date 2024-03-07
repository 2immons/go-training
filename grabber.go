// go run . --src="src.txt" --dst="test"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	// парсинг флагов терминала
	srcFilePath, dirPath := parseFlags()

	// инициализация счетчика создаваемых файлов для перечисления, начиная с ../1.txt
	fileCounter := 1

	// чтение и получение URLs из файла
	urls, err := openReadSourceFile(*srcFilePath)
	if err != nil {
		os.Exit(2)
	}

	// создание директории по заданному пути
	err = createDir(*dirPath)
	if err != nil {
		os.Exit(3)
	}

	// использование sync.waitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// увеличение счетчика горутин wg, создание для обработки каждого url своей горутины
	for _, url := range urls {
		wg.Add(1)
		go processURL(url, *dirPath, &fileCounter, &wg)
	}

	wg.Wait()

	duration := time.Since(start).Seconds()
	fmt.Printf("Длительность выполнения в секундах: %f\n", duration)
}

// parseFlags парсит флаги терминала и возвращает их
func parseFlags() (*string, *string) {
	srcFilePath := flag.String("src", "DEFAULT VALUE", "Путь к файлу источнику")
	dirPath := flag.String("dst", "DEFAULT VALUE", "Путь к директории назначения")

	flag.Parse()

	if *dirPath == "DEFAULT VALUE" {
		flag.VisitAll(func(f *flag.Flag) {})
	}

	return srcFilePath, dirPath
}

// openReadSourceFile по указанному в srcFileUrls пути и возвращает массив []string URL'ов:
func openReadSourceFile(srcFileUrls string) ([]string, error) {
	// открытие файла
	file, err := os.Open(srcFileUrls)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string

	// для универсального подхода к разным ОС для чтения используется scanner
	scanner := bufio.NewScanner(file)

	// сканирование построчно файла и добавление в urls
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
		if scanner.Err() != nil {
			fmt.Println(err)
			return nil, scanner.Err()
		}
	}

	return urls, nil
}

// createDir создает директорию по указанному в dirPath пути
func createDir(dirPath string) error {
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// processURL проверяет получен ли корректный response и вызывает функцию создания файла createFile() в директории dirPath
func processURL(url, dirPath string, fileCounter *int, wg *sync.WaitGroup) error {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil || response.Status != "200 OK" {
		if err != nil {
			fmt.Printf("Ответ не получен. Некорректный формат URL. Ошибка: %s. URL: %s\n", err, url)
			return err
		} else {
			fmt.Printf("Ответ не получен. Некорректный статус (не 200 OK). URL: %s\n", url)
			return nil
		}
	}
	defer response.Body.Close()

	// чтение байтов из response.Body
	fileBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// вызов функции создания файла и записи в него байтов
	err = createFile(url, dirPath, *fileCounter, fileBytes)
	if err != nil {
		return err
	}

	*fileCounter++

	return nil
}

// createFile создает файл по и записывает в него срез байтов
func createFile(url, dirPath string, fileCounter int, fileBytes []byte) error {
	path := fmt.Sprintf("%s/%d.txt", dirPath, fileCounter)

	err := ioutil.WriteFile(path, fileBytes, fs.ModePerm)
	if err != nil {
		fmt.Printf("Ошибка копирования тела ответа в файл: %v из URL: %s\n", err, url)
		return err
	}

	fmt.Printf("Файл с содержимым создан по пути: %s. URL: %s\n", path, url)

	return nil
}
