// input: 49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d
// convert from hex to base64
// expected: SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t

package main 

import "fmt"
import b64 "encoding/base64"
import hex "encoding/hex"

func main() {
	var input string = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	var expected string = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	var result string

	hex, _ := hex.DecodeString(input)
	hex_str := string(hex)

	fmt.Println(hex_str)
	
	sEnc := b64.StdEncoding.EncodeToString([]byte(hex_str))

//	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	result = string(sEnc)
	

	fmt.Println("Input: ", input)

	fmt.Println("Result: ", result)

	fmt.Println("Expected: ", expected)
	
	fmt.Println("Solved: ", result == expected)
}
