// in1: 1c0111001f010100061a024b53535009181c
// in2: 686974207468652062756c6c277320657965
// expected: 746865206b696420646f6e277420706c6179

package main

import "fmt"
import hex "encoding/hex"

func main() {
	in1 := "1c0111001f010100061a024b53535009181c"
	in2 := "686974207468652062756c6c277320657965"
	expected := "746865206b696420646f6e277420706c6179"

	buf1, _ := hex.DecodeString(in1)
	buf2, _ := hex.DecodeString(in2)

	fmt.Println(buf1)
	fmt.Println(buf2)

	xor := buf1
	for i, _ := range xor {
		fmt.Println("XOR: ", buf1[i], " ^ ", buf2[i])
		xor[i] = buf1[i] ^ buf2[i]
	}

	fmt.Println(string(xor))
	fmt.Println(xor)

	hex_str := hex.EncodeToString(xor)
	fmt.Println(hex_str)

	fmt.Println("Solved: ", hex_str == expected)
}
