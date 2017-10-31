package common

import (
	"fmt"
	"io/ioutil"
	"os"
)

//LoadHTMLFileAsString Load specified HTML file as one string
func LoadHTMLFileAsString(aFileName string) string {
	resourcePath := "C:/ws/ilidb/html/"
	tFileByteArray, err := ioutil.ReadFile(resourcePath + aFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	return string(tFileByteArray)
}
