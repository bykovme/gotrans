# gotrans
Localization library for GO, it uses *.json files to localize your GO application

## Getting started

JSON files should use following format:

```json
{
    "hello_world":"Hello World",
    "find_more":"Find more information about the project on the website %s"
}
```

JSON file name should represent standard [language code](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) or language-country code supported by browsers. 

There are just 3 functions in the package

### InitLocales(path string)

Use the relative path in the project to set the folder within the project where all the JSON files are located. Make sure that files have extension ".json"

### Tr(lang string, key string) string

Get translation value by the language & key 

### DetectLanguage(acceptLanguage string) string 

This function will be useful when you are creating web application, it detects the language from HTTP header Accept-Language, check the usage of the function in the example below

## Example of using gotrans package

```go
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
```

Feel free to request new features or send your pull requests.
