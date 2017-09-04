package trthrest

//TickHistoryMarketDepthViewOptions : This is an enumeration for TickHistoryMarketDepthViewOptions
type TickHistoryMarketDepthViewOptions int

//TickHistorySort : This is an enumeration for TickHistorySort
type TickHistorySort int

//TickHistoryTimeOptions :  This is an enumeration for TickHistoryTimeOptions
type TickHistoryTimeOptions int

//ReportDateRangeType : This is an enumeration for ReportDateRangeType
type ReportDateRangeType int

//PreviewMode : This is an enumeration for PreviewMode
type PreviewMode int

//TickHistoryExtractByMode : This is an enumeration for TickHistoryExtractByMode
type TickHistoryExtractByMode int

//Available Enumerations for TickHistoryMarketDepthViewOptions
const (
	ViewOptionsRawMarketByPriceEnum TickHistoryMarketDepthViewOptions = iota
	ViewOptionsRawMarketByOrderEnum
	ViewOptionsRawMarketMakerEnum
	ViewOptionsLegacyLevel2Enum
	ViewOptionsNormalizedLL2Enum
)

//Available Enumerations for TickHistoryExtractByMode
const (
	ExtractByModeRicEnum TickHistoryExtractByMode = iota
	ExtractByModeEntityEnum
)

//Available Enumerations for PreviewMode
const (
	PreviewModeNoneEnum PreviewMode = iota
	PreviewModeContentEnum
	PreviewModeInstrumentEnum
)

//Available Enumerations for TickHistorySort
const (
	SortSingleByRicEnum TickHistorySort = iota
	SortSingleByTimestampEnum
)

//Available Enumerations for TickHistoryTimeOptions
const (
	TimeOptionsLocalExchangeTimeEnum TickHistoryTimeOptions = iota
	TimeOptionsGmtUtcEnum
)

//Available Enumerations for ReportDateRangeType
const (
	ReportDateRangeTypeNoRangeEnum ReportDateRangeType = iota
	ReportDateRangeTypeInitEnum
	ReportDateRangeTypeRangeEnum
	ReportDateRangeTypeDeltaEnum
	ReportDateRangeTypeLastEnum
)

//Enumeration String of tickHistoryExtractByMode enumeration used by Marshaller while encoding to JSON
var tickHistoryExtractByMode = [...]string{"Ric", "Entity"}

//Enumeration String of previewMode enumeration used by Marshaller while encoding to JSON
var previewMode = [...]string{"None", "Content", "Instrument"}

//Enumeration String of reportDateRangeType enumeration used by Marshaller while encoding to JSON
var reportDateRangeType = [...]string{
	"NoRange",
	"Init",
	"Range",
	"Delta",
	"Last",
}

//Enumeration String of tickHistoryMarketDepthViewOptions enumeration used by Marshaller while encoding to JSON
var tickHistoryMarketDepthViewOptions = [...]string{
	"RawMarketByPrice",
	"RawMarketByOrder",
	"RawMarketMaker",
	"LegacyLevel2",
	"NormalizedLL2",
}

//Enumeration String of tickHistorySort enumeration used by Marshaller while encoding to JSON
var tickHistorySort = [...]string{
	"SingleByRic",
	"SingleByTimestamp",
}

//Enumeration String of tickHistoryTimeOptions enumeration used by Marshaller while encoding to JSON
var tickHistoryTimeOptions = [...]string{
	"LocalExchangeTime",
	"GmtUtc",
}
