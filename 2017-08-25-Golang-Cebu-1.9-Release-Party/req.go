package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	lines := make(chan *string, 5)

	go readStdinIntoChannel(lines)
	wg.Add(1)
	go getFromUrlAndPrintStatus(lines, &wg)
	wg.Wait()
}

func readStdinIntoChannel(lines chan<- *string) {
	defer close(lines)
	reader := bufio.NewReader(os.Stdin)

	for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
		newline := line // need new memory location
		lines <- &newline
	}
}

func getFromUrlAndPrintStatus(lines <-chan *string, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range lines {
		wg.Add(1)
		go getAndPrintStatus(line, wg)
	}
}

func getAndPrintStatus(line *string, wg *sync.WaitGroup) {
	defer wg.Done()

	url := strings.TrimSpace(*line)

	fmt.Printf("     (requesting %v)\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: %v", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("%v   %v\n", resp.StatusCode, url)
}
