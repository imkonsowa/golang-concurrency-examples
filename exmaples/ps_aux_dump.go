package exmaples

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type CsvWriter struct {
	writer *csv.Writer
	Mutex  *sync.Mutex
}

func NewCsvWriter(fileName string) *CsvWriter {
	file, err := os.Create(fileName + ".csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)

	return &CsvWriter{
		writer: writer,
		Mutex:  &sync.Mutex{},
	}
}

func (c *CsvWriter) Write(cells []string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	err := c.writer.Write(cells)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *CsvWriter) Flush() {
	c.writer.Flush()
}

// PSAUXDump is a simple concurrent procedure to dump linux command `ps aux` into CSV file
func PSAUXDump() {
	pidWorker := func(pids <-chan string, results chan<- string, w *CsvWriter) {
		for line := range pids {
			w.Write(strings.Fields(line))
			results <- fmt.Sprintf("Processed line: %s", line)
		}
	}

	writer := NewCsvWriter("ps-aux")

	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	// write the header
	writer.Write(strings.Fields(lines[0]))
	lines = lines[1:]

	pids := make(chan string, len(lines))
	results := make(chan string, len(lines))

	// spin up 10 workers
	for w := 0; w <= 10; w++ {
		go pidWorker(pids, results, writer)
	}

	for _, l := range lines {
		pids <- l
	}

	close(pids)

	for range lines {
		fmt.Println(<-results)
	}
	fmt.Println("done processing")
	writer.Flush()
	close(results)
}
