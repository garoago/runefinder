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

func progressDisplay(done *bool) {
	for !*done {
		fmt.Print(".")
		time.Sleep(150 * time.Millisecond)
	}
	fmt.Println()
}

func getUcdFile(fileName string) {
	url := ucdBaseUrl + ucdFileName
	fmt.Printf("%s not found\nretrieving from %s\n", ucdFileName, url)
	done := false
	go progressDisplay(&done)
	defer func() {
		done = true
	}()
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
		getUcdFile(fileName)
	}
	content, err := ioutil.ReadFile(fileName)
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

func main() {
	dir, _ := os.Getwd()
	path := path.Join(dir, ucdFileName)
	index, names := loadIndex(path)
	if len(os.Args) != 2 {
		fmt.Println("Usage:  runefinder <word>\texample: runefinder cat")
		os.Exit(1)
	}

	word := strings.ToUpper(os.Args[1])
	for _, uchar := range index[word] {
		fmt.Printf("U+%-5X %c \t%s\n", uchar, uchar, names[uchar])
	}

}
