package main

import (
	"fmt"
	"encoding/hex"
	"os"
	"bufio"
	//"strings"
	"crypto/aes"
)

func decrypt_block(line, key []byte) []byte {
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	dest := make([]byte, len(line))
	cipher.Decrypt(dest, line)

	return dest
}

// strategy:
// without knowing the key:
// the input is 160 characters
// we will just use a random 16 byte long key
// and apply it to each block (of the same size)
// to the whole line (160 byte) and look for the
// line with the most repetitions!

// TODO
// next step:
// crack the key:
// * try to detect the key len
// * find all key combinations which
// produce only ASCII chars as output

func detect_ecb(line []byte) int {
	KEY := "YELLOW SUBMARINE"

	// block size = bs
	bs := len(KEY)

	// map with counter of occurences per block decipher
	var occur map[string]int
	occur = make(map[string]int)
	for begin, end := 0, bs; begin < len(line); begin, end = begin+bs, end+bs {
		res := decrypt_block(line[begin:end], []byte(KEY))
		k := string(res)
		count, _ := occur[k]
		count += 1
		occur[k] = count
	}

	score := 0
	for _, v := range(occur) {
		if v > score {
			score = v
		}
	}

	return score
}

func main() {
	
	file, err := os.Open("./8.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
//	var buf strings.Builder

	high_score := 0
	for scanner.Scan() {
//		buf.WriteString(scanner.Text())
		decoded, err := hex.DecodeString(scanner.Text())
		if err != nil {
			panic(err)
		}

		score := detect_ecb(decoded)

		if score > high_score {
			fmt.Println("Score:", score)
			fmt.Println("Result:", string(decoded))
			high_score = score
		}
	}
}
