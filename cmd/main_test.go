package main

import (
	"fmt"
	"image"
	"os"
	"testing"

	hq2x "github.com/pokemium/hq2xgo"
)

func getImage(path string) (image.Image, error) {
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

func BenchmarkRun(b *testing.B) {
	before, err := getImage("../example/1/demo.png")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hq2x.HQ2x(before.(*image.RGBA))
		if err != nil {
			panic(err)
		}
	}
}
