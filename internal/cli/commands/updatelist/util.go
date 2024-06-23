package updatelist

import (
	"log"
	"os"
)

func openfile(infile string) *os.File {
	if infile == "" {
		log.Println("Reading from stdin.")
		return os.Stdin
	}
	in, err := os.Open(infile)
	if err != nil {
		log.Fatalf("Error opening input file %q: %s", infile, err)
	}
	return in
}

func createfile(outfile string) *os.File {
	if outfile == "" {
		return os.Stdout
	}
	out, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("Error opening output file %q: %s", outfile, err)
	}
	return out
}

func closefiles(in *os.File, out *os.File) {
	if err := in.Close(); err != nil {
		log.Printf("Warning: an error occurred during close of input file: %s", err)
	}
	if err := out.Close(); err != nil {
		log.Printf("Warning: an error occurred during close of output file: %s", err)
	}
}
