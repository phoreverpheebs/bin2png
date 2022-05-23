package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/schollz/progressbar"
)

var args struct {
	File   string `arg:"-f,required" help:"Directory or file of data to convert to image."`
	Output string `arg:"-o,required" help:"Output file of image data."`
	Invert bool   `arg:"-i" help:"Flips colours. (0 = White, 1 = Black)"`
}

func main() {
	printLogo()

	arg.MustParse(&args)

	if args.File == "" {
		fmt.Println("Make sure to supply a file argument ('-f file.txt')")
		os.Exit(1)
	}

	fileInfo, err := os.Stat(args.File)
	if err != nil {
		panic(err)
	}

	if args.Invert {
		px1 = 0x00
		px0 = 0xff
	} else {
		px1 = 0xff
		px0 = 0x00
	}

	var buf bytes.Buffer

	if fileInfo.IsDir() {
		err := filepath.Walk(args.File, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				f, err := os.ReadFile(path)
				if err != nil {
					fmt.Println(err)
				}

				buf.Write(f)
			}

			return nil
		})
		if err != nil {
			panic(err)
		}
		BinaryToPNG(buf.Bytes(), args.Output)
	} else {
		f, err := os.ReadFile(args.File)
		if err != nil {
			panic(err)
		}

		BinaryToPNG(f, args.Output)
	}
}

var (
	px1 byte
	px0 byte
)

func BinaryToPNG(b []byte, output string) {
	dimRaw := math.Sqrt(float64(len(b) * 8))

	dimEx := int(math.Round(dimRaw))

	img := image.NewGray(image.Rect(0, 0, dimEx, dimEx))

	fmt.Printf("\tWriting %.0f pixels...\n\n", math.Pow(float64(dimEx), 2))

	var x, y int

	pbpng := progressbar.New(len(b) * 8)
	pbpng.RenderBlank()
	for i := 0; i < len(b); i++ {
		bits := fmt.Sprintf("%08b", b[i])
		for n := 0; n < 8; n++ {
			pbpng.Add(1)
			// fmt.Print(fmt.Sprintf("\t\t\t\tx: %d \ty: %d \tbit: %s", x, y, string(bits[n])))

			if x >= dimEx {
				y++
				x = 0
			}

			switch bits[n] {
			case '1':
				img.SetGray(x, y, color.Gray{Y: px1})
			case '0':
				img.SetGray(x, y, color.Gray{Y: px0})
			}

			x++
		}
	}

	if !strings.HasSuffix(output, ".png") {
		args.Output += ".png"
	}

	imgFile, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}
	
	fmt.Print("\r\033[2K\r")

	bar := progressbar.DefaultBytes(-1)
	w := io.MultiWriter(imgFile, bar)
	err = png.Encode(w, img)
	if err != nil {
		panic(err)
	}
}

func printLogo() {
	fmt.Print(`

	 _     _       ____                    
	| |__ (_)_ __ |___ \ _ __  _ __   __ _ 
	| '_ \| | '_ \  __) | '_ \| '_ \ / _' |
	| |_) | | | | |/ __/| |_) | | | | (_| |
	|_.__/|_|_| |_|_____| .__/|_| |_|\__, |
                            |_|          |___/ 
   
   
   `)
}
