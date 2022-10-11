package main

import (
	"flag"
	"fmt"
	"os"

	recoder "github.com/WelcomerTeam/Recoder"
)

func main() {
	InFile := flag.String("In", "", "Input filename")
	OutFile := flag.String("Out", "", "Output filename")
	Speed := flag.Int("Speed", 3, "Speed (1 slowest, 10 fastest)")

	flag.Parse()

	attrs := recoder.NewQuantizationAttributes()
	attrs.Speed = *Speed

	src, err := os.OpenFile(*InFile, os.O_RDONLY, 0444)
	if err != nil {
		fmt.Println(fmt.Errorf("os.OpenFile(%s): %w", *InFile, err))
		os.Exit(1)
	}

	defer src.Close()

	quant, err := recoder.RecodeImage(src, attrs)
	if err != nil {
		fmt.Println(fmt.Errorf("recoder.RecodeImage: %w", err))
		os.Exit(1)
	}

	dst, err := os.OpenFile(*OutFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("os.OpenFile(%s): %w", *OutFile, err))
		os.Exit(1)
	}

	defer dst.Close()

	_, err = dst.ReadFrom(quant)
	if err != nil {
		fmt.Println(fmt.Errorf("dst.ReadFrom(quant): %w", err))
		os.Exit(1)
	}

	os.Exit(0)
}
