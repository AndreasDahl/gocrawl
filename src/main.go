package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)



func download(id int, c chan int) {
    resp, err := http.Get(fmt.Sprintf("http://www.pathofexile.com/forum/view-thread/%d", id))
    if err != nil {
    	// handle error
    }
    body, err := ioutil.ReadAll(resp.Body)
    ioutil.WriteFile(fmt.Sprintf("shops/%d.html", id), body, 0666)
    c <- id
}

func main() {
    const dl_count int = 5
    c := make(chan int)
    for i := 1588095; i > 1588095 - dl_count; i-- {
        go download(i, c)
    }

    for i := 0..; i < dl_count; i++ {
        fmt.Println("%d: done", <-c)
    }
}
