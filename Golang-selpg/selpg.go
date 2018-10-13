package main

import(
	"fmt"
	"os"
	flag "github.com/spf13/pflag"
	"bufio"
	"io"
	"os/exec"
	"errors"
)

var(
	start_page int
	end_page int
	page_length int
	in_filename string
	print_dest string
	page_type bool
	help bool
)

func init(){
	flag.IntVarP(&start_page, "start", "s",-1, "Define start page of select pages")
	flag.IntVarP(&end_page, "end", "e",-1, "Define end page of select pages")
	flag.IntVarP(&page_length, "pl", "l", 72, "Define lines of every select page")
	flag.StringVarP(&print_dest, "pd", "d", "", "Define the path of the output destination")
	flag.BoolVarP(&page_type, "pt", "f", false, "Define the type of every select page")
	flag.BoolVarP(&help, "help", "h", false, "Show usage of selpg")
}

func main(){
	flag.Usage = func(){
		fmt.Fprintf(os.Stderr, "\nUsage: selpg [-s start_page] [-e end_page] [options] [in_filename] \nOptions: \n")
		flag.PrintDefaults()
	}
	flag.Parse()
	
	//handle main arguments
	if help{
		flag.Usage()
	}else if start_page == -1 {
		fmt.Fprintf(os.Stderr, "Missing selpg 1st argument: -s start_page\n")
		flag.Usage()
		os.Exit(1)
	}else if end_page == -1 {
		fmt.Fprintf(os.Stderr, "Missing selpg 2nd argument: -e end_page\n")
		flag.Usage()
		os.Exit(2)
	}else if start_page > end_page {
		fmt.Fprintf(os.Stderr, "Invalid pair of start_page and end_page\n")
		flag.Usage()
		os.Exit(3)
	}
	
	if page_length < 1 {
		fmt.Fprintf(os.Stderr, "Invalid page_length\n")
		flag.Usage()
		os.Exit(4)
	}
	
	if flag.NArg() > 0{
		in_filename = flag.Arg(0)
		file_info, err := os.Stat(in_filename)
		if err != nil{
			fmt.Fprintf(os.Stderr, "File not exist, check file path\n")
			flag.Usage()
			os.Exit(5)
		}else{
			file_mode := file_info.Mode()
	    		perm := file_mode.Perm()
			perm_flag := perm & os.FileMode(73) 
    			if uint32(perm_flag) == uint32(73) {
    				fmt.Fprintf(os.Stderr, "File does not allow you to read or write, check file path\n")
				flag.Usage()
				os.Exit(6)
			}	 	
		}
	}

	
	//process
	line_ptr := 0
	page_ptr := 1
	
	fin := os.Stdin
	response := ""	
 
	if in_filename != ""{
		err := errors.New("")
		fin , err = os.Open(in_filename)
		if err != nil{
			fmt.Fprintf(os.Stderr, "Open file fail, try again or check file path\n")
			flag.Usage()
			os.Exit(7)	
		}
		defer fin.Close()
	}
	
	read_line := bufio.NewReader(fin)
	if !page_type {
		for {
			line, err := read_line.ReadString('\n')
			if err == io.EOF{
				break
			}else if err != nil{
				fmt.Fprintf(os.Stderr, "Read file error, try again or check file path\n")
				flag.Usage()
				os.Exit(8)	
			}
			line_ptr++
			if line_ptr > page_length{
				page_ptr++
				line_ptr = 1
			}
			if page_ptr >= start_page && page_ptr <= end_page{
				response += line
			}
		}
	} else{
		for {
			page, err := read_line.ReadString('\f')
			if err == io.EOF{
				break
			}else if err != nil{
				fmt.Fprintf(os.Stderr, "Read file error, try again or check file path\n")
				flag.Usage()
				os.Exit(8)	
			}
			if page_ptr >= start_page && page_ptr <= end_page{
				response += page
			}
			
			page_ptr++
		}
	}

	//print
	cmd := exec.Command("cat", "-n")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create pipe error\n")
		flag.Usage()
		os.Exit(9)	
	}
	if print_dest != "" {
		cmd.Stdout = os.Stdout;
		cmd.Start()
		fmt.Fprintf(stdin, response)
		stdin.Close()
		cmd.Wait()
	} else{	
		fmt.Printf("Content of select pages :\n%s", response)
	}
}
