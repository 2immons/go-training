// go run . --src="src.txt" --dst="test"

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {
	// инициализируем счетчик создаваемых файлов для перечисления, начиная с ../1.txt
	fileCounter := 1

	// парсинг флагов
	srcFileUrls := flag.String("src", "", "Путь к файлу источнику")
	dirPath := flag.String("dst", "", "Путь к директории назначения")
	flag.Parse()

	// открываем файл, читаем его в urls []string
	urls := openAndReadFile(*srcFileUrls)

	// создаем директорию по пути пользователя
	createDir(*dirPath)

	// используем sync.waitGroup для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// увеличиваем счетчик горутин wg, создавая для обработки каждого url свою горутину
	for _, url := range urls {
		wg.Add(1)
		go processURL(url, *dirPath, &fileCounter, &wg)
	}

	// ждем завершения всех горутин обработки URL в цикле for range
	wg.Wait()
}

// checkFatalError() проверяет содержится ли фатальная ошибка и завершает программу, если содержится
func checkFatalError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// openAndReadFile() []string открывает файл по указанному в srcFileUrls пути и возвращает массив []string URL'ов:
// если случается ошибка (error), завершает программу с помощью функции checkFatalError()
func openAndReadFile(srcFileUrls string) []string {
	file, err := os.Open(srcFileUrls)
	checkFatalError(err)
	defer file.Close()

	var urls []string

	// для универсального подхода к разным ОС используем scanner и считываем строки в urls[]
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	checkFatalError(scanner.Err())

	return urls

	// content, err := os.ReadFile(srcFileUrls)
	// checkFatalError(err)

	// // Разделение содержимого файла на строки
	// urls := strings.Split(string(content), "\n")

	// return urls
}

// createDir() создает директорию по указанному в dirPath пути:
// если случается ошибка (error), завершает программу с помощью функции checkFatalError()
func createDir(dirPath string) {
	err := os.Mkdir(dirPath, os.ModePerm)
	checkFatalError(err)
}

// processURL() уменьшает счетчик горутин wg на -1, отправляет GET-запрос по url и проверяет получен ли корректный response:
// если получен корректный ответ, вызывает функцию создания файла createFile() в директории dirPath;
// если ответ не получен или ответ некорректный, печатает сообщение об ошибке
func processURL(url, dirPath string, fileCounter *int, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil || response.Status != "200 OK" {
		fmt.Printf("Ответ не получен или некорректный формат URL. URL: %s\n", url)
		return
	}

	// ВОПРОС на будущее:
	// Мы обсуждали, что манипуляции с https и дескриптерами должны быть вынесены в отдельную функцию, чтобы https и дескриптер могли закрыться,
	// но ведь у меня происходит return на 79 строчке в случае ошибки и тогда defer Close() должен сработать, разве нет?
	// Или же я не так понял?
	defer response.Body.Close()

	createFile(url, dirPath, fileCounter, response)
}

// createFile() создает файл по пути [dirPath]/[fileCounter].txt с содержимым из response.Body, а также увеличивает счетчик fileCounter, если не было ошибок
// если при создании или чтении response.Body случается ошибка (error), то выводит сообщение об ошибке
func createFile(url, dirPath string, fileCounter *int, response *http.Response) {
	filePath := fmt.Sprintf("%s/%s.txt", dirPath, fmt.Sprint(*fileCounter))
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v из URL: %s\n", err, url)
		return
	}

	// такой же ВОПРОС на будущее, что и с https, та же ситуация:
	defer newFile.Close()

	*fileCounter++

	_, err = io.Copy(newFile, response.Body)
	if err != nil {
		fmt.Printf("Ошибка копирования тела ответа в файл: %v из URL: %s\n", err, url)
		return
	}

	fmt.Printf("Ответ получен. Файл с содержимым создан по пути: %s. URL: %s\n", filePath, url)
}
