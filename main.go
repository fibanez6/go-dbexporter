package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/fibanez6/go-dbexporter/dbdriver"
	"github.com/fibanez6/go-dbexporter/domain"
	"github.com/fibanez6/go-dbexporter/service"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func logError(err error, line string) {
	log.Printf(`{ "err": %s, line: %s }`, err, line)
}

/*
It will stop the normal execution of the current goroutine if there  is an
error
*/
func check(e error, m string) {
	if e != nil {
		logError(e, m)
		panic(e)
	}
}

/*
Two Steps to Calculate PPI
 - Use the Pythagorean Theorem and the screen width and height in pixels to
   calculate the diagonal length in pixels.
 - Use the formula to calculate PPI, dividing the length of the diagonal in
   pixels by the length of the diagonal in inches
*/
func calculatePPI(line []string) (float64, error) {
	d, errD := strconv.ParseFloat(line[5], 64)
	if errD != nil {
		return 0, errD
	}
	h, errH := strconv.ParseFloat(line[6], 64)
	if errH != nil {
		return 0, errH
	}
	v, errV := strconv.ParseFloat(line[7], 64)
	if errV != nil {
		return 0, errV
	}
	diagonal := math.Sqrt(h*h+v*v) / d
	ppi := math.Round(diagonal*100000) / 100000
	return ppi, nil
}

/*
This is the entry point for the application. It imports the csv data in a SQL
database  and it will start processing  from the last saved offset, otherwise
from the beginning of the file.

In the case  of any error in  processing a line, it will  log it and continue
along the next line.

Considerations:
- Database and tables must exist.
*/
func main() {
	defer dbdriver.Close()

	dataPtr := flag.String("file", "data/data.csv", "Path to the CSV")
	offsetFlagPtr := flag.Bool("set-offset-beginning", true, "Flag to read the data file from the beginning, true by default")
	flag.Parse()

	offset := int64(0)
	if !*offsetFlagPtr {
		offset = offsetService.ReadLatestOffset()
	}

	csvFile, _ := os.Open(*dataPtr)
	_, errF := csvFile.Seek(offset, io.SeekStart)
	check(errF, "Error reading data")

	scanner := bufio.NewScanner(csvFile)
	// skipping first line
	if offset == 0 {
		scanner.Scan()
		offset = offsetService.MoveOffset(offset, scanner.Text())
	}

	counter := 0

	// starting to read line by line
	for scanner.Scan() {
		line := scanner.Text()
		offset = offsetService.MoveOffset(offset, line)

		parts := strings.Split(line, ",")
		if len(parts) != 8 {
			logError(errors.New("line size is wrong"), line)
			continue
		}

		ppi, err := calculatePPI(parts)
		if err != nil {
			logError(err, line)
			continue
		}

		monitor := domain.Monitor{
			SerialNumber: parts[3],
			Resolution:   ppi,
		}
		device := domain.Device{
			Name:          parts[0],
			LastIpAddress: parts[1],
		}

		// writing to the database
		errDBW := dbdriver.Write(device, monitor)
		if errDBW == nil {
			counter++
		}
		check(errDBW, line)
	}
	log.Printf(`Total number of rows: %d`, counter)

	// writing the latest offset to the file
	errW := offsetService.WriteOffset(offset)
	check(errW, fmt.Sprintf("the latest offset is %d", offset))
	log.Print("Successfully finished")
}
