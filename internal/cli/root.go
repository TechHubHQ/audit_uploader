package cli

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	InputFile  string
	Month      int
	MonthRange string
}

func DefineArgs(args *Args) {
	flag.StringVar(&args.InputFile, "input", "", "Path to the input file (required)")
	flag.IntVar(&args.Month, "month", 0, "Month for the report, e.g., '1' for January (optional, defaults to the latest sheet)")
	flag.StringVar(&args.MonthRange, "month-range", "", "Month range for the report, e.g., '1-3' for January-March (optional, defaults to the latest sheet)")
}

func Help() {
	output := flag.CommandLine.Output()
	fmt.Fprintf(output, "Usage: %s -input <file> [options]\n\n", os.Args[0])
	fmt.Fprintln(output, "Upload audit data from an input file.")
	fmt.Fprintln(output, "\nRequired:")
	fmt.Fprintln(output, "  -input <file>  Path to the input file.")
	fmt.Fprintln(output, "\nOptions:")
	flag.PrintDefaults()
}

func ReadArgs() Args {
	var args Args
	DefineArgs(&args)
	flag.Usage = Help
	flag.Parse()

	if args.InputFile == "" {
		Help()
		os.Exit(1)
	}

	return args
}
