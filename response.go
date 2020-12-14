package rthrest
import "time"
//RequestTokenResponse : The HTTP response from Authentication/RequestToken request will be decoded to this type by json.Unmarshal
type RequestTokenResponse struct {
	//The value in @odata.content field will be decoded to Metadata field
	Metadata string `json:"@odata.context,omitempty"`
	Value    string
}

//IdentifierValidationError : If the indentifier in the request is invalid, the repsonse from extraction (RawExtractionResult) will contain this information.
type IdentifierValidationError struct {
	//The value in '@odata.content' field will be decoded to 'Metadata' field
	Metadata   string `json:"@odata.context,omitempty"`
	Identifier InstrumentIdentifier
	Message    string
}

//RawExtractionResult : The HTTP response from the completed Extractions/ExtractRaw request will be decoded to this type by json.Unmarshal
type RawExtractionResult struct {
	//The value in '@odata.content' field will be decoded to this 'Metadata' field
	Metadata string `json:"@odata.context,omitempty"`
	//The value in 'JobId' field will be decoded to this 'JobID field
	JobID                      string `json:"JobId"`
	Notes                      []string
	IdentifierValidationErrors []IdentifierValidationError
}
//ExtractedFile : Return the data file information from the extractionID
type ExtractedFile struct {
	Metadata string `json:"@odata.context,omitempty"`
	ExtractedFileId string
	ReportExtractionId string
	ScheduleId string
	FileType string
	ExtractedFileName string
	LastWriteTimeUtc *time.Time
	ContentsExists bool
	Size int64
	ReceivedDateUtc *time.Time
}