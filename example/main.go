package main

import (
	"fmt"
	"net/http"

	"github.com/bykovme/gotrans"
)

func handler(w http.ResponseWriter, r *http.Request) {
	lang := gotrans.DetectLanguage(r.Header.Get("Accept-Language"))
	fmt.Fprintf(w, "<html><head><title> %s </title></head><body>", gotrans.Tr(lang, "hello_world"))
	fmt.Fprintf(w, "<h2> %s </h2>", gotrans.Tr(lang, "hello_world"))
	githubLink := "https://github.com/bykovme/gotrans"
	link := fmt.Sprintf(`<a href="%s">%s</a>`, githubLink, githubLink)
	fmt.Fprintf(w, gotrans.Tr(lang, "find_more"), link)
	fmt.Fprint(w, "</body></html>")
}

func main() {
	err := gotrans.InitLocales("langs")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}
