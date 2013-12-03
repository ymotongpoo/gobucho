package main

import (
	"flag"
	"fmt"

	"github.com/ymotongpoo/gobucho"
)

var (
	browser = flag.Bool("browser", false, "Flag for using browser")
	episode = flag.Int("episode", -1, "Specify episode # of torumemo")
)

func main() {
	flag.Parse()

	if *episode < 0 {
		fmt.Println(bucho.LatestStatus())
	} else {
		err := bucho.Torumemo(*browser, *episode)
		if err != nil {
			panic(err)
		}
	}
}
