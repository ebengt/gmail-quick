package main

import (
	"flag"
	"log"
	"os"
)

type config struct {
	from     string
	infile   string
	receiver string
	subject  string
}

func configuration(args []string, from string) (result *config) {
	result = new(config)
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.StringVar(&result.infile, "infile", "", "help message for infile")
	f.StringVar(&result.subject, "subject", "", "help message for subject")
	f.Parse(args[1:])
	if result.infile == "" {
		log.Fatalf("missing infile")
	}
	if result.subject == "" {
		log.Fatalf("missing subject")
	}
	result.receiver = f.Arg(0)
	if result.receiver == "" {
		log.Fatalf("missing receiver")
	}
	result.from = os.Getenv(from)
	if result.from == "" {
		log.Fatalf("missing from")
	}
	return result
}
