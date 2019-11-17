package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type Student struct {
	ageOld    float64
	yearStudy float64
	degree    string
}

type Distance struct {
	distance float64
	student  Student
}

var students []Student

//https://www.thepolyglotdeveloper.com/2017/03/parse-csv-data-go-programming-language/
func parseCSVFile(filePath string) {
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var student Student
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		student.ageOld, _ = strconv.ParseFloat(line[0], 32)
		student.yearStudy, _ = strconv.ParseFloat(line[1], 32)
		student.degree = line[2]
		students = append(students, student)
	}
}

func euclideanDistance(student1 Student, student2 Student) (d float64) {
	d += math.Sqrt(math.Pow((student1.ageOld-student2.ageOld), 2) + math.Pow((student1.yearStudy-student2.yearStudy), 2))
	return
}

type ByDistance []Distance

func (s ByDistance) Len() int           { return len(s) }
func (s ByDistance) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByDistance) Less(i, j int) bool { return s[i].distance < s[j].distance }

func getNeightbors(trainingSet []Student, testInstance Student, k int) (neightbors []Distance) {
	var distances []Distance
	var dist float64

	for i := 0; i < len(trainingSet); i++ {
		dist = euclideanDistance(testInstance, trainingSet[i])
		distances = append(distances, Distance{
			distance: dist,
			student:  trainingSet[i],
		})
	}
	sort.Sort(ByDistance(distances))

	for i := 0; i < k; i++ {
		neightbors = append(neightbors, distances[i])
	}
	return neightbors
}

func getResponse(neightbors []Distance) string {

	classVote := make(map[string]int)
	var item string
	for i := 0; i < len(neightbors); i++ {
		item = neightbors[i].student.degree
		_, response := classVote[item]
		if response {
			classVote[item] += 1
		} else {
			classVote[item] = 1
		}
	}

	keys := []string{}

	// iterate over the map and append all keys to our
	// string array of keys
	for key := range classVote {
		keys = append(keys, key)
	}

	// use the sort method to sort our keys array
	sort.Strings(keys)

	return keys[0]
}

func main() {

	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 90
	maxY := 16

	parseCSVFile("data.csv")
	var s Student
	s.ageOld = float64(rand.Intn((max - min) + min))
	s.yearStudy = float64(rand.Intn((maxY - min) + min))
	s.degree = ""
	var k int = 10
	fmt.Println(getNeightbors(students, s, k))
	fmt.Println(getResponse(getNeightbors(students, s, k)))
}
