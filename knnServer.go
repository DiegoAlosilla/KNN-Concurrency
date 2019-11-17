package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

const myIp = "10.0.75.1"

type Info struct {
	Tipo     string
	NodeNum  int
	NodeAddr string
	Class    string
}

type MyInfo struct {
	cont     int
	first    bool
	nextNum  int
	nextAddr string
}

var chMyInfo chan MyInfo
var readyToStart chan bool

var addrs []string
var myNum int

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

var class string

func main() {

	rand.Seed(time.Now().UnixNano())
	parseCSVFile("data.csv")

	var s Student

	min := 1
	max := 90
	maxY := 16
	s.ageOld = float64(rand.Intn((max - min) + min))
	s.yearStudy = float64(rand.Intn((maxY - min) + min))
	s.degree = ""

	var k int = 10

	class = getResponse(getNeightbors(students, s, k))

	rand.Seed(time.Now().UTC().UnixNano())
	myNum = rand.Intn(int(1e6))
	fmt.Println(myNum)
	var n int
	fmt.Print("Ingrese la cantidad de nodos: ")
	fmt.Scanf("%d\n", &n)
	addrs = make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Ingrese nodo %d: ", i+1)
		fmt.Scanf("%s\n", &(addrs[i]))
	}
	readyToStart = make(chan bool)
	go func() {
		chMyInfo = make(chan MyInfo)
		chMyInfo <- MyInfo{0, true, int(1e7), ""}
	}()
	go func() {
		gin := bufio.NewReader(os.Stdin)
		fmt.Print("Presione enter para iniciar...")
		gin.ReadString('\n')
		info := Info{"SENDNUM", myNum, myIp, class}
		for _, addr := range addrs {
			send(addr, info)
		}
	}()
	server()
}
func server() {
	host := fmt.Sprintf("%s:8000", myIp)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}
func handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')
	var info Info
	json.Unmarshal([]byte(msg), &info)
	fmt.Println(info)
	switch info.Tipo {
	case "SENDNUM":
		myInfo := <-chMyInfo
		myInfo.cont++
		if info.NodeNum < myNum {
			myInfo.first = false
		} else if info.NodeNum < myInfo.nextNum {
			myInfo.nextNum = info.NodeNum
			myInfo.nextAddr = info.NodeAddr
		}
		go func() {
			chMyInfo <- myInfo
		}()
		if myInfo.cont == len(addrs) {
			if myInfo.first {
				fmt.Println(class)
				criticalSection()
			} else {
				readyToStart <- true
			}
		}
	case "START":
		<-readyToStart
		criticalSection()
	}
}
func criticalSection() {
	fmt.Println(class)
	myInfo := <-chMyInfo
	if myInfo.nextAddr == "" {
		fmt.Println(class)
	} else {
		info := Info{Tipo: "START"}
		fmt.Println(myInfo, info)
		send(myInfo.nextAddr, info)
	}
}
func send(remoteAddr string, info Info) {
	remote := fmt.Sprintf("%s:8000", remoteAddr)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	bytesMsg, _ := json.Marshal(info)
	fmt.Fprintln(conn, string(bytesMsg))
}
