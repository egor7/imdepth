package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"image/color"
	_ "image/gif"  // register the GIF format with the image package
	_ "image/jpeg" // register the JPEG format with the image package
	"image/png"
)

var area = flag.Int("r", 2, "area around every point to get its sharp value")
var regFile = regexp.MustCompile(`^[0-9]+\..+$`)
var regLvl = regexp.MustCompile(`^[0-9]+`)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n")
	fmt.Fprintf(os.Stderr, "  imdepth [flags] <dir_name>\n")
	fmt.Fprintf(os.Stderr, "rules:\n")
	fmt.Fprintf(os.Stderr, "  - <dir_name> must constain a files named <number>.<ext>\n")
	fmt.Fprintf(os.Stderr, "  - this <number>s used as heights: [0..255]\n")
	fmt.Fprintf(os.Stderr, "  - <ext> should be jpg or png\n")
	fmt.Fprintf(os.Stderr, "flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if flag.NArg() != 1 {
		usage()
	}

	var levels, width, height int
	imgs := make(map[int][]uint8)

	// loop throug dir
	files, err := ioutil.ReadDir(args[0])
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !regFile.MatchString(f.Name()) {
			continue
		}

		// open image
		levels += 1
		level, _ := strconv.Atoi(regLvl.FindString(f.Name()))
		imf, err := os.Open(filepath.Join(args[0], f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer imf.Close()
		fmt.Printf("read %s\n", imf.Name())

		// read image size
		im, _, err := image.Decode(imf)
		if err != nil {
			log.Fatal(err)
		}
		bounds := im.Bounds()
		w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
		if width != 0 && width != w {
			log.Fatalf("width not matched: got %d, expect %d", w, width)
		}
		if height != 0 && height != h {
			log.Fatalf("height not matched: got %d, expect %d", h, height)
		}
		width, height = w, h

		// read image data
		imgs[level] = make([]uint8, w*h)
		for c := 0; c < w; c++ {
			for r := 0; r < h; r++ {
				imgs[level][c*h+r] = color.GrayModel.Convert(im.At(c, r)).(color.Gray).Y
			}
		}

	}

	// process
	fmt.Printf("processing %d images [%dx%d]: ", levels, width, height)
	im_heights := image.NewGray(image.Rect(0, 0, width, height))
	im_mesh := image.NewGray(image.Rect(0, 0, width, height))
	for c := 0; c < width; c++ {
		for r := 0; r < height; r++ {
			var max float64
			var max_l int
			for l, _ := range imgs {
				d := disp(imgs[l], c, r, width, height)
				if d >= max {
					max = d
					max_l = l
				}
			}
			im_heights.Set(c, r, color.Gray{uint8(max_l)})
			im_mesh.Set(c, r, color.Gray{imgs[max_l][c*height+r]})
		}
		if c%10 == 0 {
			fmt.Printf(".")
		}
	}
	fmt.Printf(" DONE\n")

	// save heights
	heights, err := os.Create(filepath.Join(args[0], "heights.png"))
	defer heights.Close()
	png.Encode(heights, im_heights)
	fmt.Printf("save %s\n", heights.Name())

	// save mesh
	mesh, err := os.Create(filepath.Join(args[0], "mesh.png"))
	defer mesh.Close()
	png.Encode(mesh, im_mesh)
	fmt.Printf("save %s\n", mesh.Name())
}

func disp(img []uint8, col, row, width, height int) float64 {
	var (
		n     int
		x, x2 float64
	)
	for c := col - *area; c < col+*area; c++ {
		for r := row - *area; r < row+*area; r++ {
			if c < 0 || c >= width || r < 0 || r >= height || (col-c)*(col-c)+(row-r)*(row-r) > *area**area {
				// continue
			} else {
				n += 1
				x += float64(img[c*height+r])
				x2 += float64(img[c*height+r]) * float64(img[c*height+r])
			}
		}
	}
	return x2/float64(n) - (x/float64(n))*(x/float64(n))
}
