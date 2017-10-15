package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jeffprestes/stopwords"
	"github.com/jeffprestes/stopwords/corpus"
)

func main() {
	var countingMap = make(map[string]int, 0)
	filePathOriginal := "/Users/jeffprestes/Downloads/conversation_report_2017-10-11.csv"
	originalFile, err := os.Open(filePathOriginal)
	if err != nil {
		panic("Original file couldn't be opened")
	}
	defer originalFile.Close()
	err = os.Remove("/Users/jeffprestes/Downloads/wordcloud.txt")
	if err != nil {
		fmt.Println(err)
	}
	destFile, err := os.Create("/Users/jeffprestes/Downloads/wordcloud.txt")
	if err != nil {
		panic("Could not create the destination file: " + err.Error())
	}
	defer destFile.Close()
	csvReader := csv.NewReader(originalFile)
	csvReader.Comma = ';'
	csvReader.Comment = '#'
	records, err := csvReader.ReadAll()
	if err != nil {
		panic("Couldn't read the record " + err.Error())
	}

	var totalRecords int
	var tmp string
	for index, record := range records {
		// if index > 1000 {
		// 	break
		// }
		tmp = corpus.RemoveDiacriticMark(record[2])
		tmp, err := stopwords.Filter(tmp, corpus.Portuguese)
		if err != nil {
			fmt.Println("Error at line ", index, " ", err.Error())
			continue
		}
		if len(tmp) > 0 {
			stringsInPhrase := strings.Split(tmp, " ")
			if len(stringsInPhrase) > 0 {
				//fmt.Println("Users said: ", stringsInPhrase)
				for _, word := range stringsInPhrase {
					counter, _ := countingMap[word]
					counter++
					countingMap[word] = counter
				}
				totalRecords++
			}
		}
	}
	fmt.Println("========================================")
	fmt.Println(" ")
	var destString string
	for wordFound, value := range countingMap {
		destString = strconv.Itoa(value) + " " + wordFound + "\n"
		//fmt.Println(destString)
		destFile.WriteString(destString)
	}
	err = destFile.Sync()
	if err != nil {
		panic("Error during flush content to destination file: " + err.Error())
	}
	fmt.Println(" ")
	fmt.Println("========================================")
	fmt.Println(" ")
	fmt.Println("Total records were: ", totalRecords)
}
