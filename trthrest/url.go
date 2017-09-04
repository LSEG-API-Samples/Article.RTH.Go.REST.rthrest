package trthrest

func GetRequestTokenURL(trthapiurl string)(string){
	return trthapiurl+"Authentication/RequestToken"
}
func GetExtractRawURL(trthapiurl string)(string){
	return trthapiurl+"Extractions/ExtractRaw"
}
func GetReportExtractionFullFileURL(trthapiurl string, extractionId string)(string){
	return trthapiurl + "Extractions/ReportExtractions('" + extractionId + "')/FullFile"
}
func GetRawExtractionResultGetDefaultStreamURL(trthapiurl string, jobId string)(string){
	return  trthapiurl + "Extractions/RawExtractionResults('" + jobId + "')" + "/$value"
}