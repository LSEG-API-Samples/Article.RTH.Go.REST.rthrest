package trthrest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"sync"
	"time"
)

//HTTPPost : The function that wraps HTTP POST request. It adds the authorization token if token isn't nil
func HTTPPost(client *http.Client, url string, body *bytes.Buffer, headers map[string]string, trace bool) (*http.Response, error) {

	req, _ := http.NewRequest("POST", url, body)

	/*req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Prefer", "respond-async")
	if token != nil {
		req.Header.Add("Authorization", "Token "+*token)
	*/

	for key, value := range headers {
		//fmt.Printf("%s: %s\n", key, value)
		req.Header.Add(key, value)
	}
	if trace == true {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
	}

	resp, err := client.Do(req)

	if trace == true && err == nil {

		dumpBody := true
		contentLength, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		if contentLength > 5000 {
			dumpBody = false
		}

		dump, _ := httputil.DumpResponse(resp, dumpBody)
		log.Println(string(dump))
	}

	return resp, err
}

//HTTPGet : The function that wraps HTTP GET request. It adds the authorization token if token isn't nil
func HTTPGet(client *http.Client, url string, headers map[string]string, trace bool) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	/*
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Prefer", "respond-async")
		if token != nil {
			req.Header.Add("Authorization", "Token "+*token)

	*/
	for key, value := range headers {
		//fmt.Printf("%s: %s\n", key, value)
		req.Header.Add(key, value)
	}

	if trace == true {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(dump))
	}

	resp, err := client.Do(req)

	if trace == true && err == nil {
		dumpBody := true
		contentLength, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		if contentLength > 5000 {
			dumpBody = false
		}

		dump, _ := httputil.DumpResponse(resp, dumpBody)
		log.Println(string(dump))
	}

	return resp, err

}

//DownloadFile: Download the file by offset.
//if start == -1 means download full file
//if stop == -1 means download from start to the end of file
func DownloadFile(client *http.Client, headers map[string]string, url string, outFileName string, start int64, stop int64, tracing bool) {

	log.Printf("Download File: %s, %d, %d\n", outFileName, start, stop)
	var newHeaders map[string]string
	newHeaders = make(map[string]string)
	for k, v := range headers {
		newHeaders[k] = v
	}
	if start != -1 {
		//start == -1 means download full file
		if stop != -1 {
			newHeaders["Range"] = fmt.Sprintf("bytes=%d-%d", start, stop)

		} else {
			newHeaders["Range"] = fmt.Sprintf("bytes=%d-", start)
		}

	}
	resp, err := HTTPGet(client, url, newHeaders, tracing)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Status Code: %s\n%s ", resp.Status, string(body))
		//log.Fatalf("Status Code: %s\n ", resp.Status)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))

	if err != nil {
		log.Fatal(err)
	}

	done := make(chan int64)
	//outputFileName := "output_" + strconv.Itoa(os.Getpid()) + ".csv.gz"

	out, err := os.Create(outFileName)
	if err != nil {
		log.Fatal(err)
	}

	go PrintDownloadPercent(done, outFileName, int64(size))

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	done <- n
	resp.Body.Close()

}

//ConcurrentDownload: This function is used to download a file concurrently by the specified by the numOfConn
//Filesize of the file is required
func ConcurrentDownload(client *http.Client, headers map[string]string, url string, outFileName string, numOfConn int, fileSize int64, tracing bool) {
	var partSize, fileOffset int64
	partSize = fileSize / int64(numOfConn)
	fileOffset = 0

	log.Printf("ConcurrentDownload: %s, conn=%d\n", outFileName, numOfConn)
	var wg sync.WaitGroup

	for i := 1; i <= numOfConn; i++ {
		wg.Add(1)
		if i == numOfConn {
			log.Printf("Part %d: %d- \n", i, fileOffset)

			go func(filename string, start int64, stop int64) {
				defer wg.Done()
				DownloadFile(client, headers, url, filename, start, stop, tracing)
			}(fmt.Sprintf("part%d", i), fileOffset, -1)
		} else {
			log.Printf("Part %d: %d - %d\n", i, fileOffset, fileOffset+partSize-1)

			go func(filename string, start int64, stop int64) {
				defer wg.Done()
				DownloadFile(client, headers, url, filename, start, stop, tracing)
			}(fmt.Sprintf("part%d", i), fileOffset, fileOffset+partSize-1)
			fileOffset = fileOffset + partSize
		}
	}
	wg.Wait()
	MergeFile(numOfConn, outFileName)
}

//PrintDownloadPercent : This function shows the download progress
func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop = false
	var previousSize int64
	var maxRate, count, currentRate int

	var logInterval = 5
	previousSize, maxRate, count, currentRate = 0, 0, 0, 0

	for {
		select {
		case <-done:
			stop = true
		default:
			count = count + 1
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			currentRate = int(size - previousSize)
			if currentRate > maxRate {
				maxRate = currentRate
			}

			var percent float64
			percent = float64(size) / float64(total) * 100

			if count%logInterval == 0 {
				log.Printf("%s, Bytes: %d/Total: %d (%.0f%%)", path, size, total, percent)
			}

			previousSize = size

		}

		if stop {
			totalMB := total / 1024
			log.Printf("%s: Download Completed, Speed: Avg %.2f KB/s, Max %.2f KB/s", path, float32(totalMB)/float32(count), float32(maxRate)/float32(1024))
			break
		}

		time.Sleep(time.Second * 1)
	}
}

//MergeFile: Merge part1, part2, part3, ... files
func MergeFile(numberOfParts int, outFileName string) {
	log.Printf("Merging Files: %s\n", outFileName)
	b := make([]byte, 5000)
	destFile, _ := os.Create(outFileName)
	writer := bufio.NewWriter(destFile)
	for i := 1; i <= numberOfParts; i++ {
		filename := fmt.Sprintf("part%d", i)
		srcFile, _ := os.Open(filename)
		//fmt.Printf("Open File: %s\n", fmt.Sprintf("part%d", i))
		reader := bufio.NewReader(srcFile)
		readByte, err := reader.Read(b)
		//fmt.Printf("Read: %d bytes\n", readByte)
		for err != io.EOF || readByte > 0 {
			writer.Write(b[:readByte])
			//fmt.Printf("Write: %d bytes\n", writeByte)
			readByte, err = reader.Read(b)
			//fmt.Printf("Read: %d bytes\n", readByte)

		}
		srcFile.Close()
	}
	writer.Flush()
	destFile.Close()
}
