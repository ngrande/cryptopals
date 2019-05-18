package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"encoding/base64"
	"crypto/aes"
)

// Some links:
// why not to use ECB:
// https://crypto.stackexchange.com/questions/20941/why-shouldnt-i-use-ecb-encryption/20946#20946
// 
// How ECB works:
// https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_Codebook_.28ECB.29

// to be honest i copied this pasta from SO:
// Credits: https://stackoverflow.com/questions/24072026/golang-aes-ecb-encryption
func decrypt_aes_128_ecb(data []byte, key []byte) []byte {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	
	decrypted := make([]byte, len(data))
	size := 16 // 16 bytes blocks


	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func main() {
	aes_128_key := "YELLOW SUBMARINE"
	fmt.Println("Key:", aes_128_key)

	file, err := os.Open("./7.txt")
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

	res := decrypt_aes_128_ecb(decoded, []byte(aes_128_key))

	fmt.Println(string(res))
	//fmt.Println(dec_str)
}
