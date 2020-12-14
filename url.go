package rthrest

func GetRequestTokenURL(rthapiurl string)(string){
	return rthapiurl+"Authentication/RequestToken"
}
func GetExtractRawURL(rthapiurl string)(string){
	return rthapiurl+"Extractions/ExtractRaw"
}
func GetReportExtractionFullFileURL(rthapiurl string, extractionId string)(string){
	return rthapiurl + "Extractions/ReportExtractions('" + extractionId + "')/FullFile"
}
func GetRawExtractionResultGetDefaultStreamURL(rthapiurl string, jobId string)(string){
	return  rthapiurl + "Extractions/RawExtractionResults('" + jobId + "')" + "/$value"
}