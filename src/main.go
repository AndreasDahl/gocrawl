package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

func getShops() (ids [2]int) {
	resp, err := http.Get("http://www.pathofexile.com/forum/view-forum/standard-trading-shops")
	if err != nil {
		// handle error
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", body)
	return [2]int{1, 2}
}

func main() {
	getShops()
	// const dlCount int = 5
	// c := make(chan int)
	// for i := 1588095; i > 1588095-dlCount; i-- {
	//     go download(i, c)
	// }
	//
	// for i := 0; i < dlCount; i++ {
	//     fmt.Printf("done: %d\n", <-c)
	// }
}
