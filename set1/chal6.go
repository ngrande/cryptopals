package main

import "fmt"
import "math/bits"
import "os"
import "bufio"
import b64 "encoding/base64"
import "strings"

func calc_diff(left string, right string) int {
	diff := 0

	for i, _ := range left {
		diff += bits.OnesCount8(left[i] ^ right[i])
	}

	return diff
}

func detect_key_size(text string, min int, max int) int {

	//slice := text[:]
	best_score := 10000000
	best_size := 0
	for possible := min; possible <= max; possible++ {
		first_slice := text[:possible]
		second_slice := text[possible:possible * 2]
		if len(first_slice) != possible || len(second_slice) != possible {
			panic("Wrong code")
		}
		//slice = text[possible * 2:]

		edit_distance := calc_diff(first_slice, second_slice)
		normalized := edit_distance / possible
		fmt.Printf("d: %d, n: %d, key_size: %d\n", edit_distance, normalized, possible)
		if (normalized < best_score) {
			best_size = possible
			best_score = normalized
		}
	}

	return best_size
}

func is_valid_char(in byte) bool {
	lower := strings.ToLower(string(in))
	if strings.Contains("abcdefghijklmnopqrstuvwxyz\"\n&()!=?", lower) {
		return true
	}

	return false
}

func check_all_blocks(text string, key_size int, key_ind int, candidate byte) bool {
	// go through all blocks
	for i := 0; i < len(text); i+= key_size {
		block := text[i:i+key_size]
		encr := block[key_ind]
		decr := encr ^ candidate

		if is_valid_char(decr) == false {
			return false
		}
	}

	// is ok for all blocks!
	return true
}

func crypt(text string, key string) string {
	key_len := len(key)

	var result string

	for i := 0; i < len(text); i+= key_len {

		block := text[i:i + key_len]

		for k := 0; k < key_len; k++ {
			result += string(block[k] ^ key[k])
		}
	}

	return result
}

func score_text(text string, key string) int {
	decrypted := crypt(text, key)

	fmt.Println(decrypted)
	return 0
}

func search_key_candidates(text string, key_size int) []string {
	key_cand := make([]string, key_size)

	for i := 0; i < key_size; i++ {

		for char := 0; char < 256; char++ {
			possible := check_all_blocks(text, key_size, i, byte(char))

			if possible == false {
				continue
			} else {
				key_cand[i] += string(byte(char))
			}
		}
	}

	return key_cand
}

func crack_key(text string, key_size int) string {
	var key string
	//best_score := 0

	key_cand := search_key_candidates(text, key_size)
	fmt.Println("Candidates: ", key_cand)
	// key_cand looks like:
	// "HELLO"
	// "YOU"
	// "THERE"
	// -> possible keys: HYT, HYH, HYE, ... EYT, EYH, etc...

	// "H"
	// "IX"
	// "ABC"

	// 


	//var possible_key string

	possibilities := 1
	for i := 0; i < key_size; i++ {
		possibilities *= len(key_cand[i])
	}

	keys_arr := make([]string, possibilities)
	for i := 0; i < key_size; i++ {
		for char_ind := 0; char_ind < len(key_cand[i]); char_ind++ {
			keys_arr[i + char_ind] += string(key_cand[i][char_ind])
		}
	}


	fmt.Println("Possible keys: ", keys_arr)

	return key
}

func main() {
	test  := "this is a test"
	wokka := "wokka wokka!!!"
	expected_diff := 37

	diff := calc_diff(test, wokka)

	fmt.Println("Calced: ", diff)
	fmt.Println("Diff correct: ", diff == expected_diff)
	if diff != expected_diff {
		panic("False hamming distance")
	}


	input_test := "hello you"
	encrypted_test := crypt(input_test, "lol")
	decrypted_test := crypt(encrypted_test, "lol")

	if input_test != decrypted_test {
		panic("False crypt impl")
	}

	file, _ := os.Open("6.txt")

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var text string
	for scanner.Scan() {
		text += scanner.Text()
	}

	enc, _ := b64.StdEncoding.DecodeString(text)
	decoded := string(enc)
	//fmt.Println(decoded)

	//fmt.Println(text)

	key_size := detect_key_size(decoded, 2, 40)

	fmt.Println("Keysize detected: ", key_size)

	crack_key(decoded, key_size)


	crack_key(encrypted_test, 3)
}
