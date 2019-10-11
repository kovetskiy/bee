package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/docopt/docopt-go"
)

var (
	version = "[manual build]"
	usage   = "bee " + version + `

bee makes your text look like a bee

Usage:
  bee [options] [--] [<cmd>...]
  bee -h | --help
  bee --version

Options:
  -h --help  Show this screen.
  -t <time>  Time to make it look like a bee (ms). [default: 1000]
  --version  Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	ms, err := strconv.Atoi(args["-t"].(string))
	if err != nil {
		log.Fatalln(err)
	}

	timeout := time.Millisecond * time.Duration(ms)

	cmdline, _ := args["<cmd>"].([]string)
	if len(cmdline) > 0 {
		cmd := exec.Command(cmdline[0], cmdline[1:]...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalln(err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatalln(err)
		}
		cmd.Stdin = os.Stdin
		err = cmd.Start()
		if err != nil {
			log.Fatalln(err)
		}

		work := &sync.WaitGroup{}

		work.Add(3)

		go func() {
			defer work.Done()
			err := bee(stdout, os.Stdout, timeout)
			if err != nil {
				log.Fatalln(err)
			}
		}()

		go func() {
			defer work.Done()
			err := bee(stderr, os.Stderr, timeout)
			if err != nil {
				log.Fatalln(err)
			}
		}()

		go func() {
			defer work.Done()
			err := cmd.Wait()
			if err != nil {
				log.Fatalln(err)
			}
		}()

		work.Wait()
		return
	}

	err = bee(os.Stdin, os.Stdout, timeout)
	if err != nil {
		log.Fatalln(err)
	}
}

func bee(input io.Reader, output io.Writer, timeout time.Duration) error {
	var text string
	flushed := make(chan struct{}, 1)
	highlighted := false

	go func() {
		for {
			after := time.After(timeout)
			select {
			case <-after:
				if highlighted {
					break
				}

				err := highlight(output, text)
				if err != nil {
					log.Fatalln(err)
				}

				highlighted = true
			case <-flushed:
				highlighted = false
				break
			}
		}
	}()

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		text = scanner.Text()
		if text[len(text)-1] == '\n' {
			text = text[:len(text)-1]
		}
		_, err := output.Write([]byte("\n" + text))
		if err != nil {
			return err
		}
		flushed <- struct{}{}
	}
	return nil
}

func highlight(output io.Writer, text string) error {
	text = "\r\x1b[43m" + text + "\x1b[0m"
	_, err := output.Write([]byte(text))
	return err
}
