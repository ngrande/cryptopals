package main

import (
	"fmt"
)

func len_to_pad(size int, bs int) int {
	if (size == bs) {
		return bs * 2
	}

	return bs - size
}

func pad(in []byte, bs int) []byte {
	to_pad := len_to_pad(len(in), bs)

	var res []byte
	res = make([]byte, len(in) + to_pad)

	for b := range res {
		if (b < len(in)) {
			res[b] = in[b]
		} else {
			res[b] = byte(to_pad)
		}
	}

	return res
}

func main() {
	input := "YELLOW SUBMARINE" // 16 bytes
	block_size := 20 // don't know why 20...

	padded := pad([]byte(input), block_size)

	fmt.Println(padded)
	fmt.Printf("'%s'\n", string(padded))
}
