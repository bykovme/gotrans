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
	locales      []string
	translations map[string]map[string]string
}

var trans *translation

// Tr - translate for current locale
func Tr(locale string, trKey string) string {
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

// InitLocales - initiate locales from the folder
func InitLocales(trPath string) error {
	trans = &translation{translations: make(map[string]map[string]string)}
	return loadTranslations(trPath)
}

// LoadTranslations - load translations files from the folder
func loadTranslations(trPath string) error {
	files, err := filepath.Glob(trPath + "/*.json")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("No translations found")
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
	var objmap map[string]string

	localName := strings.Replace(filepath.Base(filename), ".json", "", 1)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &objmap)
	if err != nil {
		return err
	}
	trans.translations[localName] = objmap
	trans.locales = append(trans.locales, localName)
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

// DetectLanguage - parse to find the most preferable language
func DetectLanguage(acceptLanguage string) string {

	langStrs := strings.Split(acceptLanguage, ",")
	for _, langStr := range langStrs {
		lang := strings.Split(strings.Trim(langStr, " "), ";")
		if checkLocale(lang[0]) {
			return lang[0]
		}
	}

	return "en"
}
