package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func seedColumnIds() []string {
	columnIds := make([]string, 0)
	f, err := os.Open("columnIds.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		//fmt.Println(scanner.Text())
		columnIds = append(columnIds, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return columnIds
}

func main() {
	var wg sync.WaitGroup
	work := make(chan string)

	var columnIdsToDelete []string = seedColumnIds()
	fmt.Println(len(columnIdsToDelete))

	// write columnIds to the work channel, blocking until a worker goroutine
	// is able to start work
	go func(ids []string) {
		for _, columnId := range ids {
			fmt.Println(columnId)
			work <- columnId
		}
		// close channel to tell workers data is complete
		close(work)
	}(columnIdsToDelete)

	// Start up pool of goroutines and read columnIds from work channel.
	// Honeycomb doesn't explicitly list a rate limit for column operations,
	// but Events are limited to 2000/s for free tier, 7000/s for pro/enterprise.
	// https://intercom.help/honeycomb/en/articles/2121867-is-there-a-limit-to-how-fast-i-can-send-in-data
	// To compromise between speed for deleting up to thousands of columns and polite usage of the API,
	// we limit to 100 column deletions per second.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(w chan string) {
			defer wg.Done()
			for {
				select {
				case columnId, ok := <-w:
					// if channel is closed, no more requests need to be made
					if !ok {
						return
					}
					url := fmt.Sprintf("https://api.honeycomb.io/1/columns/<dataset>/%s", columnId)
					sendRequest(url)
					time.Sleep(1 * time.Second)
					fmt.Println(url)
				}
			}
		}(work)
	}

	wg.Wait()

	fmt.Println("Done")

}

func sendRequest(url string) {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("X-Honeycomb-Team", "<HC token>")

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
