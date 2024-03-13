// go run . --src="src.txt" --dst="test"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	// парсинг флагов терминала
	srcFilePath, dirPath, err := parseFlags()
	if err != nil {
		os.Exit(1)
	}

	// инициализация счетчика создаваемых файлов для перечисления, начиная с ../1.txt
	fileCounter := 1

	// чтение и получение URLs из файла
	urls, err := openReadSourceFile(srcFilePath)
	if err != nil {
		os.Exit(2)
	}

	// создание директории по заданному пути
	err = createDir(dirPath)
	if err != nil {
		os.Exit(3)
	}

	// использование sync.waitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// увеличение счетчика горутин wg, создание для обработки каждого url своей горутины
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			processURL(url, dirPath, &fileCounter)
			wg.Done()
		}(url)
	}

	wg.Wait()

	duration := time.Since(start).Seconds()
	fmt.Printf("Длительность выполнения в секундах: %f\n", duration)
}

// parseFlags парсит флаги терминала и возвращает их
func parseFlags() (string, string, error) {
	var srcFilePath string
	flag.StringVar(&srcFilePath, "src", "DEFAULT VALUE", "путь к файлу источнику")
	var dirPath string
	flag.StringVar(&dirPath, "dst", "DEFAULT VALUE", "путь к директории назначения")

	flag.Parse()

	if srcFilePath == "DEFAULT VALUE" || dirPath == "DEFAULT VALUE" {
		flag.Usage()
	}

	return srcFilePath, dirPath, nil
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

// createDir создает директорию по указанному в dirPath пути, если ее не существует
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
func processURL(url, dirPath string, fileCounter *int) error {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ответ не получен. Ошибка: %s. URL: %s\n", err, url)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Ответ не получен. Ошибка: некорректный адрес. URL: %s\n", url)
		err = fmt.Errorf(errMsg)
		fmt.Println(errMsg)
		return err
	}

	// чтение байтов из response.Body
	fileBytes, err := io.ReadAll(response.Body)
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
	creationTime := time.Now().Format("02-01-2006__15-04-05")
	fileName := fmt.Sprintf("%d__%s.txt", fileCounter, creationTime)
	path := filepath.Join(dirPath, fileName)

	err := os.WriteFile(path, fileBytes, fs.ModePerm)
	if err != nil {
		fmt.Printf("Ошибка копирования тела ответа в файл: %v из URL: %s\n", err, url)
		return err
	}

	fmt.Printf("Файл создан по пути: %s. URL: %s\n", path, url)

	return nil
}
