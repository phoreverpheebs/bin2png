package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/schollz/progressbar"
)

var formats []string
var outputs = make(map[string]string)

/**
 *	CLI Arguments <3
**/
var args struct {
	File   string `arg:"-f" help:"Directory or file of data to convert to image."`
	Output string `arg:"-o" help:"Output file of image data."`

	PNG  bool `help:"PNG Output."`
	JPEG bool `help:"JPEG Output."`

	Invert bool `arg:"-i" help:"Flips colours. \n\t\t\t\t(white represents a 0; black represents a 1)"`

	Compression int `arg:"-c" help:"PNG Level of compression. \n\t\t\t\t(0 = Default, 1 = None, 2 = Fastest, 3 = Best)" default:"0" placeholder:"LEVEL"`
	Quality     int `arg:"-q" help:"JPEG Quality. \n\t\t\t\t(1 - 100; 1 is lowest, 100 is highest)" default:"100" placeholder:"QUALITY"`
}

func main() {
	printLogo()

	arg.MustParse(&args)

	if args.File == "" {
		fmt.Println("Make sure to supply a file argument ('-f file.txt')")
		os.Exit(1)
	}

	if !args.PNG && !args.JPEG {
		args.PNG = true
	}

	var compLevel int

	_, filename := filepath.Split(args.File)
	if strings.HasPrefix(filename, ".") {
		filename = "output-" + filename
	}

	if args.Invert {
		px1 = 0x00
		px0 = 0xff
	} else {
		px1 = 0xff
		px0 = 0x00
	}

	if args.PNG {
		formats = append(formats, "png")

		switch args.Compression {
		case 0:
			compLevel = 0
		case 1:
			compLevel = -1
		case 2:
			compLevel = -2
		case 3:
			compLevel = -3
		default:
			panic("Unknown compression level, must be in range 0 - 3")
		}
	}
	if args.JPEG {
		formats = append(formats, "jpg")
	}

	if len(formats) > 1 {
		for _, x := range formats {
			if args.Output == "" {
				outputs[x] = filename + "." + x
			} else {
				outputs[x] = args.Output + "." + x
			}
		}
	} else {
		if args.Output != "" && !strings.HasSuffix(args.Output, formats[0]) {
			outputs[formats[0]] = args.Output + "." + formats[0]
		} else if args.Output == "" {
			outputs[formats[0]] = filename + "." + formats[0]
		} else {
			outputs[formats[0]] = args.Output
		}
	}

	fileInfo, err := os.Stat(args.File)
	if err != nil {
		panic(err)
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

		BinaryToPNG(buf.Bytes(), outputs, compLevel)
	} else {
		f, err := os.ReadFile(args.File)
		if err != nil {
			panic(err)
		}

		BinaryToPNG(f, outputs, compLevel)
	}
}

var (
	px1 byte
	px0 byte
)

func BinaryToPNG(b []byte, outputTo map[string]string, compression int) {

	dimRaw := math.Sqrt(float64(len(b) * 8))

	dimEx := int(math.Round(dimRaw))

	img := image.NewGray(image.Rect(0, 0, dimEx, dimEx))

	fmt.Printf("\tWriting %.0f pixels...\n", math.Pow(float64(dimEx), 2))

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

	for key, obj := range outputTo {
		fmt.Printf("\r\033[2K\r\tEncoding %s...\n", key)
		bar := progressbar.DefaultBytes(-1)

		imgFile, err := os.OpenFile(obj, os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			panic(err)
		}
		w := io.MultiWriter(imgFile, bar)

		switch key {
		case "png":
			enc := &png.Encoder{
				CompressionLevel: png.CompressionLevel(compression),
			}

			err = enc.Encode(w, img)
			if err != nil {
				panic(err)
			}
		case "jpg":
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: args.Quality})
			if err != nil {
				panic(err)
			}
		}

		fmt.Print("\r\033[2K\r\tEncoded!\n")
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
