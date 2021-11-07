package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/fatih/color"
	"goURL/parser"
	"goURL/utils"
)

func init() {
	flag.StringVar(&utils.HttpMethod, "X", "GET", "HTTP method to use")
	flag.BoolVar(&utils.HttpResponseHead, "I", false, "show response head and source code of page")
	flag.BoolVar(&utils.HttpConnectInfo, "v", false, "show connect process")
	flag.BoolVar(&utils.ShowVersion, "V", false, "show goURL version")
	flag.Usage = usage
}

func usage()  {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] URL\n\n", os.Args[0])
	_, _ = fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
}

func main() {
	// parse command-line flags from os.Args[1:].
	flag.Parse()

	// show goURL version or warning.
	if utils.ShowVersion {
		if utils.Version == "Dev" {
			// print information with red color !
			color.HiRed("This is a %s version! Please do not use in a product environment! \n (runtime: %s)\n", utils.Version, runtime.Version())
		} else {
			fmt.Printf("goURL version: %s \n(runtime: %s)\n", utils.VisitURL, runtime.Version())
		}
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		log.Fatalf(color.HiRedString("Too few arguments"))
	}

	// parse url argument.
	url, err := parser.ParseURL(args[0])
	if err != nil {
		log.Fatalf(color.HiRedString("Something wrong while parsing url:" + err.Error()))
	}
	// do connect with target URL.
	err = utils.VisitURL(url)
	if err != nil {
		log.Fatalf(color.HiRedString(err.Error()))
	}
}
