package common

import (
	"testing"
)

func TestLoadExistingHTMLFileAsString(t *testing.T) {
	indexPage := LoadHTMLFileAsString("mainIndex.html")
	if len(indexPage) == 0 {
		t.Fail()
		println("Could not load index.html as string")
	}
}

func TestLoadMissingHTMLFileAsString(t *testing.T) {
	indexPage := LoadHTMLFileAsString("i222ndex.html")
	if len(indexPage) > 0 {
		t.Fail()
		println("Should not be able to load i222ndex.html as string")
	}
}
