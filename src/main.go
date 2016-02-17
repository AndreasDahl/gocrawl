package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

func getShopIds(url string) [][]byte {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	expr, _ := regexp.Compile("thread_title\\D+(?P<id>\\d+)")
	matches := expr.FindAllSubmatch(body, -1)
	retval := make([][]byte, len(matches))
	for i, e := range matches {
		retval[i] = e[1]
	}

	return retval
}

func bytesAsInt(bts []byte) (i int) {
	i = 0
	for _, b := range bts {
		i *= 10
		i += int(b) - 48
	}
	return
}

func getShops(page int, snc chan int) {
	shopIds := getShopIds(
		fmt.Sprintf("http://www.pathofexile.com/forum/view-forum/standard-trading-shops/page/%d",
			page))
	c := make(chan int)
	for _, id := range shopIds {
		go download(bytesAsInt(id), c)
	}

	for range shopIds {
		fmt.Printf("done: %d\n", <-c)
	}

	snc <- page
}

func main() {
	const pageCount = 2
	c := make(chan int)
	for i := 0; i < pageCount; i++ {
		getShops(i, c)
	}

	for i := 0; i < pageCount; i++ {
		fmt.Printf("done: %d\n", <-c)
	}
}
