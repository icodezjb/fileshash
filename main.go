// fileshash project main.go
package main

import (
	"bufio"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const DEBUG = true

func file_sha1(filename string) string {
	var ret string

	file, err := os.Open(filename)
	if err != nil {
		return ret
	}
	defer file.Close()

	hash := sha1.New()
	_, erro := io.Copy(hash, file)
	if erro != nil {
		return ret
	}
	ret = fmt.Sprintf("%x", hash.Sum(nil))

	return ret
}

func recursion_dir(dirname string, pattern *string, outputWriter *bufio.Writer) {
	dir_list, e := ioutil.ReadDir(dirname)
	if e != nil {
		fmt.Printf("read the dir(%s) error\n", dirname)
		return
	}

	for i, v := range dir_list {
		isMatched, _ := filepath.Match(*pattern, v.Name())
		if isMatched {
			continue
		}

		var outputstring string
		path := fmt.Sprintf("%s/%s", dirname, v.Name())

		switch v.IsDir() {
		case true:
			if DEBUG {
				fmt.Printf("%s is directory\n", v.Name())
			}
			recursion_dir(path, pattern, outputWriter)
		case false:
			/* skip thess files*/
			if (".result" == v.Name()) || ((v.Mode() & os.ModeType) != 0) {
				continue
			}

			/* only deal with regular file */
			hashstring := file_sha1(path)
			if "" != hashstring {
				/*wins:\r\n, unix:\n*/
				outputstring = fmt.Sprintf("%d,%s,%s,%d,%t\r\n", i, v.Name(), hashstring, v.Size(), ((v.Mode() & os.ModeType) == 0))
				outputWriter.WriteString(outputstring)
				if DEBUG {
					fmt.Println(outputstring)
				}
			}
		}
	}

	return
}

func main() {

	var outputWriter *bufio.Writer
	var DIR string
	var RESULT string
	var PATTERN string

	flag.StringVar(&DIR, "d", "./", "the destination directory")
	flag.StringVar(&RESULT, "o", "./.result", "output path of the reslut file")
	flag.StringVar(&PATTERN, "i", "", "the ignore dirs or/and files")

	flag.Parse()

	if DEBUG {
		fmt.Println(DIR)
		fmt.Println(RESULT)
	}
	/* reclear .result file */
	outputfile, openerr := os.Create(RESULT)
	if openerr != nil {
		fmt.Println("open the result file fail!")
		return
	} else {
		outputWriter = bufio.NewWriter(outputfile)
		defer outputfile.Close()
	}

	isMatched, _ := filepath.Match(PATTERN, DIR)
	if isMatched {
		fmt.Println("ignore thr dir:%s", DIR)
	} else {
		recursion_dir(DIR, &PATTERN, outputWriter)
	}

	outputWriter.Flush()
}
