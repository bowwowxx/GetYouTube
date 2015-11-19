package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type PassThru struct {
	io.Reader
	total    int64
	length   int64
	progress float64
}

func getRemoteFileSize(url string) int {
	responseHeader, err := http.Head(url)
	if err != nil {
		log.Fatal("Couldn't retreive source file header:", err)
	}
	srcFileSize, err := strconv.Atoi(responseHeader.Header.Get("Content-Length"))
	if err != nil {
		log.Fatal("Couldn't convert the Content-Length to integer:", err)
	}
	return srcFileSize
}

func printDownloadProgress(total int64, length int64) {
	var maxProgress int = 20
	var progressBar string = ""
	var percentage int = int((total * 100) / length)

	// Calculate and create the progress bar
	var progress = (percentage * maxProgress) / 100
	for i := 0; i < progress; i++ {
		progressBar = progressBar + "="
	}
	for i := 0; i < (maxProgress - progress); i++ {
		progressBar = progressBar + " "
	}

	fmt.Printf("\t[%s] %d %d/%d", progressBar, percentage, total, length)
	if percentage < 100 {
		fmt.Print("\r")
	} else {
		fmt.Print("\n")
	}
}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		printDownloadProgress(pt.total, pt.length)
	}

	return n, err
}

func DownloadFile(url string, fileName string) {
	log.Println("Downloading", fileName)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading -", err)
		return
	}
	defer response.Body.Close()

	readerpt := &PassThru{Reader: response.Body, length: response.ContentLength}
	body, err := ioutil.ReadAll(readerpt)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(fileName, body, 0644)
}
