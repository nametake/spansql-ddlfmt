package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	ddlfmt "github.com/nametake/spansql-ddlfmt"
)

func main() {
	write := flag.Bool("w", false, "write file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("no filename specified")
		flag.Usage()
		os.Exit(1)
	}

	filename := args[0]
	if err := run(filename, *write); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(filename string, write bool) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	ddl, err := ddlfmt.FormatDDL(filename, string(content))
	if err != nil {
		return fmt.Errorf("failed to format DDL: %v", err)
	}

	out, err := output(filename, write)
	if err != nil {
		return fmt.Errorf("failed to open output: %v", err)
	}
	defer out.Close()

	if _, err := out.Write([]byte(ddl)); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}

func output(filename string, write bool) (io.WriteCloser, error) {
	if !write {
		return os.Stdout, nil
	}
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	return file, nil
}
