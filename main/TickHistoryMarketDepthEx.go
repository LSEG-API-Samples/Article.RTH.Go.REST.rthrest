package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/TR-API-Samples/Article.TRTH.Go.REST.trthrest"
	"github.com/howeyc/gopass"
)

//Enter username and password here
var dssUserName = ""
var dssPassword = ""
var trthURL = "https://hosted.datascopeapi.reuters.com/RestApi/v1/"

//GetExtractionIDFromNote : Get Extraction ID number from note in the response
func GetExtractionIDFromNote(note string) string {
	extractionIDReg := regexp.MustCompile("Extraction ID: ([0-9]+)")
	IDReg := regexp.MustCompile("[0-9]+")
	return IDReg.FindString(extractionIDReg.FindString(note))

}

func main() {
	//heaers is map used to store HTTP headers for the request
	var headers map[string]string
	var outputFilename string
	var fileSize int64
	var step = 0

	//All available arguments of the example
	directDownloadFlag := flag.Bool("aws", false, "Download from AWS (false)")
	numOfConnection := flag.Int("n", 1, "Number of concurent download channels")
	traceFlag := flag.Bool("X", false, "Enable HTTP tracing (false)")
	username := flag.String("u", "", "DSS Username ('')")
	password := flag.String("p", "", "DSS Password ('')")
	proxy := flag.String("proxy", "", "Proxy: http://username:password@proxy:port")
	flag.Parse()

	dssUserName = *username
	dssPassword = *password

	//Print the values in the arguments for verification
	if *directDownloadFlag == true {
		log.Printf("X-Direct-Download: true \n")
	}
	if *traceFlag == true {
		log.Printf("Tracing: true \n")
	}
	log.Printf("Number of concurrent download: %d\n", *numOfConnection)

	//Create and set common headers of the HTTP request
	headers = make(map[string]string)

	headers["Content-Type"] = "application/json"
	headers["Prefer"] = "respond-async"

	//Prepare the TickHistoryMarketDepthExtractionRequest
	request := new(trthrest.TickHistoryMarketDepthExtractionRequest)
	request.Condition.View = trthrest.ViewOptionsNormalizedLL2Enum
	request.Condition.SortBy = trthrest.SortSingleByRicEnum
	request.Condition.NumberOfLevels = 10
	request.Condition.MessageTimeStampIn = trthrest.TimeOptionsGmtUtcEnum
	request.Condition.DisplaySourceRIC = true
	request.Condition.ReportDateRangeType = trthrest.ReportDateRangeTypeRangeEnum
	startdate := time.Date(2017, 7, 1, 0, 0, 0, 0, time.UTC)
	request.Condition.QueryStartDate = &startdate
	enddate := time.Date(2017, 8, 23, 0, 0, 0, 0, time.UTC)
	request.Condition.QueryEndDate = &enddate
	request.ContentFieldNames = []string{
		"Ask Price",
		"Ask Size",
		"Bid Price",
		"Bid Size",
		"Domain",
		"History End",
		"History Start",
		"Instrument ID",
		"Instrument ID Type",
		"Number of Buyers",
		"Number of Sellers",
		"Sample Data",
	}

	request.IdentifierList.InstrumentIdentifiers = append(request.IdentifierList.InstrumentIdentifiers, trthrest.InstrumentIdentifier{Identifier: "IBM.N", IdentifierType: "Ric"})
	request.IdentifierList.ValidationOptions = &trthrest.InstrumentValidationOptions{AllowHistoricalInstruments: true}

	//Define the HTTP transport and client used by the example
	var tr http.Transport
	if *proxy == "" {
		tr = http.Transport{
			DisableCompression: true}
	} else {
		proxyURL, _ := url.Parse(*proxy)
		tr = http.Transport{
			DisableCompression: true,
			Proxy:              http.ProxyURL(proxyURL),
		}
	}

	client := &http.Client{
		Transport: &tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	//If the username and password are not specified in the arguments, the example will request from the users.
	if dssUserName == "" {
		fmt.Print("Enter DSS Username: ")
		fmt.Scanln(&dssUserName)
	}
	if dssPassword == "" {
		fmt.Print("Enter DSS Password: ")
		temp, _ := gopass.GetPasswdMasked()
		dssPassword = string(temp)
	}

	//Create JSON byte array for the token request
	loginreq, err := json.Marshal(struct {
		Credentials trthrest.Credential
	}{
		Credentials: trthrest.Credential{
			Username: dssUserName,
			Password: dssPassword,
		},
	})

	step++
	log.Printf("Step %d: RequestToken\n", step)

	//Request to get the token
	resp, err := trthrest.HTTPPost(client, trthrest.GetRequestTokenURL(trthURL), bytes.NewBuffer(loginreq), headers, *traceFlag)

	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalf("Status Code: %s\n%s ", resp.Status, string(body))

	}

	//Process the token in the reponse
	var tokentResponse = &trthrest.RequestTokenResponse{}

	err = json.Unmarshal(body, tokentResponse)
	resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	//Add the Authorization header with the retreived token
	token := tokentResponse.Value
	headers["Authorization"] = "Token " + token

	//Prepare JSON object for TickHistoryMarketDepthExtractionRequest
	req1, _ := json.Marshal(struct {
		ExtractionRequest *trthrest.TickHistoryMarketDepthExtractionRequest
	}{
		ExtractionRequest: request,
	})
	step++
	log.Printf("Step %d: ExtractRaw for TickHistoryMarketDepthExtractionRequest\n", step)

	//Send the TickHistoryMarketDepthExtractionRequest to ExtractRaw endpoint
	resp, err = trthrest.HTTPPost(client, trthrest.GetExtractRawURL(trthURL), bytes.NewBuffer(req1), headers, *traceFlag)

	if err != nil {
		log.Fatal(err)
	}

	//Check the status of the extraction
	var statusCount = 0
	for resp.StatusCode == 202 {
		time.Sleep(3000 * time.Millisecond)
		statusCount++
		location := resp.Header.Get("Location")
		//Change the protocol to https if it is http
		location = strings.Replace(location, "http:", "https:", 1)
		if statusCount == 1 {
			step++
		}
		log.Printf("Step %d: Checking Status (%d) of Extraction (%d)\n", step, resp.StatusCode, statusCount)
		resp, err = trthrest.HTTPGet(client, location, headers, *traceFlag)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Status Code: %s\n%s ", resp.Status, string(body))
	}

	//Process in the extraction response
	extractRawResult := &trthrest.RawExtractionResult{}
	err = json.Unmarshal(body, extractRawResult)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	//if the client uses concurrent downloads (n > 1), the example will get the extraction ID from the notes,
	//and then send a request to get the filename and filesize
	if *numOfConnection > 1 {
		extractionID := GetExtractionIDFromNote(extractRawResult.Notes[0])

		log.Printf("ExtractionID: %q\n", extractionID)
		//If there is no extraction ID in the notes, the concurrent download will be diable
		if extractionID == "" {
			log.Println("ExtractionID is nil: Disable Concurrent Download")
			*numOfConnection = 1
			outputFilename = fmt.Sprintf("output_%s.csv.gz", extractRawResult.JobID)
			fileSize = 0
		}

		if extractionID != "" {
			step++
			log.Printf("Step %d: Get File information\n", step)
			resp, err = trthrest.HTTPGet(client, trthrest.GetReportExtractionFullFileURL(trthURL, extractionID), headers, *traceFlag)
			if err != nil {
				log.Fatal(err)
			}
			body, err = ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != 200 {

				log.Fatalf("Status Code: %s\n%s ", resp.Status, string(body))
			}
			extractedFile := &trthrest.ExtractedFile{}
			err = json.Unmarshal(body, extractedFile)
			if err != nil {
				log.Fatal(err)
			}

			outputFilename = extractedFile.ExtractedFileName
			fileSize = extractedFile.Size

		}
	} else {
		outputFilename = fmt.Sprintf("output_%s.csv.gz", extractRawResult.JobID)
		fileSize = 0
	}
	log.Printf("File: %s, Size: %d\n", outputFilename, fileSize)

	//Set the download url to Extractions/RawExtractionResults('{{jobId}}')/$value
	downloadURL := trthrest.GetRawExtractionResultGetDefaultStreamURL(trthURL, extractRawResult.JobID)
	//Set the time to measure the download time
	start := time.Now()
	//If -aws is set, the application will download the result file from aws
	if *directDownloadFlag == true {

		//Clone the TRTH headers to newHeaders and then add X-Direct-Download to the new header
		newHeaders := make(map[string]string)
		for k, v := range headers {
			newHeaders[k] = v
		}
		newHeaders["X-Direct-Download"] = "true"
		step++
		log.Printf("Step %d: Get AWS URL\n", step)
		resp, err = trthrest.HTTPGet(client, downloadURL, newHeaders, *traceFlag)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode == 302 {
			//GET AWS URL used to download a file and change the URL in downloadURL variable
			downloadURL = resp.Header.Get("Location")
			log.Printf("AWS: %s\n", downloadURL)
			//Clear all headers before sending GET request to AWS. Otherwise, it will return an error
			for k := range headers {
				delete(headers, k)
			}
		}

	}

	if *numOfConnection > 1 {
		//if we get the filename and filesize from Extractions/ReportExtractions, it will use the concurrent download
		step++
		log.Printf("Step %d: Concurrent Download: %s, Size: %d, Connection: %d\n", step, outputFilename, fileSize, *numOfConnection)
		trthrest.ConcurrentDownload(client, headers, downloadURL, outputFilename, *numOfConnection, fileSize, *traceFlag)
	} else {
		//if we can't get the filename and filesize from Extractions/ReportExtractions, it will download with one connection
		step++
		log.Printf("Step %d: Download: %s\n", step, outputFilename)
		trthrest.DownloadFile(client, headers, downloadURL, outputFilename, -1, -1, *traceFlag)
	}
	elapsed := time.Since(start)
	log.Printf("Download Time: %s\n", elapsed)
}
