package offsetService

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

/*
This file reads and stores the latest offset of the csv
*/

const (
	offsetFileName = "offset.txt"
	endLineSize    = int64(len("\\\n"))
)

/*
Reads the offset from 'offset.txt', if the file does not exist,
then it will return offset=0.
*/
func ReadLatestOffset() int64 {
	file, errFile := os.Open(offsetFileName)
	if errFile != nil {
		log.Println(errFile)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		lastOffset, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			return lastOffset
		}
	}
	return 0
}

/*
Moves the offset len(string) positions
*/
func MoveOffset(offset int64, line string) int64 {
	return offset + int64(len(line)) + endLineSize
}

/*
Writes the offset to 'offset.txt', if the file does not exist,
then it will create a new file.
*/
func WriteOffset(offset int64) error {
	var file *os.File
	var err error

	if _, err = os.Stat(offsetFileName); err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(offsetFileName)
		}
	}

	file, err = os.OpenFile(offsetFileName, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", offset))
	return err
}
