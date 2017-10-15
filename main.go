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
	var totalRecords int
	var tmp string

	//Define your file paths here
	var filePathOriginal = "/Users/jeffprestes/Downloads/conversation_report_2017-10-11.csv"
	var filePathDestination = "/Users/jeffprestes/Downloads/wordcloud.txt"

	originalFile, err := os.Open(filePathOriginal)
	if err != nil {
		panic("Original file couldn't be opened")
	}
	defer originalFile.Close()

	//Remove old destination file if exists
	err = os.Remove(filePathDestination)
	if err != nil {
		fmt.Println(err)
	}

	destFile, err := os.Create(filePathDestination)
	if err != nil {
		panic("Could not create the destination file: " + err.Error())
	}
	defer destFile.Close()

	csvReader := csv.NewReader(originalFile)
	//This programs assumes you use semi-colon as column separator. You can change it here.
	csvReader.Comma = ';'
	csvReader.Comment = '#'
	records, err := csvReader.ReadAll()
	if err != nil {
		panic("Couldn't read the record " + err.Error())
	}

	//If your CSV contains several columns, below you define which column has the text to be analysed
	var columnToRead = 2

	for index, record := range records {
		//If you want to define a limit to read the file content you can uncomment here
		// if index > 1000 {
		// 	break
		// }
		tmp = corpus.RemoveDiacriticMark(record[columnToRead])

		//You can change the corpus set from Portuguese to English, Spanish and several other languages
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
		//Below format are compatible with https://www.wordclouds.com/ . The result you must copy to it using 'Word list'
		//button.
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
