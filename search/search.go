package search

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func extensionsFactory(ext []string) ([]*regexp.Regexp, error) {
	//do case insensitive match for filename ending. Example: "(?i)\.txt$"
	if len(ext) > 25 {
		return nil, errors.New("Surpassed limit of 25 file extensions")
	}
	var backary [25]*regexp.Regexp
	extensions := backary[0:len(ext)]
	for i, extstr := range ext {
		if strings.HasPrefix(extstr, ".") {
			extstr = "(?i)\\" + extstr + "$"
		} else {
			extstr = "(?i)\\." + extstr + "$"
		}
		thisregexp, err := regexp.Compile(extstr)
		if err != nil {
			return nil, errors.New(extstr + " is not a valid extension")
		}
		extensions[i] = thisregexp
	}
	return extensions, nil
}

func isValidFile(filename string, ext []*regexp.Regexp) bool {
	for _, extension := range ext {
		if extension.MatchString(filename) {
			return true
		}
	}
	return false
}

//Factory returns search function for filepath.Walk that uses file extensions and regex to search for and logfile or console for output
func Factory(extensions []string, reg *regexp.Regexp, logname string, foundlimit int, foundcount *int, echoslice []string) (func(string, os.FileInfo, error) error, error) {
	ext, err := extensionsFactory(extensions)
	if err != nil {
		fmt.Println("Error parsing file extension: " + err.Error())
		return nil, err
	}
	return func(thispath string, finfo os.FileInfo, err error) error {
		if *foundcount > foundlimit {
			return fmt.Errorf("Found files exceeded the limit of %d", foundlimit)
		}
		if !finfo.IsDir() {
			if !isValidFile(finfo.Name(), ext) {
				return nil
			}

			testfile, err := os.Open(thispath)
			if err != nil {
				return err
			}
			defer testfile.Close()

			if reg.MatchReader(bufio.NewReader(testfile)) {
				*foundcount = *foundcount + 1
				if logname == "console" {
					fmt.Println(thispath)
				} else if logname == "echoback" {
					if *foundcount > foundlimit {
						return fmt.Errorf("Found files exceeded the limit of %d", foundlimit)
					}
					echoslice[*foundcount-1] = thispath
				} else {
					f, err := os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					defer f.Close()

					if _, err = f.WriteString(thispath + "\r\n"); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}, nil
}
