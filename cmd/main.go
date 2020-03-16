package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"

	hq2x "github.com/Akatsuki-py/hq2xgo"

	_ "image/jpeg"
	"image/png"
	_ "image/png"
)

func main() {
	os.Exit(Run())
}

// Run - run app
func Run() int {
	flag.Parse()

	input := flag.Arg(0)
	if input == "" {
		help()
		return 0
	}

	output := flag.Arg(1)
	if output == "" {
		base, ext := getFileNameWithoutExt(input)
		output = base + "_hq2x" + ext
	}

	before, err := openImage(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	after, err := hq2x.HQ2x(before.(*image.RGBA))

	if err := saveImage(output, after); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func help() {
	fmt.Println("hq2x <input> [<output>]")
}

func getFileNameWithoutExt(path string) (filename string, ext string) {
	filename = path[:len(path)-len(filepath.Ext(path))]
	ext = filepath.Ext(path)
	return filename, ext
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return img, err
	}

	return img, nil
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		return err
	}
	return nil
}
