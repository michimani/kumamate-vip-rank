package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/michimani/kumamate-vip-rank/kvr"
)

const cachePath = "./testdata/cache"

func main() {
	os.Setenv("KVR_CACHE_FILE", cachePath)
	client, err := kvr.New(http.DefaultClient)
	if err != nil {
		panic(err)
	}

	vrs, err := client.ListVIPRank()
	if err != nil {
		panic(err)
	}

	for _, vr := range vrs {
		fmt.Printf("%s: %d 〜 ", vr.Title, vr.Min)
		if vr.Max > 0 {
			fmt.Printf("%d", vr.Max)
		}
		fmt.Println()
	}
}
