package controller

var recordTypes = []string{
	"A",
	"AAAA",
	"CNAME",
	"MX",
	"NS",
	"SOA",
	"SRV",
	"TXT",
}

func initializeZoneFile() map[string][]map[string]string {
	zoneFile := make(map[string][]map[string]string)
	for _, recordType := range recordTypes {
		zoneFile[recordType] = make([]map[string]string, 0)
	}
	return zoneFile
}
