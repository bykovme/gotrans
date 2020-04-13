package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bykovme/gotrans"
)

const cHtmlTemplate = `
<html>
	<head>
		<title>%s</title>
	</head>
	<body>
		<h2> %s </h2>
		%s 
	</body>
</html>
`
const cGitHubLink = "https://github.com/bykovme/gotrans"
const cLink = `<a href="%s">%s</a>`

func handler(w http.ResponseWriter, r *http.Request) {
	lang := gotrans.DetectLanguage(r.Header.Get("Accept-Language"))
	helloWorld := gotrans.Tr(lang, "hello_world")
	link := fmt.Sprintf(cLink, cGitHubLink, cGitHubLink)
	findMore := fmt.Sprintf(gotrans.Tr(lang, "find_more"), link)
	_, err := fmt.Fprintf(w, cHtmlTemplate,
		helloWorld,
		helloWorld,
		findMore)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	err := gotrans.InitLocales("langs")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(":3000", nil)
}
