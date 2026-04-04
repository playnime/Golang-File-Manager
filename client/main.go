package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func UploadFile(login string, password string, filename string) (*http.Response, error) {
	//создание "коробки" и структуры чтоб в нее писать
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//далее замена на jwt токен
	writer.WriteField("Login", login)
	writer.WriteField("Password", password)

	//занесение файла в форму
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error on opening file: ", err)
		return nil, err
	}
	defer f.Close()

	path, err := writer.CreateFormFile("File", filename)
	if err != nil {
		fmt.Println("error on creating form file: ", err)
		return nil, err
	}

	_, err = io.Copy(path, f)
	if err != nil {
		fmt.Println("error on copying file in form: ", err)
		return nil, err
	}

	//заклеить коробку
	err = writer.Close()
	if err != nil {
		fmt.Println("error on closing writer (box): ", err)
		return nil, err
	}

	//отправка

	req, err := http.NewRequest("POST", "http://localhost:9091/add", body)
	if err != nil {
		fmt.Println("error on creating request: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error on sending request: ", err)
		return nil, err
	}

	return r, nil
}

func main() {
	// мало ли Илья Клаленик захочет посмотреть что этот клиент делает
	//просто выполняет один запрос на добавление файла (чтобы запрос выполнился, файл должен лежать возле main.go)
	r, err := UploadFile("Vanya", "4321", "test2.txt")
	if err != nil {
		fmt.Println("main: error on qyery: ", err)
	}

	_ = r
}
