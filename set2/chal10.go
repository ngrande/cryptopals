package main

import (
	"fmt"
	"crypto/aes"
	"encoding/base64"
	"os"
	"bufio"
	"strings"
)

// Some links:
// why not to use ECB:
// https://crypto.stackexchange.com/questions/20941/why-shouldnt-i-use-ecb-encryption/20946#20946
// 
// How ECB works:
// https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_Codebook_.28ECB.29

// to be honest i copied this pasta from SO:
// Credits: https://stackoverflow.com/questions/24072026/golang-aes-ecb-encryption


func xor(left []byte, right []byte) []byte {
	res := make([]byte, len(left))

	for i, _ := range left {
		res[i] = left[i] ^ right[i]
	}

	return res
}

func cbc_encrypt(IV []byte, data []byte, key[]byte) []byte {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted := make([]byte, len(data))

	// CBC = Cipher Block Chaining.
	// each cipher block depends on the previous block.
	// each block (of plaintext) is XORed with the previous ciphertext block
	// before being encrypted.

	prev_block := make([]byte, len(key))
	if (len(key) != len(IV)) {
		panic("len key != len IV")
	}
	prev_block = IV
	bsize := len(key)

	for bs, be := 0, bsize; bs < len(data); bs, be = bs+bsize, be+bsize {
		xored := xor(prev_block, data[bs:be])
		cipher.Encrypt(encrypted[bs:be], xored)
		prev_block = encrypted[bs:be]
	}

	return encrypted
}

func cbc_decrypt(IV []byte, data []byte, key[]byte) []byte {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	var decrypted []byte

	// CBC = Cipher Block Chaining.
	// each cipher block depends on the previous block.
	// each block (of plaintext) is XORed with the previous ciphertext block
	// before being encrypted.

	prev_block := make([]byte, len(key))
	decrypted_block := make([]byte, len(key))
	if (len(key) != len(IV)) {
		panic("len key != len IV")
	}
	prev_block = IV
	bsize := len(key)

	for bs, be := 0, bsize; bs < len(data); bs, be = bs+bsize, be+bsize {
		cipher.Decrypt(decrypted_block, data[bs:be])
		xored := xor(prev_block, decrypted_block)
		decrypted = append(decrypted, xored...)
		prev_block = data[bs:be]
	}

	return decrypted
}

func main() {
	key := "YELLOW SUBMARINE" // 16 bytes
	IV := make([]byte, len(key))
	for i, _ := range key {
		IV[i] = 0
	}

	block_size := len(key)

	file, err := os.Open("./10.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var buf strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
//		fmt.Println(scanner.Text())
		buf.WriteString(scanner.Text())
	}

	decoded, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		panic(err)
	}

	res := cbc_decrypt(IV, []byte(decoded), []byte(key))

	fmt.Printf("start: %d\n", block_size)
	fmt.Printf("res: %s\n", string(res))
}
