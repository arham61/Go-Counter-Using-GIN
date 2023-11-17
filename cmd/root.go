package cmd

import (
	"fmt"
	"gin/models"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
)

func FileData(c *gin.Context) {
	var data models.Data
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	totalTime, totalLines, totalWords, totalVowels, totalPunctuations := FileReader(data.Path, data.Routines)
	totalTimeString := totalTime.String()
	c.JSON(http.StatusOK, gin.H{"File Path": data.Path,
		"Go Routines":        data.Routines,
		"Execution Time":     totalTimeString,
		"Total Words ":       totalWords,
		"Total Lines":        totalLines,
		"Total Vowels":       totalVowels,
		"Total Punctuations": totalPunctuations})
}

func FileReader(path string, routines int) (time.Duration, int, int, int, int) {
	var Lines, Words, Vowels, Punctuations int
	fileContent, err := Read(path)
	if err != nil {
		fmt.Print("\nError : ", err)
	}
	startTime := time.Now()

	channel := make(chan models.Counter)

	chunk := len(fileContent) / routines

	fmt.Printf("File Chunks = %v ", routines)
	fmt.Printf("\n")

	for i := 0; i < routines; i++ {
		start := i * chunk
		end := (i + 1) * chunk
		go Count(fileContent[start:end], channel)
	}

	for i := 0; i < routines; i++ {
		Counts := <-channel
		fmt.Printf("No of Words of Chunk %d: %d \n", i+1, Counts.Words)
		fmt.Printf("No of Lines of Chunk %d: %d \n", i+1, Counts.Lines)
		fmt.Printf("No of Vowels of Chunk %d: %d \n", i+1, Counts.Vowels)
		fmt.Printf("No of Punctuation of Chunk %d: %d \n", i+1, Counts.Punctuations)
		fmt.Printf("\n\n")
		Lines = Lines + Counts.Lines
		Words = Words + Counts.Words
		Vowels = Vowels + Counts.Vowels
		Punctuations = Punctuations + Counts.Punctuations
	}
	return time.Since(startTime), Lines, Words, Vowels, Punctuations
}

func Read(path string) (string, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	Content := string(fileContent)
	return Content, nil
}

func Count(fileContent string, channel chan models.Counter) {
	count := models.Counter{}
	for _, char := range fileContent {
		switch {
		case char == ' ' || char == '\t' || char == '\r' || char == '.' || char == ',' || char == ';' || char == ':' || char == '!' || char == '?' || char == '(' || char == ')' || char == '[' || char == ']' || char == '{' || char == '}':
			count.Words++
		case char == 'A' || char == 'E' || char == 'I' || char == 'O' || char == 'U' || char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u':
			count.Vowels++
		case char == '.' || char == '!' || char == '?' || char == ',' || char == ':' || char == ';' || char == '(' || char == ')' || char == '[' || char == ']' || char == '{' || char == '}':
			count.Punctuations++
		case char == '\n':
			count.Lines++
		}
	}
	channel <- count
}
