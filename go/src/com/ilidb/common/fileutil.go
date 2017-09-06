package common

import (
	"io/ioutil"
)

//LoadHTMLFileAsString Load specified HTML file as one string
func LoadHTMLFileAsString(aFileName string) string {
	resourcePath := "C:/ws/ilidb/html/"
	tFileByteArray, err := ioutil.ReadFile(resourcePath + aFileName)
	if err != nil {
		panic(err)
	}
	return string(tFileByteArray)
}
