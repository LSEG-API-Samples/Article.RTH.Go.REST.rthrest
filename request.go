package rthrest

//TickHistoryMarketDepthExtractionRequest : defined type for TickHistoryMarketDepthExtractionRequest. This type is used to request RTH Market Depth data
//This type will be encoded to Json by Marshaller
type TickHistoryMarketDepthExtractionRequest struct {
	//It uses 'json' metadata to change the fieldname from Metadata to @data.type
	//It uses user-defined 'odata' metadata to define the default value
	Metadata          string                          `json:"@odata.type" odata:"#DataScope.Select.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest"`
	ContentFieldNames []string                        `json:",omitempty"`
	IdentifierList    InstrumentIdentifierList        `json:",omitempty"`
	Condition         TickHistoryMarketDepthCondition `json:",omitempty"`
}
