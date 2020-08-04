package main

import (
	"encoding/hex"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	switch cmd {
	case "upload":
		exec := upload{}
		r := reader()
		defer r.Close()
		data, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatalln(err)
		}
		ret := exec.run(data)
		log.Printf("file hash will be shown below, use the hash to download")
		n, err := hex.NewEncoder(os.Stdout).Write(ret)
		os.Stdout.Write([]byte("\n"))
		if err != nil {
			log.Fatalln(err, n)
		}
	case "download":
		exec := download{}
		w := writer()
		defer w.Close()
		data, err := hex.DecodeString(hash)
		if err != nil {
			log.Fatalln(err)
		}

		ret := exec.run(data)
		n, err := w.Write(ret)
		if err != nil {
			log.Fatalln(err, n)
		}
	default:
		log.Fatalln("unsupported command:", cmd)
	}
}

func init() {
	flag.StringVar(&cmd, "c", "download", "upload / download")
	flag.StringVar(&filename, "f", "-", "filename, `-` for stdin / stdout")
	flag.StringVar(&hash, "s", "", "hash")

	flag.IntVar(&threads, "threads", 16, "thread num")
	flag.Parse()
}

func reader() io.ReadCloser {
	if filename == "-" {
		return os.Stdin
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func writer() io.WriteCloser {
	if filename == "-" {
		return os.Stdout
	}
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

var (
	cmd      string
	filename string
	hash     string

	threads int
)
