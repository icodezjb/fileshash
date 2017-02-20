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
)

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

func main() {

	var outputWriter *bufio.Writer
	var DIR string
	flag.StringVar(&DIR, "dir", "./", "the destination directory")
	flag.Parse()

	dir_list, e := ioutil.ReadDir(DIR)
	if e != nil {
		fmt.Printf("read the dir(DIR) error\n")
		return
	}
	fmt.Println(DIR)
	if nil == os.Chdir(DIR) {
		/* reclear .result file */
		outputfile, openerr := os.Create(".result")
		if openerr != nil {
			fmt.Println("open the result file fail!")
			return
		} else {
			outputWriter = bufio.NewWriter(outputfile)
			defer outputfile.Close()
		}

	} else {
		fmt.Printf("can't enter this directory:%s\n ")
		return
	}

	for i, v := range dir_list {

		var outputstring string

		switch v.IsDir() {
		case true:
			fmt.Printf("%s is directory\n", v.Name())
			break
		case false:
			/* skip .result file */
			if ".result" == v.Name() {
				continue
			}

			hashstring := file_sha1(v.Name())
			if "" != hashstring {
				/*wins:\r\n, unix:\n*/
				outputstring = fmt.Sprintf("%d,%s,%s,%d\r\n", i, v.Name(), hashstring, v.Size())
				outputWriter.WriteString(outputstring)
				//fmt.Println(outputstring)
			}
		}
	}

	outputWriter.Flush()
}
