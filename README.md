# Using Tick History V2 REST API with Go Programming Language

## Introduction

Thomson Reuters Tick History (TRTH) is an Internet-hosted product on the DataScope Select platform that provides SOAP-based and a REST API for unparalleled access to historical high frequency data across global asset classes dating to 1996. However a legacy SOAP-based API is also available, and is scheduled to be sunset. Therefore client who still uses SOAP-based API may need to upgrade their application to use REST API instead.

This article demonstrates problems and solutions that developers should aware when using TRTH V2 On Demand data extraction with Go programming language. It uses Tick History Market Depth On Demand data extraction as an example to demonstrate the usage and solutions. However, the methods mentioned in this article can be applied to other types of data extractions and programming languages.

## Prerequisite

The following knowledge is required to follow this article.

* You must know how to use On Demand extraction in TRTH V2. This article doesn't explain TRTH V2 REST API On Demand data extraction request in detail. Fortunately, there is a [REST API Tutorial 3: On Demand Data extraction workflow](https://developers.thomsonreuters.com/thomson-reuters-tick-history-trth/thomson-reuters-tick-history-trth-rest-api/learning?content=11307&type=learning_material_item) tutorial available in the Developer Community which thoroughly explains On Demand data extraction

* You must have basic knowledge of Go programming language. This article doesn't cover the installation, settings, and usage of Go programming language. You can refer to the official [Go Programming Language Website](https://golang.org/) for more information

## Overview

Go is an open source project under a BSD-style license developed by a team at Google in 2007 and many contributors from the open source community. Its binary distributions are available for Linux, Max OS X, Windows, and more. Go is a statically typed and compiled language with a simple syntax. It features garbage collection, concurrency, type safety and large standard library.

Developers can use Go programming language to consume Tick History data via TRTH V2 REST API. This article lists several problems and solutions which developers may find during development. The problems mentioned in this article include:

* Encode and decode JSON object
* Encode enumeration
* Concurrently download a gzip file
* Download a gzip file from Amazon Web Services

## Encode and Decode JSON Object
TRTH V2 REST API requires JSON  (JavaScript Object Notation) in request and response messages. JSON is a lightweight data-interchange format. It is easy for humans to read and write and machines to parse and generate. In Go programming language, there are several ways to encode and decode JSON objects. 

### Using a String to Encode and Decode JSON Object
JSON is a text format so the application can directly construct a JSON string for the HTTP request and process a JSON string in the HTTP response by using a string parser or regular expression. However, this method is ineffective and prone to error.

### Using a Map to Encode and Decode JSON Object
JSON is also a key and value pair data so **map[string]interface{}** can be used with **json.Marshal** and **json.Unmarshal** functions in the **encoding/json** package to encode and decode JSON objects. 

```
jsonMap := map[string]interface{}{
    "field1": "value1",
    "field2": 2,
    "a":      "1",
    "b":      2,
}
jsonByte, _ := json.Marshal(jsonMap)
fmt.Println(string(jsonByte))
```
The above code uses **map[string]interface{}** to store key/value pair data. Then, it uses **json.Marshal** function to encode the map to JSON byte array. After that, it prints the encoded JSON string. 

```
{"a":"1","b":2,"field1":"value1","field2":2}
```
To decode a JSON string, **json.Unmarshal** function can be used. 

```
var jsonMap map[string]interface{}
jsonStr := `{"field1":"value1","field2":2,"a":"1","b":2}`
json.Unmarshal([]byte(jsonStr), &jsonMap)
for k, v := range jsonMap {
    fmt.Printf("%s: %v\n", k, v)
}
```
The above code defines a JSON string in a string variable. Then, it calls the **json.Unmarshal** function to decode the JSON string to a map. After that, it prints keys and values in the map.
```
b: 2
field1: value1
field2: 2
a: 1
```

The drawback from this method is that the order of fields when encoding and decoding may not be out of order, as shown in the previous examples. This could be the problem when using with the API that the order of fields in the HTTP request must be preserved. 

### Using a Type to Encode and Decode JSON Object

In addition to a map, **json.Marshal** and **json.Unmarshal** functions can also be used with user-defined types. Therefore, JSON objects in HTTP body can be defined as types in Go programming language. Then, the types can be used with those functions to encode and decode JSON objects. This method is recommended and used by the example in this article.

In the example, the types for JSON objects in the request and response messages are defined as:
```
type TickHistoryMarketDepthExtractionRequest struct {	
	Metadata          string                          `json:"@odata.type" odata:"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest"`
	ContentFieldNames []string                        `json:",omitempty"`
	IdentifierList    InstrumentIdentifierList        `json:",omitempty"`
	Condition         TickHistoryMarketDepthCondition `json:",omitempty"`
}

type RawExtractionResult struct {
	Metadata string `json:"@odata.context,omitempty"`
	JobID                      string `json:"JobId"`
	Notes                      []string
	IdentifierValidationErrors []IdentifierValidationError
}
``` 
The first type is used in the HTTP request message to extract Tick History Market Depth data. The second type is for the JSON object sent in the HTTP response message when the extraction is completed. These types will be encoded and decoded as JSON objects. Each field represents a member of the JSON object by using a type's field name as a JSON object key.

JSON object in the HTTP request and response of TRTH V2 REST API contains **@data.type** field which defines a type name of OData. 
```
{
    "ExtractionRequest":{
        "@odata.type":"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest",
        "ContentFieldNames":[...]
        ...
    }
}
```
However, **@data.type** is an invalid field name in Go programming language. To solve this issue, the **json** tag is used in the **Metadata** field to customize the field name for JSON object.
```
Metadata string `json:"@odata.context,omitempty"`
```
The **omitempty** option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.

The value of **@odata.type** is unique and constant for each request type. For example, the value of **@odata.type** field for **TickHistoryMarketDepthExtractionRequest** is **#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest**. Therefore, it is inconvenient and prone to error, if this value will be set by users. Therefore, a custom field tag (**odata**) is defined for the this **Metadata** field so the user doesn't need to specify its value when using the TickHistoryMarketDepthExtractionRequest type.

```
type TickHistoryMarketDepthExtractionRequest struct {	
	Metadata          string                          `json:"@odata.type" odata:"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest"`
    ...
}
```
To use this **odata** tag, the custom JSON marshaller is defined for this type. The custom marshaller will be called when **json.Marshal** method is called on the defined type, such as **TickHistoryMarketDepthExtractionRequest**.  
```
func (r TickHistoryMarketDepthExtractionRequest) MarshalJSON() ([]byte, error) {
	type _TickHistoryMarketDepthExtractionRequest TickHistoryMarketDepthExtractionRequest
	if r.Metadata == "" {
		st := reflect.TypeOf(r)
		field, _ := st.FieldByName("Metadata")
		r.Metadata = field.Tag.Get("odata")
	}
	return json.Marshal(_TickHistoryMarketDepthExtractionRequest(r))
}
```
This marshaller uses reflection to get the value from **odata** tag, and set the value back to **Metadata** field. It also defines a new type (**_TickHistoryMarketDepthExtractionRequest**) as an alias for the marshaled type. After setting the value in the **Metadata** field, it casts the marshaled type to the new type and calls the **json.Marshal** function on this new type. Thus, the default marshaller of this new type will be used to marshal the data. This prevents the recursive call on this same custom marshaller. 

The following code shows how to use this user-defined type and marshaller to encode JSON object.

```
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

req1, _ := json.Marshal(struct {
    ExtractionRequest *trthrest.TickHistoryMarketDepthExtractionRequest
}{
    ExtractionRequest: request,
})
```
The above code is from the example which shows how to use **TickHistoryMarketDepthExtractionRequest** type and its marshaller to encode JSON object. The returned JSON object looks like:
```
{
    "ExtractionRequest":{
        "@odata.type":"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest",
        "ContentFieldNames":[
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
                "Sample Data"
        ],
        "IdentifierList":{
                "@odata.type":"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.InstrumentIdentifierList",
                "InstrumentIdentifiers":[
                        {
                            "Identifier":"IBM.N",
                            "IdentifierType":"Ric"
                        }
                ],
                "ValidationOptions":{
                        "AllowHistoricalInstruments":true
                }
        },
        "Condition":{
                "View":"NormalizedLL2",
                "NumberOfLevels":10,
                "SortBy":"SingleByRic",
                "MessageTimeStampIn":"GmtUtc",
                "ReportDateRangeType":"Range",
                "QueryStartDate":"2017-07-01T00:00:00Z",
                "QueryEndDate":"2017-08-23T00:00:00Z",
                "Preview":"None",
                "ExtractBy":"Ric",
                "DisplaySourceRIC":true
        }
    }
}
```
The JSON object shows that the **Metdata** field in **TickHistoryMarketDepthExtractionRequest** type is encoded as a **@odata.type** field with the value specified in the **odata** tag.

```
"@odata.type":"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest",
```

After the HTTP request is sent, the HTTP response will return with the JSON object. To decode the returned JSON object, **json.Unmarshal** function is called with **RawExtractionResult** type.

```
extractRawResult := &trthrest.RawExtractionResult{}
err = json.Unmarshal(body, extractRawResult)
```
The above code decodes the following JSON object to **RawExtractionResults** type.

```
{  
   "@odata.context":"https://hosted.datascopeapi.reuters.com/RestApi/v1/$metadata#RawExtractionResults/$entity",
   "JobId":"0x05db4e3626eb2f86",
   "Notes":[  
      "Extraction Services Version 11.1.37239 (5fcaa4f4395d), Built Aug  9 2017 15:35:02
      User ID: 9008895
      Extraction ID: 2000000002034110
      Schedule: 0x05db4e3626eb2f86 (ID = 0x0000000000000000)
      Input List (1 items):  (ID = 0x05db4e3626eb2f86) Created: 09/02/2017 11:21:10 Last Modified: 09/02/2017 11:21:10
      Report Template (12 fields): _OnD_0x05db4e3626eb2f86 (ID = 0x05db4e363d5b2f86) Created: 09/02/2017 11:19:20 Last Modified: 09/02/2017 11:19:20
      Schedule dispatched via message queue (0x05db4e3626eb2f86), Data source identifier (59BBF01D63F444CEB2E64CEE05F2ED4C)
      Schedule Time: 09/02/2017 11:19:21
      Processing started at 09/02/2017 11:19:23
      Processing completed successfully at 09/02/2017 11:21:12
      Extraction finished at 09/02/2017 04:21:12 UTC, with servers: tm04n01, TRTH (94.023 secs)
      Instrument <RIC,IBM.N> expanded to 1 RIC: IBM.N.
      Quota Message: INFO: Tick History Cash Quota Count Before Extraction: 1956; Instruments Extracted: 1; Tick History Cash Quota Count After Extraction: 1956, 391.2% of Limit; Tick History Cash Quota Limit: 500
      Manifest: #RIC,Domain,Start,End,Status,Count
      Manifest: IBM.N,Market Price,2017-07-03T11:30:01.198715182Z,2017-08-22T20:45:07.095340511Z,Active,489036"
   ]
}
```
The value of **@odata.type** field in JSON object is decoded to **Metadata** field in **RawExtractionResults** type.

In conclusion, using a type to encode and decode JSON object is effective and flexible. It is also useful when using with IDE that supports Intellisense, such as Visual Studio Code. Moreover, the user-defined types can be reused by other examples. 

## Encode enumeration
TRTH V2 REST API defines enumerations used in JSON object, such as **TickHistoryExtractByMode**, **TickHistoryMarketDepthViewOptions**, and **ReportDateRangeType**. Enumerations can also be defined in Go programming language and they can be used when constructing the request message.

```
type TickHistoryMarketDepthViewOptions int

const (
	ViewOptionsRawMarketByPriceEnum TickHistoryMarketDepthViewOptions = iota
	ViewOptionsRawMarketByOrderEnum
	ViewOptionsRawMarketMakerEnum
	ViewOptionsLegacyLevel2Enum
	ViewOptionsNormalizedLL2Enum
)
```
The above code defines an enumeration type called **TickHistoryMarketDepthViewOptions** and all enumeration values of this type are defined as constants.

The following shows how to use this enumeration.

```
request.Condition.View = trthrest.ViewOptionsNormalizedLL2Enum
```
Condition.View is **TickHistoryMarketDepthViewOptions** type and its value is set to **ViewOptionsNormalizedLL2Enum**. 

However, in JSON object, these enumeration fields are encoded as strings, not integers. To encode each enumeration as as string, an array of string and custom text marshaller are defined.

```
var tickHistoryMarketDepthViewOptions = [...]string{
	"RawMarketByPrice",
	"RawMarketByOrder",
	"RawMarketMaker",
	"LegacyLevel2",
	"NormalizedLL2",
}

func (d TickHistoryMarketDepthViewOptions) MarshalText() ([]byte, error) {
	return []byte(tickHistoryMarketDepthViewOptions[d]), nil
}
```
The above code defines an array of strings called **tickHistoryMarketDepthViewOptions** which contains a string for each enumeration value.  This array is used by the custom text marshaller of **TickHistoryMarketDepthViewOptions** type to convert an integer to a string while encoding to JSON object. For example, if the application sets the value of **TickHistoryMarketDepthViewOptions** type to **ViewOptionsNormalizedLL2Enum (4)**, when encoding JSON object, the custom text marshaller of this type will return a **"NormalizedLL2"** string which is the string at the fourth index in the array and this string will be used by the JSON marshaller for encoding, as shown below.

```
"Condition":{
                "View":"NormalizedLL2",
...
```

## Concurrently download a gzip file

The result file of **ExtractRaw** extraction is in **.csv.gz** format and the HTTP response when downloading the result file typically contains **Content-Encoding: gzip** in the header. With this header, the **net/http** library in Go programming language typically decompresses the gzip file and then returns the csv to the application. To download the raw gzip file, the decommpression must be disabled by using the following code.
```
tr := &http.Transport{
    DisableCompression: true,    
}
```
Depending on the number of instruments or the range of periods specified in the extraction request, the size of gzip file could be huge  According to TRTH V2 REST API User Guide, download speed is limited to 1 MB/s for each connection. Therefore, downloading the huge gzip file can take more than several hours with a single connection. 

To speed up the download, the file can download concurrently with multiple connections. Each connection will download a specific range of a file by defining a range (offset) in the HTTP request header. 

```
Range: bytes=0-3079590
```
The above header indicates that the HTTP request would like to download the first 3079591 bytes of the file. Fortunately, DSS server supports **Range** header so the status code of the HTTP response from DSS will be **206 Partial Content**. 

```
HTTP/1.1 206 Partial Content
Content-Length: 3079591
Cache-Control: no-cache
Content-Range: bytes 0-3079590/12318367
Content-Type: text/plain
Date: Sun, 03 Sep 2017 07:34:05 GMT
```
The response also indicates the content size, starting, and ending offset. 

However, in order to download file concurrently, the staring and ending offset of each download connection must be calculated so the size of result file must be known. There are several ways to get the size of the extracted file. The example in this article shows the way to get the size of result file by using the **Extraction ID** appearing in the **Notes** field when the job is completed. 

```
{
  "@odata.context": "https://hosted.datascopeapi.reuters.com/RestApi/v1/$metadata#RawExtractionResults/$entity",
  "JobId": "0x05dbaba5eceb2f76",
  "Notes": [
    "Extraction Services Version 11.1.37239 (5fcaa4f4395d), Built Aug 21 2017 20:06:16
    User ID: 9008895
    Extraction ID: 2000000002049332
    Schedule: 0x05dbaba5eceb2f76 (ID = 0x0000000000000000)
    Input List (1 items):  (ID = 0x05dbaba5eceb2f76) Created: 09/03/2017 14:33:56 Last Modified: 09/03/2017 14:33:56
    Report Template (12 fields): _OnD_0x05dbaba5eceb2f76 (ID = 0x05dbaba63edb2f76) Created: 09/03/2017 14:32:16 Last Modified: 09/03/2017 14:32:16
    Schedule dispatched via message queue (0x05dbaba5eceb2f76), Data source identifier (19624AF632374B9B8613138BEDA99FC6)
    Schedule Time: 09/03/2017 14:32:17
    Processing started at 09/03/2017 14:32:17
    Processing completed successfully at 09/03/2017 14:33:56
    Extraction finished at 09/03/2017 07:33:56 UTC, with servers: tm01n01
    Instrument <RIC,IBM.N> expanded to 1 RIC: IBM.N.
    Quota Message: INFO: Tick History Cash Quota Count Before Extraction: 1956; Instruments Extracted: 1; Tick History Cash Quota Count After Extraction: 1956, 391.2% of Limit; Tick History Cash Quota Limit: 500
    Manifest: #RIC,Domain,Start,End,Status,Count
    Manifest: IBM.N,Market Price,2017-07-03T11:30:01.198715182Z,2017-08-22T20:45:07.095340511Z,Active,489036"
  ]
}
```
From the above response, the **Extraction ID** in the **Notes** field is 2000000002049332. To get the file description, the following HTTP GET request is used.

```
GET /RestApi/v1/Extractions/ReportExtractions('2000000002049332')/FullFile
```
The response for this request is the description of the data file.
```
{
  "@odata.context": "https://hosted.datascopeapi.reuters.com/RestApi/v1/$metadata#ExtractedFiles/$entity",
  "ExtractedFileId": "VjF8MHgwNWRiYWJiZWI2YWIzMDE2fA",
  "ReportExtractionId": "2000000002049332",
  "ScheduleId": "0x05dbaba5eceb2f76",
  "FileType": "Full",
  "ExtractedFileName": "_OnD_0x05dbaba5eceb2f76.csv.gz",
  "LastWriteTimeUtc": "2017-09-03T07:33:56.000Z",
  "ContentsExists": true,
  "Size": 12318367
}
```
The **Size** field in the response contains the size of file. Then, the download byte offset can be calculated for each connection by dividing the size of file by the number of connections. For example, if the above file is downloaded concurrently with four connections, the download size for each connection will be 3079591 bytes (12318367 / 4) and the download offets for four connections will be:

```
Connection 1: Range: Bytes=0-3079590
Connection 2: Range: Bytes=3079591-6159181
Connection 3: Range: Bytes=6159182-9238772
Connection 4: Range: Bytes=9238773 -
```
The fourth connection will start downloading the file starting at 9238773 offset until the end of file. After all connections complete downloading files, all files must be merged in offset order to get the completed file.

The following test results compare the download times between a single connection and four connections.

|No.|Total download time (seconds) with a single connection|Total download time (seconds) with four concurrent connections|
| ------------- |-------------|-----|
|1|43.832|24.675|
|2|111.683|26.654|
|3|63.658|29.655|
|4|46.807|33.009|
|5|89.659|25.013|
|6|66.757|20.037|
|7|54.846|25.15|
|8|106.874|18.664|
|9|56.865|19.841|
|10|55.628|45.135|

After testing ten times, downloading a file with four concurrent connections is faster than download a file with a single connection. The test results may vary according to machine and netowrk performance.

## Download a gzip file from Amazon Web Services

In addition to download extracted files directly from DSS server, the application can download the files faster by retrieving them directly from the Amazon Web Services (AWS) cloud in which they are hosted. This feature is available for VBD (Venue by Day) data, Tick History Time and Sales, Tick History Market Depth, Tick History Intraday Summaries, and Tick History Raw reports.

To use this feature, the application must include the HTTP header **X-Direct-Download: true** in the request. If the file is available on AWS, the status code of HTTP response will be **302 Found** with the new AWS URL in the **Location** HTTP header field. The new URL is the pre-signed URL to get data directly from AWS. 

```
HTTP/1.1 302 Found
Cache-Control: no-cache
Date: Sun, 03 Sep 2017 10:34:17 GMT
Expires: -1
Location: https://s3.amazonaws.com/tickhistory.query.production.hdc-results/xxx/data/merged/merged.csv.gz?AWSAccessKeyId=xxx&Expires=1504456458&response-content-disposition=attachment%3B%20filename%3D_OnD_0x05dbb5f5a62b3016.csv.gz&Signature=xxx&x-amz-request-payer=requester
```

Then, the application can use this new AWS URL to download the file.

However, when retrieving the HTTP status code 302, the **http** library in Go programming language will automatically redirect to the new URL with the same HTTP headers which have fields for TRTH V2 REST API. This causes AWS returning **403 Forbidden** status code.

To avoid this issue, the application should disable this automatic redirect by using the following code.

```
client := &http.Client{
    Transport: &tr,
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
}
```
Then, the application can remove TRTH V2 headers and optionally add its own HTTP headers in the request. Concurrent downloads mentioned in the previous section can also be used with AWS by specifying **Range** header in the request.

## Go Get and Run the Example
**TickHistoryMarketDepthEx.go** is implemented to demonstrate methods mentioned in this article. It uses **ExtractRaw** endpoint to send **TickHistoryMarketDepthExtractionRequest** to extract normallized legacy level 2 data of IBM.N from 1 Jul 2017 to 23 Aug 2017. All settings are hard-coded. This example supports the following features:
* Concurrent Downloads
* Download a file from AWS
* Request and response tracing
* Proxy setting

This example depends on the **github.com/howeyc/gopass** package in order to retrieve the DSS password from the console.

The optional arguments for this example are:

|Argument Name|Description|Argument Type (Default Value)|
|-------------|-----------|-------------|
|--help|List all valid arguments||
|-u|Specify the DSS user name|String ("")|
|-p|Specify the DSS password|String ("")|
|-n|Specify the number of concurrent downloads|Integer (1)|
|-aws|Flag to download from  AWS|Boolean (false)|
|-X|Flag to trace HTTP request and response|Boolean (false)|

To download the example, please run the following command.

```
go get github.com/TR-API-Samples/Article.TRTH.Go.REST.trthrest/main
```

The example can be run with the following command.

```
go run github.com/TR-API-Samples/Article.TRTH.Go.REST.trthrest/main/TickHistoryMarketDepthEx.go -aws -n 4
```

The above command runs the example to download the result file from AWS with four concurrent connections. The output is shown below.
```
2017/09/05 15:21:28 X-Direct-Download: true
2017/09/05 15:21:28 Number of concurrent download: 4
Enter DSS Username: 9008895
Enter DSS Password: **********
2017/09/05 15:21:43 Step 1: RequestToken
2017/09/05 15:21:46 Step 2: ExtractRaw for TickHistoryMarketDepthExtractionRequest
2017/09/05 15:22:19 Step 3: Checking Status (202) of Extraction (1)
2017/09/05 15:22:52 Step 3: Checking Status (202) of Extraction (2)
2017/09/05 15:23:25 Step 3: Checking Status (202) of Extraction (3)
2017/09/05 15:23:47 ExtractionID: "2000000002101180"
2017/09/05 15:23:47 Step 4: Get File information
2017/09/05 15:23:48 File: _OnD_0x05dc5346defb3036.csv.gz, Size: 12318367
2017/09/05 15:23:48 Step 5: Get AWS URL
2017/09/05 15:23:49 AWS: https://s3.amazonaws.com/tickhistory.query.production.hdc-results/E72B426B7E554FA48E15D96B60162A93/data/merged/merged.csv.gz?AWSAccessKeyId=xxx&Expires=1504621429&response-content-disposition=attachment%3B%20filename%3D_OnD_0x05dc5346defb3036.csv.gz&Signature=xxx&x-amz-request-payer=requester
2017/09/05 15:23:49 Step 6: Concurrent Download: _OnD_0x05dc5346defb3036.csv.gz, Size: 12318367, Connection: 4
2017/09/05 15:23:49 ConcurrentDownload: _OnD_0x05dc5346defb3036.csv.gz, conn=4
2017/09/05 15:23:49 Part 1: 0 - 3079590
2017/09/05 15:23:49 Part 2: 3079591 - 6159181
2017/09/05 15:23:49 Part 3: 6159182 - 9238772
2017/09/05 15:23:49 Download File: part1, 0, 3079590
2017/09/05 15:23:49 Download File: part2, 3079591, 6159181
2017/09/05 15:23:49 Part 4: 9238773-
2017/09/05 15:23:49 Download File: part4, 9238773, -1
2017/09/05 15:23:49 Download File: part3, 6159182, 9238772
2017/09/05 15:23:54 part2, Bytes: 956719/Total: 3079591 (31%)
2017/09/05 15:23:54 part4, Bytes: 1183022/Total: 3079594 (38%)
2017/09/05 15:23:54 part3, Bytes: 1095983/Total: 3079591 (36%)
2017/09/05 15:23:54 part1, Bytes: 1600821/Total: 3079591 (52%)
2017/09/05 15:23:55 part3: Download Completed, Speed: Avg 601.40 KB/s, Max 663.00 KB/s
2017/09/05 15:23:55 part1: Download Completed, Speed: Avg 601.40 KB/s, Max 816.00 KB/s
2017/09/05 15:23:56 part2: Download Completed, Speed: Avg 501.17 KB/s, Max 1377.00 KB/s
2017/09/05 15:23:56 part4: Download Completed, Speed: Avg 501.17 KB/s, Max 1105.00 KB/s
2017/09/05 15:23:56 Merging Files: _OnD_0x05dc5346defb3036.csv.gz
2017/09/05 15:23:56 Download Time: 8.231s
```

## References

* [Go Programming Language](https://golang.org/)
* [Thomson Reuters Tick History (TRTH) - REST API](https://developers.thomsonreuters.com/thomson-reuters-tick-history-trth/thomson-reuters-tick-history-trth-rest-api)
* [JavaScript Object Notation](www.json.org/)
* [Go: Package json](https://golang.org/pkg/encoding/json/)
* [Go: Package http](https://golang.org/pkg/net/http/) 