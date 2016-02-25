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
		fmt.Printf("ERR: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(fmt.Sprintf("shops/%d.html", id), body, 0666)
	if err != nil {
		fmt.Printf("ERR: %s", err)
	}
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

func getItemJSON(shopID int) {
	expr, _ := regexp.Compile("new R\\((.*)\\);")
	shopHTML, err := ioutil.ReadFile(fmt.Sprintf("shops/%d.html", shopID))
	if err != nil {
		fmt.Printf("ERR: %s", err)
	}
	matches := expr.FindSubmatch(shopHTML)
	if matches == nil {
		fmt.Printf("shop '%d' has no items\n", shopID)
	} else {
		ioutil.WriteFile(fmt.Sprintf("shops/%d.json", shopID), matches[1], 0666)
	}
}

func bytesAsInt(bts []byte) (i int) {
	i = 0
	for _, b := range bts {
		i *= 10
		i += int(b) - 48
	}
	return
}

func getShops(page int) {
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
}

func main() {
	// const pageCount = 1
	// for i := 1; i < pageCount+1; i++ {
	// 	getShops(i)
	// }

	getItemJSON(389691)

	fmt.Println("FINISHED")
}

// frameType property:
// #  0 Normal Item
// #  1 Magic Item
// #  2 Rare Item
// #  3 Unique Item
// #  4 Gem
// #  5 Currency
// #  6 Quest Item
