package web

import (
	"bytes"
	"com/ilidb/common"
	"com/ilidb/db"
	"html/template"
)

//GenerateIndexPage Index page.
func GenerateIndexPage() string {
	topString := generatePageTop()
	tIndexMainBeginning := "<div id=\"main\">"

	tContentTrendingBooks := GenerateContentTrendingBooks()
	tContentNewBooks := GenerateContentNewBooks()

	tIndexMainEnd := "</div>"
	footString := generatePageFoot()
	return topString + tIndexMainBeginning + tContentTrendingBooks + tContentNewBooks + tIndexMainEnd + footString
}

//GenerateContentTrendingBooks Popular trending books.
func GenerateContentTrendingBooks() string {
	tBooksIter := db.FetchPopularBooks(20)
	tContentTrendingBooksString := executeTemplate("contentTrendingBooks.html", tBooksIter)
	return tContentTrendingBooksString
}

//GenerateContentNewBooks Popular new books.
func GenerateContentNewBooks() string {
	tBooksIter := db.FetchPopularBooks(20)
	tContentNewBooksString := executeTemplate("contentNewBooks.html", tBooksIter)
	return tContentNewBooksString
}

//GenerateLoginPage Login page.
func GenerateLoginPage() string {
	topString := generatePageTop()
	tIndexMainString := common.LoadHTMLFileAsString("mainLogin.html")
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + tIndexMainString + pageSideString + footString
}

//GenerateBooksPage Books page.
func GenerateBooksPage() string {
	tBooksIter := db.FetchPopularBooks(20)
	pageMainString := executeTemplate("mainBooks.html", tBooksIter)

	topString := generatePageTop()
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + pageMainString + pageSideString + footString
}

//GenerateSearchPage Search page.
func GenerateSearchPage(aQuery string) string {
	tSearchResult := db.SearchBookTitle(aQuery)
	pageMainString := executeTemplate("mainSearch.html", tSearchResult)

	topString := generatePageTop()
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + pageMainString + pageSideString + footString
}

//GenerateUserVotesPage User votes page.
func GenerateUserVotesPage(aUser db.User) string {
	tUserBookVotes := aUser.BookVotes
	pageMainString := executeTemplate("mainUserVotes.html", tUserBookVotes)

	topString := generatePageTop()
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + pageMainString + pageSideString + footString
}

//GenerateBookCategoryPage Page with books of a specific category.
func GenerateBookCategoryPage(aBookCategory string) string {
	tBooksIter := db.FetchPopularBooksCategory(aBookCategory, 20)
	// Uses same main page template as books page (books = all categories)
	pageMainString := executeTemplate("mainBooks.html", tBooksIter)

	topString := generatePageTop()
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + pageMainString + pageSideString + footString
}

//GenerateBookPage Book page.
func GenerateBookPage(aBookID string) string {
	var tBook db.Book
	tBook, _ = db.FetchBook(aBookID)

	pageMainString := executeTemplate("mainBook.html", tBook)

	topString := generatePageTop()
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + pageMainString + pageSideString + footString
}

//GenerateContributePage Contribute page.
func GenerateContributePage() string {
	topString := generatePageTop()
	pageMainString := common.LoadHTMLFileAsString("contribute.html")
	pageSideString := generatePageSideBookCategories()
	footString := generatePageFoot()
	return topString + pageMainString + pageSideString + footString
}

//generatePageTop Page top.
func generatePageTop() string {
	pageTopString := common.LoadHTMLFileAsString("top.html")
	return pageTopString
}

//generatePageSideBookCategories Side div container with book categories.
func generatePageSideBookCategories() string {
	pageSideString := common.LoadHTMLFileAsString("sideBookCategories.html")
	return pageSideString
}

//generatePageFoot Page foot content.
func generatePageFoot() string {
	pageFootString := common.LoadHTMLFileAsString("foot.html")
	return pageFootString
}

//executeTemplate Execute a template HTML file with given data
func executeTemplate(aHTMLFile string, aTemplateData interface{}) string {
	tTemplateFileString := common.LoadHTMLFileAsString(aHTMLFile)
	tmpl, err := template.New("").Parse(tTemplateFileString)
	if err != nil {
		panic(err)
	}

	var tBuffer bytes.Buffer
	err = tmpl.Execute(&tBuffer, aTemplateData)
	if err != nil {
		panic(err)
	}
	resultString := tBuffer.String()

	return resultString
}
