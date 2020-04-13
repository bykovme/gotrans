package gotrans

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

const jsonLanguageFileTemplate = `
{
	"hello":"%s",
	"test":"%s",
	"optional":"%s"
}
`
const jsonLanguageFileTemplate2 = `
{
	"hello":"%s",
	"test":"%s"
}
`

const jsonBrokenFile = `
{
	"hello":"bla-bla",
	"test":
}
`

var helloWorldValueEn string
var testValueEn string
var optionalValueEn string
var helloWorldValueRu string
var testValueRu string
var jsonFilesFolder string

func TestMain(m *testing.M) {
	errSetup := setup()

	if errSetup != nil {
		log.Printf("setup failure: %s", errSetup.Error())
		os.Exit(1)
		return
	}

	errInit := InitLocales(jsonFilesFolder)
	if errInit != nil {
		log.Printf("init failure: %s", errInit.Error())
		os.Exit(1)
		return
	}

	retCode := m.Run()

	errFinish := teardown()
	if errFinish != nil {
		log.Printf("teardown failure: %s", errFinish.Error())
		os.Exit(1)
		return
	}

	os.Exit(retCode)
}

func TestEnTranslations(t *testing.T) {
	enHello := Tr("en", "hello")
	if enHello != helloWorldValueEn {
		t.Errorf("expected %s but retrieved %s", helloWorldValueEn, enHello)
	}
	enTest := Tr("en", "test")
	if enTest != testValueEn {
		t.Errorf("expected %s but retrieved %s", testValueEn, enTest)
	}
}

func TestRuTranslations(t *testing.T) {
	ruHello := Tr("ru", "hello")
	if ruHello != helloWorldValueRu {
		t.Errorf("expected %s but retrieved %s", helloWorldValueRu, ruHello)
	}
	ruTest := Tr("ru", "test")
	if ruTest != testValueRu {
		t.Errorf("expected %s but retrieved %s", testValueRu, ruTest)
	}
}

func TestMissingTranslations(t *testing.T) {
	ruOptional := Tr("ru", "optional")
	if ruOptional != optionalValueEn {
		t.Errorf("expected %s but retrieved %s", optionalValueEn, ruOptional)
	}
}

const nonExistingTag = "non_existing"

func TestNonExistingTranslations(t *testing.T) {
	ruNonExisting := Tr("ru", nonExistingTag)
	if ruNonExisting != nonExistingTag {
		t.Errorf("expected %s but retrieved %s", nonExistingTag, ruNonExisting)
	}
}

func TestCheckExistingLocales(t *testing.T) {
	if !checkLocale("en") {
		t.Error("locale en is not found")
	}
	if !checkLocale("ru") {
		t.Error("locale en is not found")
	}
}

func TestCheckNonExistingLocales(t *testing.T) {
	if checkLocale("de") {
		t.Error("locale de is found, but it does not exist")
	}
}

func TestLoadNonExistingFile(t *testing.T) {
	err := loadFileToMap(generateRandomString())
	if err == nil {
		t.Error("expected file load failure, retrieved no error")
	}
}

func TestLoadTranslationFromNonExistingFolder(t *testing.T) {
	err := loadTranslations(generateRandomString())
	if err == nil {
		t.Error("expected translation load failure, retrieved no error")
	}
}

const cAcceptLanguageSample1 = "ru, en-GB;q=0.8, en;q=0.7"
const cAcceptLanguageSample2 = "de, en-GB;q=0.8, en;q=0.7"
const cAcceptLanguageSample3 = "de, pt-BR;q=0.8"

func TestDetectMainLanguage(t *testing.T) {
	langDetected := DetectLanguage(cAcceptLanguageSample1)
	if langDetected != "ru" {
		t.Errorf("expected 'ru' but retrieved '%s'", langDetected)
	}
}
func TestDetectSecondaryLanguage(t *testing.T) {
	langDetected := DetectLanguage(cAcceptLanguageSample2)
	if langDetected != "en" {
		t.Errorf("expected 'en' but retrieved '%s'", langDetected)
	}
}
func TestDetectNotSupportedLanguage(t *testing.T) {
	langDetected := DetectLanguage(cAcceptLanguageSample3)
	if langDetected != "en" {
		t.Errorf("expected 'en' but retrieved '%s'", langDetected)
	}
}

func TestLoadBrokenFile(t *testing.T) {

	brokenFilePath := generateTempFolder()
	log.Println("full path to broken language json files: " + brokenFilePath)

	if _, err := os.Stat(brokenFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(brokenFilePath, os.ModePerm)
		if err != nil {
			t.Error("broken folder creation failure")
			return
		}
	}

	brokenJsonBytes := []byte(jsonBrokenFile)
	err := ioutil.WriteFile(brokenFilePath+"/en.json", brokenJsonBytes, 0644)
	if err != nil {
		t.Error("broken file creation failure")
		return
	}

	err = loadFileToMap(brokenFilePath + "/en.json")
	if err == nil {
		t.Error("expected file load failure due to broken json file, retrieved no error")
	}
	_ = os.RemoveAll(brokenFilePath)
}

func TestLoadBrokenTranslations(t *testing.T) {

	brokenFilePath := generateTempFolder()
	log.Println("full path to broken language json files: " + brokenFilePath)

	if _, err := os.Stat(brokenFilePath); os.IsNotExist(err) {
		err = os.MkdirAll(brokenFilePath, os.ModePerm)
		if err != nil {
			t.Error("broken folder creation failure")
			return
		}
	}

	brokenJsonBytes := []byte(jsonBrokenFile)
	err := ioutil.WriteFile(brokenFilePath+"/en.json", brokenJsonBytes, 0644)
	if err != nil {
		t.Error("broken file creation failure")
		return
	}

	err = loadTranslations(brokenFilePath)
	if err == nil {
		t.Error("expected file load failure due to broken json file, retrieved no error")
	}
	_ = os.RemoveAll(brokenFilePath)
}

func setup() error {
	log.Println("Tests initial setup")

	jsonFilesFolder = generateTempFolder()
	log.Println("full path to language json files: " + jsonFilesFolder)

	if _, err := os.Stat(jsonFilesFolder); os.IsNotExist(err) {
		err = os.MkdirAll(jsonFilesFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}

	helloWorldValueEn = generateRandomString()
	testValueEn = generateRandomString()
	optionalValueEn = generateRandomString()
	enJSON := fmt.Sprintf(jsonLanguageFileTemplate, helloWorldValueEn, testValueEn, optionalValueEn)
	enJsonBytes := []byte(enJSON)
	err := ioutil.WriteFile(jsonFilesFolder+"/en.json", enJsonBytes, 0644)
	if err != nil {
		return err
	}

	helloWorldValueRu = generateRandomString()
	testValueRu = generateRandomString()
	ruJSON := fmt.Sprintf(jsonLanguageFileTemplate2, helloWorldValueRu, testValueRu)
	ruJsonBytes := []byte(ruJSON)
	err = ioutil.WriteFile(jsonFilesFolder+"/ru.json", ruJsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func teardown() error {
	return os.RemoveAll(jsonFilesFolder)
}

func generateTempFolder() string {
	fullPathDBFile := os.TempDir()
	tempFileName := generateRandomString()
	if strings.HasSuffix(fullPathDBFile, "/") {
		fullPathDBFile += tempFileName
	} else {
		fullPathDBFile += "/" + tempFileName
	}
	return fullPathDBFile
}

func generateRandomString() string {
	length := 8
	rb := make([]byte, length)
	_, _ = rand.Read(rb)

	b64 := base64.URLEncoding.EncodeToString(rb)
	if utf8.RuneCountInString(b64) < length {
		length = utf8.RuneCountInString(b64)
	}
	return b64[0:length]
}
