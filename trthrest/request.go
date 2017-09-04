package trthrest

//TickHistoryMarketDepthExtractionRequest : defined type for TickHistoryMarketDepthExtractionRequest. This type is used to request TRTH Market Depth data
//This type will be encoded to Json by Marshaller
type TickHistoryMarketDepthExtractionRequest struct {
	//It uses 'json' metadata to change the fieldname from Metadata to @data.type
	//It uses user-defined 'odata' metadata to define the default value
	Metadata          string                          `json:"@odata.type" odata:"#ThomsonReuters.Dss.Api.Extractions.ExtractionRequests.TickHistoryMarketDepthExtractionRequest"`
	ContentFieldNames []string                        `json:",omitempty"`
	IdentifierList    InstrumentIdentifierList        `json:",omitempty"`
	Condition         TickHistoryMarketDepthCondition `json:",omitempty"`
}
