package main

import (
	"bufio"
	"fmt"
	"io"
	"github.com/spf13/pflag"
	"math"
	"os"
	"os/exec"
)

type selpg_args struct {
	startPage int
	endPage int
	inFile string
	pageLen int
	pageType bool
	//ture for -f, false for -l
	printDestination string
}

var programName []byte
// parse for program name

func main() {
	args := new(selpg_args)
	ReceiveArgs(args)
	CheckArgs(args)
	HandleArgs(args)
}

func usage(){
	fmt.Fprintf(os.Stderr, "Usage error!\n")
	fmt.Fprintf(os.Stderr, "Usage: ")
	fmt.Fprintf(os.Stderr, "\tselpg --s=Number --e=Number [options] [filename]\n\n")
	fmt.Fprintf(os.Stderr, "\t--s=Number\tStart Page(Start <= End)\n")
	fmt.Fprintf(os.Stderr, "\t--e=Number\tEnd Page(Start <= End)\n")
	fmt.Fprintf(os.Stderr, "\t--l=Number\tLength of per Page, default 72\n")
	fmt.Fprintf(os.Stderr, "\t--f\t\tWhether using the form feeds\n")
	fmt.Fprintf(os.Stderr, "\t[filename]\tRead from file, default stadern input\n\n")
}

func ReceiveArgs(args *selpg_args) {
	pflag.Usage = usage;
	//add error informations
	pflag.IntVar(&(args.startPage), "s", -1, "start page")
	pflag.IntVar(&(args.endPage), "e", -1, "end page")
	pflag.IntVar(&(args.pageLen), "l", 72, "page len")
	pflag.StringVar(&(args.printDestination), "d", "", "print destionation")
	pflag.BoolVar(&(args.pageType), "f", false, "type of print")
	pflag.Parse()  
	//parse for input file names
	othersArg := pflag.Args()
	if len(othersArg) > 0 {
		args.inFile = othersArg[0]
	} else {
		args.inFile = ""
	}
}  

func CheckArgs(args *selpg_args) {
	if args.startPage == -1 || args.endPage == -1 {
		os.Stderr.Write([]byte("You should input --s --e at least\n"))
		pflag.Usage()
		os.Exit(0)
	}
	if args.startPage < 1 || args.startPage > (math.MaxInt32-1) {
		os.Stderr.Write([]byte("Invalid start page\n"))
		pflag.Usage()
		os.Exit(1)
	}
	if args.endPage < 1 || args.endPage > (math.MaxInt32-1) || args.endPage < args.startPage {
		os.Stderr.Write([]byte("Invalid end page\n"))
		pflag.Usage()
		os.Exit(2)
	}
	if args.pageLen < 1 || args.pageLen > (math.MaxInt32-1) {
		os.Stderr.Write([]byte("Invalid page length\n"))
		pflag.Usage()
		os.Exit(3)
	}
}

func HandleArgs(args *selpg_args) {
	var (
		reader  *bufio.Reader
		lineCtr int
		pageCtr int
	)

	//init reader
	if args.inFile == "" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		fileIn, err := os.Open(args.inFile)
		defer fileIn.Close()
		if err != nil {
			os.Stderr.Write([]byte("open file error\n"))
			os.Exit(4)
		}
		reader = bufio.NewReader(fileIn)
	}

	if args.printDestination == "" { 
	//output to os.stdout
		writer := bufio.NewWriter(os.Stdout)
		if args.pageType == true { 
			//-f type
			HandleArgs_f(reader, writer, args, &pageCtr) 
		} else { 
			//-l type
			HandleArgs_l(reader, writer, args, &pageCtr, &lineCtr)
		}
	} else { 
	//output to another command by pipe
		cmdGrep := exec.Command("./" + args.printDestination)
		stdinGrep, grepError := cmdGrep.StdinPipe()
		if grepError != nil {
			fmt.Println("Error happened about standard input pipe ", grepError)
			os.Exit(30)
		}
		writer := stdinGrep
		if grepError := cmdGrep.Start(); grepError != nil {
			fmt.Println("Error happened in execution ", grepError)
			os.Exit(30)
		}
		if args.pageType == true { 
			//-d type
			HandleArgs_f_d(reader, writer, args, &pageCtr)
		} else { 
			//-l type
			HandleArgs_l_d(reader, writer, args, &pageCtr, &lineCtr)
		}
		stdinGrep.Close()
		//make sure all the infor in the buffer could be read
		if err := cmdGrep.Wait(); err != nil {
			fmt.Println("Error happened in Wait process")
			os.Exit(30)
		}
	}
	if pageCtr < args.startPage {
		os.Stderr.Write([]byte("start page is greater than total page\n"))
		os.Exit(9)
	}
	if pageCtr < args.endPage {
		os.Stderr.Write([]byte("end page is greater than total page\n"))
		os.Exit(10)
	}
}

func HandleArgs_f(reader *bufio.Reader, writer *bufio.Writer, args *selpg_args, pageCtr *int) {
	*pageCtr = 1
	for {
		char, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			os.Stderr.Write([]byte("read byte from Reader fail\n"))
			os.Exit(7)
		}
		if *pageCtr >= args.startPage && *pageCtr <= args.endPage {
			errW := writer.WriteByte(char)
			if errW != nil {
				os.Stderr.Write([]byte("Write byte to out fail\n"))
				os.Exit(8)
			}
			writer.Flush()
		}
		if char == '\f' {
			(*pageCtr)++
		}
	}
}

func HandleArgs_l(reader *bufio.Reader, writer *bufio.Writer, args *selpg_args, pageCtr *int, lineCtr *int) {
	*lineCtr = 0
	*pageCtr = 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil { 
			if err == io.EOF {
				break
			}
			os.Stderr.Write([]byte("read bytes from Reader error\n"))
			os.Exit(5)
		}
		*lineCtr++
		if *pageCtr >= args.startPage && *pageCtr <= args.endPage {
			_, errW := writer.Write(line)
			if errW != nil {
				os.Stderr.Write([]byte("Write to file fail\n"))
				os.Exit(6)
			}
			writer.Flush()
		}
		if *lineCtr >= args.pageLen {
			*lineCtr = 0
			*pageCtr++
		}
	}
}

func HandleArgs_f_d(reader *bufio.Reader, writer io.WriteCloser, args *selpg_args, pageCtr *int) {
	*pageCtr = 1
	for {
		char, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			os.Stderr.Write([]byte("read byte from Reader fail\n"))
			os.Exit(7)
		}
		if *pageCtr >= args.startPage && *pageCtr <= args.endPage {
			writer.Write([]byte{char}) 
		}
		if char == '\f' {
			*pageCtr++
		}
	}
}

func HandleArgs_l_d(reader *bufio.Reader, writer io.WriteCloser, args *selpg_args, pageCtr *int, lineCtr *int) {
	*lineCtr = 0
	*pageCtr = 1
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil { 
			if err == io.EOF {
				break
			}
			os.Stderr.Write([]byte("read bytes from Reader error\n"))
			os.Exit(5)
		}
		*lineCtr++
		if *pageCtr >= args.startPage && *pageCtr <= args.endPage {
			_, errW := writer.Write(line)
			if errW != nil {
				os.Stderr.Write([]byte("Write to file fail\n"))
				os.Exit(6)
			}
		}
		if *lineCtr >= args.pageLen {
			*lineCtr = 0
			*pageCtr++
			writer.Write([]byte("\f"))
		}
	}
}
