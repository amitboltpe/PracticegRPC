package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {

	CompressPdf("./30MBFILE.pdf", "OutputFilePath.pdf")
}

func CompressPdf(inputFilePath, OutputFilePath string) {

	input := inputFilePath
	output := OutputFilePath
	cmd := exec.Command("gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/screen", // /screen, /ebook, /printer, /prepress
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile="+output,
		input,
	)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ghostscript compression failed: %v", err)
	}
	fmt.Println("PDF compressed:", output)
}
