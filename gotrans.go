package gotrans

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Translation - localization structure
type translation struct {
	locales       []string
	translations  map[string]map[string]string
	defaultLocale string
}

var trans *translation

// InitLocales - initiate locales from the folder.
//
// Parameters:
//
// 'trPath' - path to the folder with translations files
//
// Use the relative or absolute path to set the folder where all the JSON files
// with translations are located.
// Make sure that all the files with translations have extension ".json".
//
// Examples:
//
// err := gotrans.InitLocales("/home/user/project/languages") // absolute path
//
// err := gotrans.InitLocales("languages") // relative path
func InitLocales(trPath string) error {
	trans = &translation{translations: make(map[string]map[string]string)}
	return loadTranslations(trPath)
}

// Tr - find translation for provided locale and translation key
//
// Parameters
//
// 'locale' - locale value, for example: "en", "jp", "de", "ru"
//
// 'trKey' - translation key from json file
//
// IMPORTANT! Call gotrans.InitLocale() to initiate translations before calling this function
func Tr(locale string, trKey string) string {
	if trans == nil {
		return ""
	}
	trValue, ok := trans.translations[locale][trKey]
	if ok {
		return trValue
	}
	trValue, ok = trans.translations["en"][trKey]
	if ok {
		return trValue
	}
	return trKey
}

// SetDefaultLocale - set new default locale
func SetDefaultLocale(newLocale string) error {
	if trans == nil {
		return errors.New("translations are not initialized")
	}
	if checkLocale(newLocale) {
		trans.defaultLocale = newLocale
		return nil
	}
	return errors.New("locale is not found")
}

// GetDefaultLocale - return current default locale
func GetDefaultLocale() string {
	if trans == nil {
		return ""
	}
	return trans.defaultLocale
}

// T - find translation for default locale and provided translation key
//
// Parameters
//
// 'trKey' - translation key from json file
//
// IMPORTANT! Call gotrans.InitLocale() to initiate translations before calling this function
// and gotrans.SetDefaultLocale() to set up proper translation lacale (if it is not set then
// the library will use locale "en")
func T(trKey string) string {
	if trans == nil {
		return ""
	}
	return Tr(trans.defaultLocale, trKey)
}

// DetectLanguage - parse to find the most preferable language
func DetectLanguage(acceptLanguage string) string {

	langList := strings.Split(acceptLanguage, ",")
	for _, langStr := range langList {
		lang := strings.Split(strings.Trim(langStr, " "), ";")
		if checkLocale(lang[0]) {
			return lang[0]
		}
	}

	return "en"
}

// loadTranslations - load translations files from the folder
func loadTranslations(trPath string) error {
	files, _ := filepath.Glob(trPath + "/*.json")

	if len(files) == 0 {
		return errors.New("no translations found")
	}

	for _, file := range files {
		err := loadFileToMap(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadFileToMap(filename string) error {
	var trMap map[string]string

	localName := strings.Replace(filepath.Base(filename), ".json", "", 1)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &trMap)
	if err != nil {
		return err
	}
	trans.translations[localName] = trMap
	trans.locales = append(trans.locales, localName)
	if trans.defaultLocale == "" || localName == "en" { // Set 'en' by default if it is found
		trans.defaultLocale = localName
	}
	return nil
}

func checkLocale(localeName string) bool {
	for _, locale := range trans.locales {
		if locale == localeName {
			return true
		}
	}
	return false
}

// GetLocales - get available locales
func GetLocales() []string {
	if trans == nil {
		return nil
	}
	return trans.locales
}
