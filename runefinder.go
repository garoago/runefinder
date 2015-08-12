package main

import (
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

const ucdFileName = "UnicodeData.txt"
const ucdBaseUrl = "http://www.unicode.org/Public/UCD/latest/ucd/"

func progress(done *bool) {
	for !*done {
		fmt.Print(".")
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println()
}

func getUcdFile(fileName string) {
	url := ucdBaseUrl + ucdFileName
	fmt.Printf("%s not found\nretrieving from %s\n", ucdFileName, url)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	file.Close()
}

func loadIndex(fileName string) (map[string][]rune, map[rune]string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		done := false
		go progress(&done)
		getUcdFile(fileName)
		done = true
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(content), "\n")

	index := make(map[string][]rune)
	names := make(map[rune]string)

	for _, line := range lines {
		var code rune
		fields := strings.Split(line, ";")
		if len(fields) >= 2 {
			code64, _ := strconv.ParseInt(fields[0], 16, 0)
			code = rune(code64)
			names[code] = fields[1]
			for _, word := range strings.Split(fields[1], " ") {
				var entries []rune
				if len(index[word]) < 1 {
					entries = make([]rune, 0)
				} else {
					entries = index[word]
				}
				index[word] = append(entries, code)
			}
		}

	}
	return index, names
}

func main() {
	dir, _ := os.Getwd()
	path := path.Join(dir, ucdFileName)
	index, names := loadIndex(path)
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	word := strings.ToUpper(os.Args[1])
	for _, code := range index[word] {
		fmt.Printf("%c %s\n", code, names[code])
	}

}
