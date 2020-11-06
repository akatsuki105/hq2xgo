# hq2xgo

![Go](https://github.com/pokemium/hq2xgo/workflows/Go/badge.svg) [![GoDoc](https://godoc.org/github.com/pokemium/hq2xgo?status.svg)](https://godoc.org/github.com/pokemium/hq2xgo)

Enlarge image by 2x with hq2x algorithm

## Example(Before -> After)

<img src="./example/1/demo.png" width="320" height="288" />&nbsp;&nbsp;&nbsp;&nbsp;<img src="./example/1/demo_hq2x.png" />

<br />

<img src="./example/2/demo.png" width="320" height="288" />&nbsp;&nbsp;&nbsp;&nbsp;<img src="./example/2/demo_hq2x.png" />

## Usage

### command line

```sh
$ make build # require make and go
$ hq2x input.png output.png
```

### golang package

See [godoc](https://godoc.org/github.com/Akatsuki-py/hq2xgo) for details. 

```sh
$ go get github.com/pokemium/hq2xgo
```

```go

import (
	hq2x "github.com/pokemium/hq2xgo"
)

after, err := hq2x.HQ2x(before) // var before *image.RGBA

```