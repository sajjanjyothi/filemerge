package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// Solution
// ========
// Read the file sequentially, process N records at a time in memory (N is arbitrary based on memory) - here it is MAX_LINES_FILE
// Sort the N records in memory, write them to a temp file. Loop on the file until it is done.
// Open all the temp files at the same time, but read only one record per file.
// For each of these temp file records, find the smaller, write it to the final file, and advance only in that file.

// This can be calibrated based on the memory restriction and prefomrmance tuning- number of records read to the memory
// Since memory restriction is there we need more disk based I/O here

const MAX_LINES_FILE int = 10

func cleanupFiles() {
	files, err := filepath.Glob("sorttemp*")
	if err == nil {
		for _, file := range files {
			os.Remove(file)
		}
	}
}

func mergeFiles(currentFile string, prevFile string, tempFileName string) {

	file1, _ := os.Open(prevFile)
	file2, _ := os.Open(currentFile)
	file3, _ := os.Create("sorttemp.tmp")
	defer file1.Close()
	defer file2.Close()
	defer file3.Close()

	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)
	scanner1.Scan()
	scanner2.Scan()
Loop:
	for {
		val1, err1 := strconv.Atoi(scanner1.Text())
		val2, err2 := strconv.Atoi(scanner2.Text())
		switch {
		case err1 == nil && err2 == nil:
			if val1 <= val2 {
				fmt.Fprintln(file3, val1)
				scanner1.Scan() // advance only corresponding file reader
			} else {
				fmt.Fprintln(file3, val2)
				scanner2.Scan() // advance only corresponding file reader
			}
		case err1 == nil:
			fmt.Fprintln(file3, val1)
			scanner1.Scan() // advance only corresponding file reader
		case err2 == nil:
			fmt.Fprintln(file3, val2)
			scanner2.Scan() // advance only corresponding file reader
		case err1 != nil && err2 != nil:
			break Loop
		}
	}
	os.Rename("sorttemp.tmp", tempFileName) //Just rename the sorted temporary file
}

func main() {

	if len(os.Args) < 3 {
		panic("Not enough command line arguements.....")
	}

	cleanupFiles() //clean up files from previous run
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		panic(err)
	}

	tempFileName := os.Args[2]
	scanner := bufio.NewScanner(file)

	fmt.Println("Sorting.....Please wait")

	lines := 0
	fileCount := 0
	dataArray := make([]int, 0)
	fileList := make([]string, 0)

	for scanner.Scan() {
		lines++
		data, _ := strconv.Atoi(scanner.Text())
		dataArray = append(dataArray, data)

		if lines >= MAX_LINES_FILE { //read only max line records because of memory restrictions
			fileCount++
			tmpFileName := fmt.Sprintf("sorttemp-%d", fileCount)
			fileList = append(fileList, tmpFileName)
			f, _ := os.Create(tmpFileName)
			sort.Ints(dataArray)
			for _, value := range dataArray {
				fmt.Fprintln(f, value)
			}
			f.Close()
			dataArray = dataArray[:0] //Clear the slice
			lines = 0
		}
	}

	if len(dataArray) > 0 { // read the rest of the records
		fileCount++
		tmpFileName := fmt.Sprintf("sorttemp-%d", fileCount)
		fileList = append(fileList, tmpFileName)
		f, _ := os.Create(tmpFileName)
		defer f.Close()
		sort.Ints(dataArray)
		for _, value := range dataArray {
			fmt.Fprintln(f, value)
		}
		dataArray = dataArray[:0]
	}

	prevFile := ""
	for _, file := range fileList {
		if prevFile != "" {
			mergeFiles(file, tempFileName, tempFileName)
		} else {
			os.Rename(file, tempFileName)
			prevFile = file
		}
	}

	fmt.Println("File Sorted....")
	cleanupFiles()
}
