package main

import (
	flag "github.com/spf13/pflag"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type selpg_args struct {
	start_page, end_page, lines_per_page, page_type int
	inFilename,	printDest                       string
}

//参数绑定
var inputS = flag.IntP("start_page", "s", -1, "(Mandatory) Input Your start_page")
var inputE = flag.IntP("end_page", "e", -1, "(Mandatory) Input Your end_page")
var inputL = flag.IntP("lines_per_page", "l", 20, "(Optional) Choosing lines_per_page mode, enter lines_per_page")
var inputF = flag.BoolP("pageBreak", "f", false, "(Optional) Choosing pageBreaks mode")
var inputD = flag.StringP("printDest", "d", "default", "(Optional) Enter printing destination")

var prog_name string

func processArgs(selpg *selpg_args) {
	//判断是否有页码输入，缺省值为-1
	if *inputS == -1 || *inputE == -1 {
		fmt.Fprintf(os.Stderr, "\nError: --start_page(-s) and --end_page(-e) are necessary\n", prog_name)
		flag.PrintDefaults()
		os.Exit(1)
	}
	// handle mandatory arg
	selpg.start_page = *inputS
	selpg.end_page = *inputE
	selpg.lines_per_page = *inputL
	if *inputF == true {
		selpg.page_type = 'f'
	}
	selpg.printDest = *inputD

	if flag.NArg() >= 1 {
		if flag.NArg() > 1 {
			fmt.Fprintf(os.Stderr, "%v: You should have one file input\n", prog_name)
			os.Exit(1)
		}
		_, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		selpg.inFilename = flag.Arg(0)
	}

	//错误处理
	switch {
	case selpg.start_page < 1 || selpg.end_page < 1:
		fmt.Fprintf(os.Stderr, "/n The page number should be greater than 0\n", prog_name)
		os.Exit(1)
	case selpg.lines_per_page <= 1:
		fmt.Fprintf(os.Stderr, "/n The lines_per_page should be greater than 1\n", prog_name)
		os.Exit(1)
	case selpg.end_page < selpg.start_page:
		fmt.Fprintf(os.Stderr, "/n The end_page should be greater than start_page\n", prog_name)
		os.Exit(1)
	case selpg.page_type != 'l' && selpg.page_type != 'f':
		fmt.Fprintf(os.Stderr, "/n There are only two page_types for you to choose: lines_per_page and pageBreaks\n", prog_name)
		os.Exit(1)
	}
}

func processInput(selpg *selpg_args) {
	var inputReader *bufio.Reader
	var outputWriter *bufio.Writer
	var err error
	var cmd *exec.Cmd
	var stdin io.WriteCloser
	var file *os.File
	//输入
	if selpg.inFilename == "0" {
		inputReader = bufio.NewReader(os.Stdin)
	}else {
		file, err = os.Open(selpg.inFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		inputReader = bufio.NewReader(file)
	}
	//输出
	if selpg.printDest == "default" {
		outputWriter = bufio.NewWriter(os.Stdout)
	} else {
		cmd = exec.Command("lp", "-d", selpg.printDest)
		stdin, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	//begin two input & output loops
	lineCount, pageCount := 0, 1
	var line []byte
	for {
		if selpg.page_type == 'l' {
			line, err = inputReader.ReadBytes('\n')
		} else {
			line, err = inputReader.ReadBytes('\f')
		}
		if err != nil {
			break
		}
		if selpg.page_type == 'l' {
			lineCount++
			if lineCount > selpg.lines_per_page {
				lineCount = 1
				pageCount++
			}
		}
		if pageCount >= selpg.start_page && pageCount <= selpg.end_page {
			if selpg.printDest == "default" {
				outputWriter.Write(line)
				outputWriter.Flush()
			} else {
				_, err := io.WriteString(stdin, string(line))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
		}
		if selpg.page_type == 'f' {
			pageCount++
		}
	}
	if selpg.printDest != "default" {
		stdin.Close()
		stderr, _ := cmd.CombinedOutput()
		fmt.Fprintln(os.Stderr, string(stderr))
	}
}

func main() {
	//给程序的参数绑定缺省值
	selpg := selpg_args{
		start_page:  -1,
		end_page:    -1,
		lines_per_page:    20,
		page_type:   'l',
		inFilename: "0",
		printDest:  "default",
	}
	prog_name = os.Args[0]
	flag.Parse()
	processArgs(&selpg)
	processInput(&selpg)
}