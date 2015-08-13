package main

import (
	//"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const indexFileName = "runefinder-index.gob"
const ucdFileName = "UnicodeData.txt"
const ucdBaseUrl = "http://www.unicode.org/Public/UCD/latest/ucd/"

func progressDisplay(running *bool) {
	for *running {
		fmt.Print(".")
		time.Sleep(150 * time.Millisecond)
	}
	fmt.Println()
}

func getUcdFile(ucdPath string) {
	url := ucdBaseUrl + ucdFileName
	fmt.Printf("%s not found\nretrieving from %s\n", ucdFileName, url)
	running := true
	go progressDisplay(&running)
	defer func() {
		running = false
	}()
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	file, err := os.Create(ucdPath)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func buildIndex(ucdPath string) (map[string][]rune, map[rune]string) {
	if _, err := os.Stat(ucdPath); os.IsNotExist(err) {
		getUcdFile(ucdPath)
	}
	content, err := ioutil.ReadFile(ucdPath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")

	index := make(map[string][]rune)
	names := make(map[rune]string)

	for _, line := range lines {
		var uchar rune
		fields := strings.Split(line, ";")
		if len(fields) >= 2 {
			code64, _ := strconv.ParseInt(fields[0], 16, 0)
			uchar = rune(code64)
			names[uchar] = fields[1]
			for _, word := range strings.Split(fields[1], " ") {
				var entries []rune
				if len(index[word]) < 1 {
					entries = make([]rune, 0)
				} else {
					entries = index[word]
				}
				index[word] = append(entries, uchar)
			}
		}

	}
	return index, names
}

func loadIndex() (map[string][]rune, map[rune]string) {
	dir, _ := os.Getwd()
	indexPath := path.Join(dir, indexFileName)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		fmt.Print("Index not found")
		ucdPath := path.Join(dir, ucdFileName)
		index, names := buildIndex(ucdPath)
	}
	return index, names
}

func main() {
	index, names := loadIndex()
	if len(os.Args) != 2 {
		fmt.Println("Usage:  runefinder <word>\texample: runefinder cat")
		os.Exit(1)
	}

	word := strings.ToUpper(os.Args[1])
	for _, uchar := range index[word] {
		fmt.Printf("U+%-5X %c \t%s\n", uchar, uchar, names[uchar])
	}

}
