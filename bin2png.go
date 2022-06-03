package main

import (
	"bin2png/encode/png"
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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
	Output string `arg:"-o" help:"Output file of image data.\n" placeholder:"OUT"`

	Verbose bool `arg:"-v" help:"Verbose printing.\n"`

	PNG  bool `help:"PNG Output."`
	JPEG bool `help:"JPEG Output.\n"`

	Invert bool `arg:"-i" help:"Flips colours. \n\t\t\t\t(white represents a 0; black represents a 1)\n"`

	Compression int `arg:"-c" help:"PNG Level of compression. \n\t\t\t\t(0 = Default, 1 = None, 2 = Fastest, 3 = Best)" default:"0" placeholder:"LEVEL"`
	Quality     int `arg:"-q" help:"JPEG Quality. \n\t\t\t\t(1 - 100; 1 is lowest, 100 is highest)" default:"100" placeholder:"QUALITY"`
}

func main() {
	printLogo()

	arg.MustParse(&args)

	verbose(fmt.Sprintf(`Arguments:

	File:				%s
	Output:				%s

	Verbose:			%t

	PNG:				%t
	JPEG:				%t

	Invert:				%t

	Compression:			%d
	Quality:			%d
		
`, args.File, args.Output, args.Verbose, args.PNG, args.JPEG, args.Invert, args.Compression, args.Quality))

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
		verbose("Adding PNG to formats array.\n")
		formats = append(formats, "png")

		switch args.Compression {
		case 0:
			verbose("Compression setting set to default compression.\n")
			compLevel = 0
		case 1:
			verbose("Compression setting set to no compression.\n")
			compLevel = -1
		case 2:
			verbose("Compression setting set to fastest compression.\n")
			compLevel = -2
		case 3:
			verbose("Compression setting set to best compression.\n")
			compLevel = -3
		default:
			panic("Unknown compression level, must be in range 0 - 3")
		}
	}

	if args.JPEG {
		verbose("Adding JPEG to formats array.\n")
		formats = append(formats, "jpg")
	}

	verbose("Formatting outputs.\n")
	if len(formats) > 1 {
		for _, x := range formats {
			if args.Output == "" {
				outputs[x] = filepath.Clean(filename + "." + x)
			} else {
				outputs[x] = args.Output + "." + x
			}
		}
	} else {
		if args.Output != "" && !strings.HasSuffix(args.Output, formats[0]) {
			outputs[formats[0]] = args.Output + "." + formats[0]
		} else if args.Output == "" {
			outputs[formats[0]] = filepath.Clean(filename + "." + formats[0])
		} else {
			outputs[formats[0]] = args.Output
		}
	}

	fileInfo, err := os.Stat(args.File)
	if err != nil {
		panic(err)
	}

	var buf = new(bytes.Buffer)

	if fileInfo.IsDir() {
		verbose("Recursively parsing input directory.\n")
		err := filepath.Walk(args.File, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && args.Verbose {
				verbose("Current directory: " + info.Name() + "\n")
			} else if !info.IsDir() {
				verbose("Reading file: " + info.Name() + "\n")
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

		BinaryToPNG(buf, outputs, compLevel)
	} else {
		BinaryToPNG(readFile(args.File), outputs, compLevel)
	}
}

var (
	px1 byte
	px0 byte
)

func newMultiWriter(file string) io.Writer {
	bar := progressbar.DefaultBytes(-1)
	imgFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}
	return io.MultiWriter(imgFile, bar)
}

func BinaryToPNG(b *bytes.Buffer, outputTo map[string]string, compression int) {

	verbose("Calculating 1:1 image resolution.\n\n")
	dimRaw := math.Sqrt(float64(b.Len() * 8))

	dimEx := int(math.Round(dimRaw))

	var bufbytes []byte

	if dimEx >= 1<<16 && args.JPEG {
		fmt.Println("\tImage is too large for JPEG. Skipping JPEG encoding.")
		delete(outputTo, "jpg")
		fmt.Println()
	} else {
		bufbytes = b.Bytes()
	}

	if ol := len(outputTo); ol == 0 {
		os.Exit(1)
	}

	rect := image.Rect(0, 0, dimEx, dimEx)

	fmt.Printf("\tWriting pixels of %dx%d image.", dimEx, dimEx)
	if args.Verbose {
		fmt.Printf(" (%.0f pixels)\n", math.Pow(float64(dimEx), 2))
	}

	fmt.Println()

	for key, obj := range outputTo {
		verbose("Format: " + key + "\n")
		verbose("Output: " + obj + "\n")

		if obj == "."+key {
			obj = "output." + key
		}

		switch key {
		case "png":
			doPNG(b, obj, rect, compression)
		case "jpg":
			doJPEG(bufbytes, obj, rect, dimEx)
		}

		fmt.Printf("\r\x1b[2K\r\tEncoded to %s!\n\n", obj)
	}
}

func readFile(path string) (buffer *bytes.Buffer) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer = bytes.NewBuffer(make([]byte, 0))
	p := make([]byte, 1024)

	var count int

	for {
		if count, err = reader.Read(p); err != nil {
			break
		}
		buffer.Write(p[:count])
	}
	if err != io.EOF {
		panic("Error reading " + path + ": " + err.Error())
	}

	return
}

func doPNG(b *bytes.Buffer, file string, rect image.Rectangle, compression int) {
	var buf = new(bytes.Buffer)
	verbose("Writing data to buffer.\n\n")
	pbpng := progressbar.New(b.Len() * 8)
	pbpng.RenderBlank()

	for {
		curb, err := b.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		for _, n := range fmt.Sprintf("%08b", curb) {
			pbpng.Add(1)
			switch n {
			case '1':
				err = buf.WriteByte(px1)
			case '0':
				err = buf.WriteByte(px0)
			}
			if err != nil {
				panic(err)
			}
		}
	}

	m := newMultiWriter(file)

	fmt.Println("Done!")
	fmt.Println("\tEncoding png...")

	enc := &png.Encoder{
		CompressionLevel: png.CompressionLevel(compression),
	}
	enc.Encode(m, buf, rect)
}

func doJPEG(b []byte, file string, rect image.Rectangle, dime int) {
	img := image.NewGray(rect)

	var x, y int

	verbose("Writing data to image.Image\n\n")

	pbpng := progressbar.New(len(b) * 8)
	pbpng.RenderBlank()

	for _, by := range b {
		for _, n := range fmt.Sprintf("%08b", by) {
			pbpng.Add(1)
			// fmt.Print(fmt.Sprintf("\t\t\t\tx: %d \ty: %d \tbit: %s", x, y, string(bits[n])))

			if x >= dime {
				y++
				x = 0
			}

			switch n {
			case '1':
				img.SetGray(x, y, color.Gray{Y: px1})
			case '0':
				img.SetGray(x, y, color.Gray{Y: px0})
			}

			x++
		}
	}

	w := newMultiWriter(file)

	fmt.Println("Done!")
	fmt.Println("\tEncoding jpg...")

	err := jpeg.Encode(w, img, &jpeg.Options{Quality: args.Quality})
	if err != nil {
		panic(err)
	}
}

func verbose(str string) {
	if args.Verbose {
		fmt.Printf("[V] %s", str)
	}
	return
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
