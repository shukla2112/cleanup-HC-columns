package main

import (
    "time"
    "net/http"
    "io/ioutil"
    "sync"
    "bufio"
    "fmt"
    "log"
    "os"
)

func seedColumnIds() []string{
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
    work := make(chan string, 10000)

    var columnIdsToDelete []string = seedColumnIds()
    fmt.Println(len(columnIdsToDelete))

    // write columnIds to the work channel, blocking until a worker goroutine
    // is able to start work
    for _, columnId := range columnIdsToDelete {
        fmt.Println(columnId)
        work <- columnId
    }

    // startup pool of 10 go routines and read columnIds from work channel 
    for {
        for i := 0; i<=10; i++ {
            wg.Add(1)
            go func(w chan string) {
                defer wg.Done()
                columnId := <-w
                url := fmt.Sprintf("https://api.honeycomb.io/1/columns/vector_multitenant_prod_filtered/%s", columnId)
                sendRequest(url)
                time.Sleep(1*time.Second)
                fmt.Println(url)
            }(work)
        }
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
