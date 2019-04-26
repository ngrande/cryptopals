package main

import "fmt"
import hex "encoding/hex"
//import "reflect"

func score(in string) int {
	high_score := 0
	// most common english characters
	most_common_letters := "aAeEiIoOuU"
	for _, char := range in {
		for _, common := range most_common_letters {
			if char == common {
				high_score += 1
			}
		}
	}
	return high_score
}

func main() {
	input := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	dec, _ := hex.DecodeString(input)
	
//	fmt.Println(reflect.TypeOf(dec))

	var res string
	res_score := 0
	res_key := 0
	for key := 0; key < 256; key++ {
		tmp := dec
		for i, _ := range tmp {
			tmp[i] = tmp[i] ^ byte(key)
		}

		sc := score(string(tmp))

		if sc > res_score {
			res = string(tmp)
			res_score = sc
			res_key = key
			fmt.Println("new high score: ", res_score, " for '", res, "' with key: ", res_key)
		}
	}

	fmt.Println("Result: ", res)
	fmt.Println("Key: ", res_key)

}
