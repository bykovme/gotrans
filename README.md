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
Initiate the library using the folder with localization files with the function gotrans.InitLocales(), see detailed documentation below

The folder content should look like that:
```
    en.json
    de.json
    ru.json
```

## Example of using gotrans package

Simple usage example (retrieving translation with a locale and translation key)
```go
package main

import (
	"fmt"

	"github.com/bykovme/gotrans"
)

func main() {
    _ := gotrans.InitLocales("/home/user/project/languages")  //  Path to the folder with localization files
    fmt.Println(gotrans.Tr("en", "hello_world"))  // Using english translation from the file 'en.json'
    fmt.Println(gotrans.Tr("ru", "hello_world"))  // Using russian translation from the file 'ru.json'
}
```

Using default language and shorter function call (retrieving translation without a locale)
```go
package main

import (
	"fmt"

	"github.com/bykovme/gotrans"
)

func main() {
	_ = gotrans.InitLocales("/home/user/project/languages")  //  Path to the folder with localization files     
	_ = gotrans.SetDefaultLocale("ru") // Setting default locale
    
	fmt.Println(gotrans.T("hello_world"))  // Using russian translation from the file 'ru.json'
}
```


More complicated usage example for dynamic usage on the webserver.
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

## Quick documentation

#### FUNCTIONS

func DetectLanguage(acceptLanguage string) string
DetectLanguage - parse to find the most preferable language

func GetDefaultLocale() string
GetDefaultLocale - return current default locale

func GetLocales() []string
GetLocales - get available locales

func InitLocales(trPath string) error
InitLocales - initiate locales from the folder.

    Parameters:

    'trPath' - path to the folder with translations files

    Use the relative or absolute path to set the folder where all the JSON files
    with translations are located. Make sure that all the files with
    translations have extension ".json".

    Examples:

    err := gotrans.InitLocales("/home/user/project/languages") // absolute path

    err := gotrans.InitLocales("languages") // relative path

func SetDefaultLocale(newLocale string) error
SetDefaultLocale - set new default locale

func T(trKey string) string
T - find translation for default locale and provided translation key


    Parameters

    'trKey' - translation key from json file

    IMPORTANT! Call gotrans.InitLocale() to initiate translations before calling
    this function and gotrans.SetDefaultLocale() to set up proper translation
    lacale (if it is not set then the library will use locale "en")

func Tr(locale string, trKey string) string
Tr - find translation for provided locale and translation key


    Parameters

    'locale' - locale value, for example: "en", "jp", "de", "ru"

    'trKey' - translation key from json file

    IMPORTANT! Call gotrans.InitLocale() to initiate translations before calling
    this function


**[Alex Bykov](https://bykovsoft.com) © 2015 - 2020**


