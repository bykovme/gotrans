[![Build Status](https://travis-ci.org/bykovme/gotrans.svg?branch=master)](https://travis-ci.org/bykovme/gotrans)
[![codecov](https://codecov.io/gh/bykovme/gotrans/branch/master/graph/badge.svg)](https://codecov.io/gh/bykovme/gotrans)
[![Go Report Card](https://goreportcard.com/badge/github.com/bykovme/gotrans)](https://goreportcard.com/report/github.com/bykovme/gotrans)

# gotrans - localization package for golang

Localization package for GO, it uses JSON files to localize your GO application

## Getting started

### Installation

Install the package with the command
```bash
go get github.com/bykovme/gotrans
```

### Prepare translation files

JSON files should use following format:

```json
{
    "hello_world":"Hello World",
    "find_more":"Find more information about the project on the website %s"
}
```

JSON file name should use standard [language code](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) or language-country code supported by browsers. 
At least one file "en.json" should be in the localization folder.

The folder content should look like that:
```
    en.json
    de.json
    ru.json
```

### Quick documentation  

There are just 3 functions in the package

#### InitLocales(path string)

Use the relative or absolute path to set the folder where all the JSON files with translations are located. Make sure that all the files with translations have extension ".json"

#### Tr(lang string, key string) string

Get translation value by the language & key 

#### DetectLanguage(acceptLanguage string) string 

This function will be useful when you are creating web application, it detects the language from HTTP header Accept-Language, check the usage of the function in the example below

## Example of using gotrans package

The same example is located within the package [here](https://github.com/bykovme/gotrans/tree/master/example)

```go
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
```

## Behaviour

If the key is not found in the localization file, it will try to find the same key in English localization ("en.json"), if the key is not found there as well, the key will be returned instead of value.

## The End


