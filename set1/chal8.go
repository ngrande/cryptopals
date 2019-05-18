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
// applying a key that only differs int he last byte
// shall produce the same result except for the last byte
// but this only applies for ECB -> this way we detect it!

func detect_ecb(line []byte) bool {
	KEY_1 := "YELLOW SUBMARINE"
	KEY_2 := "YELLOW SUBMARINO"

	res_1 := decrypt_block(line, []byte(KEY_1))
	res_2 := decrypt_block(line, []byte(KEY_2))

	if string(res_1)[:len(res_1)-1] == string(res_2)[:len(res_2)-1] {
		return true
	}

	return false
}

func main() {
	
	file, err := os.Open("./8.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
//	var buf strings.Builder

	for scanner.Scan() {
//		buf.WriteString(scanner.Text())
		decoded, err := hex.DecodeString(scanner.Text())
		if err != nil {
			panic(err)
		}

		is_ecb := detect_ecb(decoded)

		if is_ecb {
			fmt.Println("Result:", string(decoded))
		}
	}
}
