package main

import "fmt"
import hex "encoding/hex"

func rep_xor(key string, in string) string {
	key_ind := 0
	var res string
	for _, char := range in {
		res += string(byte(char) ^ key[key_ind])
		key_ind++
		if (key_ind >= len(key)) {
			key_ind = 0
		}
	}

	return hex.EncodeToString([]byte(res))
}

func main() {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"

	key := "ICE"

	expected := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"

	out := rep_xor(key, input)

	fmt.Println("Result: ", out)
	fmt.Println("Corret: ", out == expected)
}
