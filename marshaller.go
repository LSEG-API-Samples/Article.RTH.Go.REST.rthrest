package trthrest

import (
	"encoding/json"
	"reflect"
)

//MarshalText : JSON Marshaller for TickHistoryExtractByMode enumeration.
//It uses tickHistoryExtractByMode string array variable to convert int (enum) to string
func (d TickHistoryExtractByMode) MarshalText() ([]byte, error) {
	return []byte(tickHistoryExtractByMode[d]), nil
}

//MarshalText : JSON Marshaller for PreviewMode enumeration.
//It uses previewMode string array variable to convert int (enum) to string
func (d PreviewMode) MarshalText() ([]byte, error) {
	return []byte(previewMode[d]), nil
}

//MarshalText : JSON Marshaller for ReportDateRangeType enumeration.
//It uses reportDateRangeType string array variable to convert int (enum) to string
func (d ReportDateRangeType) MarshalText() ([]byte, error) {
	return []byte(reportDateRangeType[d]), nil
}

//MarshalText : JSON Marshaller for TickHistoryMarketDepthViewOptions enumeration.
//It uses tickHistoryMarketDepthViewOptions string array variable to convert int (enum) to string
func (d TickHistoryMarketDepthViewOptions) MarshalText() ([]byte, error) {
	return []byte(tickHistoryMarketDepthViewOptions[d]), nil
}

//MarshalText : JSON Marshaller for TickHistorySort enumeration.
//It uses tickHistorySort string array variable to convert int (enum) to string
func (d TickHistorySort) MarshalText() ([]byte, error) {
	return []byte(tickHistorySort[d]), nil
}

//MarshalText : JSON Marshaller for TickHistoryTimeOptions enumeration.
//It uses tickHistoryTimeOptions string array variable to convert int (enum) to string
func (d TickHistoryTimeOptions) MarshalText() ([]byte, error) {
	return []byte(tickHistoryTimeOptions[d]), nil
}

//MarshalJSON : The custom JSON Marshaller for InstrumentIdentifierList. It uses reflection to set the value for 'Metadata' field.
//The default value is from 'odata" metadata
func (r InstrumentIdentifierList) MarshalJSON() ([]byte, error) {
	//This type is defined to avoid recursive while marshaling modified InstrumentIdentifierList
	//The modified InstrumentIdentifierList will be copied to this type
	//Therefore, json.Marshal can encode it to JSON with the value in 'Metatdata' field
	type _InstrumentIdentifierList InstrumentIdentifierList
	if r.Metadata == "" {
		st := reflect.TypeOf(r)
		field, _ := st.FieldByName("Metadata")
		r.Metadata = field.Tag.Get("odata")
	}
	return json.Marshal(_InstrumentIdentifierList(r))

}

//MarshalJSON : The custom JSON Marshaller for TickHistoryMarketDepthExtractionRequest. It uses reflection to set the value for 'Metadata' field.
//The default value is from 'odata" metadata
func (r TickHistoryMarketDepthExtractionRequest) MarshalJSON() ([]byte, error) {
	//This type is defined to avoid recursive while marshaling modified _TickHistoryMarketDepthExtractionRequest
	//The modified _TickHistoryMarketDepthExtractionRequest will be copied to this type.
	//Therefore, json.Marshal can encode it to JSON with the value in 'Metatdata' field
	type _TickHistoryMarketDepthExtractionRequest TickHistoryMarketDepthExtractionRequest
	if r.Metadata == "" {
		st := reflect.TypeOf(r)
		field, _ := st.FieldByName("Metadata")
		r.Metadata = field.Tag.Get("odata")
	}
	return json.Marshal(_TickHistoryMarketDepthExtractionRequest(r))
}
