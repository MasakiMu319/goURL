package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"
)

var (
	// Command line flags
	httpMethod       string // http method
	httpResponseHead bool   // response head
	httpConnectInfo  bool   // connect information

	showVersion bool	// show program version

	version = "Dev"
)

func init() {
	flag.StringVar(&httpMethod, "X", "GET", "HTTP method to use")
	flag.BoolVar(&httpResponseHead, "I", false, "show response head and source code of page")
	flag.BoolVar(&httpConnectInfo, "v", false, "show connect process")
	flag.BoolVar(&showVersion, "V", false, "show goURL version")
	flag.Usage = usage
}

func usage()  {
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] URL\n\n", os.Args[0])
	_, _ = fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
}

func main() {
	// parse command-line flags from os.Args[1:]
	flag.Parse()

	if showVersion {
		if version == "Dev" {
			// print information with red color !
			color.HiRed("This is a %s version! Please do not use in a product environment! \n (runtime: %s)\n\n", version, runtime.Version())
		} else {
			fmt.Printf("goURL version: %s \n(runtime: %s)\n\n", version, runtime.Version())
		}
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(2)
	}
}
