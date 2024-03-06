// go run . --src="src.txt" --dst="test"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	srcFileUrls := flag.String("src", "DEFAULT VALUE", "Путь к файлу источнику")
	dirPath := flag.String("dst", "DEFAULT VALUE", "Путь к директории назначения")

	flag.Parse()

	if *dirPath == "DEFAULT VALUE" {
		fmt.Println("Введите параметры:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf(" --%s - %s\n", f.Name, f.Usage)
		})
		os.Exit(1)
	}

	// инициализация счетчика создаваемых файлов для перечисления, начиная с ../1.txt
	fileCounter := 1

	// чтение URLs из файла
	urls, err := openReadSourceFile(*srcFileUrls)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// создание директории по пути пользователя
	err = createDir(*dirPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// использование sync.waitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// увеличение счетчика горутин wg, создавая для обработки каждого url свою горутину
	for _, url := range urls {
		wg.Add(1)
		go processURL(url, *dirPath, &fileCounter, &wg)
	}

	// ожидание завершения всех горутин обработки URL в цикле for range
	wg.Wait()
	duration := time.Since(start)
	fmt.Println(duration)
}

// openReadSourceFile() []string открывает файл по указанному в srcFileUrls пути и возвращает массив []string URL'ов:
// если случается ошибка (error), возвращает ее
func openReadSourceFile(srcFileUrls string) ([]string, error) {
	file, err := os.Open(srcFileUrls)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string

	// для универсального подхода к разным ОС используется scanner
	scanner := bufio.NewScanner(file)

	// сканирование каждой строчки и добавление в urls
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
		if scanner.Err() != nil {
			fmt.Println(err)
			return nil, scanner.Err()
		}
	}

	return urls, nil
}

// createDir() создает директорию по указанному в dirPath пути:
// если случается ошибка (error), возвращает ее
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

// processURL() уменьшает счетчик горутин wg на -1, отправляет GET-запрос по url и проверяет получен ли корректный response:
// если получен корректный ответ, вызывает функцию создания файла createFile() в директории dirPath;
// если ответ не получен или ответ некорректный, печатает сообщение об ошибке
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

	err = createFile(url, dirPath, *fileCounter, response)
	if err != nil {
		return err
	}

	*fileCounter++

	return nil
}

// createFile() создает файл по пути [dirPath]/[fileCounter].txt с содержимым из response.Body, а также увеличивает счетчик fileCounter, если не было ошибок
// если при создании или чтении response.Body случается ошибка (error), то выводит сообщение об ошибке
func createFile(url, dirPath string, fileCounter int, response *http.Response) error {
	filePath := fmt.Sprintf("%s/%s.txt", dirPath, strconv.Itoa(fileCounter))
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v. URL: %s\n", err, url)
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, response.Body)
	if err != nil {
		fmt.Printf("Ошибка копирования тела ответа в файл: %v из URL: %s\n", err, url)
		return err
	}

	fmt.Printf("Файл с содержимым создан по пути: %s. URL: %s\n", filePath, url)

	return nil
}
