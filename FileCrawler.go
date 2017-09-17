package main

import (
	"FileCrawler/search"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	termptr := flag.String("term", "", "The term you're searching for. Example: -term=\"find me\"")
	rootptr := flag.String("root", "./", "The starting folder for searching. Example -root=\"c:/Looky Here\"")
	extptr := flag.String("ext", "txt", "Up to 25 file extension(s) to search. Example -ext=txt,doc")
	caseptr := flag.String("case", "n", "y for case sensitive")
	regptr := flag.String("regexp", "", "will search by regexp instead of term and case. Example -regexp=(?i)^startswith")
	logptr := flag.String("log", "console", "Where to log names of files containing the search. Can use \"console\". Example: -log=C:/Logs/Log.txt or -log=console")
	flag.Parse()
	searchterm := *termptr
	regexpterm := *regptr
	if searchterm == "" && regexpterm == "" {
		fmt.Println("There is no search or regexp term! Example: FileCrawler -term=\"find me\" -regexp=^startswith")
		return
	}

	var reg *regexp.Regexp
	var err error
	if regexpterm != "" {
		reg, err = regexp.Compile(regexpterm)
		if err != nil {
			fmt.Println("regexp is not valid: ", regexpterm)
			os.Exit(1)
		}
	} else {
		caseflag := *caseptr
		//default is case-insensitive
		caseterm := "(?i)"
		if caseflag == "y" || caseflag == "Y" {
			caseterm = ""
		}
		reg, err = regexp.Compile(caseterm + searchterm)
		if err != nil {
			fmt.Println("term is not valid: ", searchterm)
			os.Exit(1)
		}
	}

	root := *rootptr
	logname := *logptr
	extflag := *extptr
	ext := strings.Split(extflag, ",")
	foundlimit := 250
	foundcount := 0
	searchfunc, err := search.Factory(ext, reg, logname, foundlimit, &foundcount, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = filepath.Walk(root, searchfunc)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Done!")
}
