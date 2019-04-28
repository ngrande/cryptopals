package main

import "fmt"
import "math/bits"
import "os"
import "bufio"
import b64 "encoding/base64"
import "strings"

var key_gen_sep byte = 1

func calc_diff(left string, right string) int {
	diff := 0

	for i, _ := range left {
		diff += bits.OnesCount8(left[i] ^ right[i])
	}

	return diff
}

func detect_key_size(text string, min int, max int) []int {

	//slice := text[:]
	best_score := 10000000
	var res []int
	for possible := min; possible <= max; possible++ {
		if (possible * 2 > len(text)) {
			break
		}
		first_slice := text[:possible]
		second_slice := text[possible:possible * 2]
		if len(first_slice) != possible || len(second_slice) != possible {
			panic("Wrong code")
		}
		//slice = text[possible * 2:]

		edit_distance := calc_diff(first_slice, second_slice)
		normalized := edit_distance / possible
		fmt.Printf("d: %d, n: %d, key_size: %d\n", edit_distance, normalized, possible)
		if (normalized <= best_score) {
			best_score = normalized
			res = append(res, []int{ possible }...)
		}
	}

	return res 
}

func is_valid_char(in byte) bool {
	// all printable ASCII chars
	if in >= 32 && in <= 126 {
		if len(string(in)) == 0 {
			panic("Not a valid char!")
		}
		return true
	}
	return false

//	lower := strings.ToLower(string(in))
//	if strings.Contains("abcdefghijklmnopqrstuvwxyz\"\n&()!=? ,", lower) {
//		return true
//	}
//
//	return false
}

func check_all_blocks(text string, key_size int, key_ind int, candidate byte) bool {
	// go through all blocks
	for i := 0; i < len(text); i+= key_size {
		end := i + key_size
		if end > len(text) {
			end = len(text)
		}
		block := text[i:end]

		if key_ind + i < len(text) {
			encr := block[key_ind]
			decr := encr ^ candidate

			if is_valid_char(decr) == false {
				return false
			}
		}
	}

	// is ok for all blocks!
	return true
}

func crypt(text string, key string) string {
	key_len := len(key)
	if key_len == 0 {
		panic("Key len can not be 0!")
	}

	var result string

	for i := 0; i < len(text); i+= key_len {
		end := i + key_len
		if end > len(text) {
			end = len(text)
		}
		block := text[i:end]

		for k := 0; k < key_len; k++ {
			if (k + i >= len(text)) {
				// text block < key_len
				break
			}
			result += string(block[k] ^ key[k])
		}
	}

	return result
}

func score(in string) int {
	sc := 0

	common_chars := "!,zqxjkvbpygfwmucldrhsnioate "

	for _, char := range in {
		
		if is_valid_char(byte(char)) == false {
			fmt.Println("Wrong score: ", in)
			return -1
		}

		for worth, comm := range common_chars {
			if char == comm || byte(char) == byte(comm) - 22 {
				sc += worth
			}
		}
	}

	return sc
}

func score_text(text string, key string) (int, string) {
	fmt.Println("Decrypting with key ", key, " text: ", text)
	decrypted := crypt(text, key)

	fmt.Println("Decrypted")
	sc := score(decrypted)
	return sc, decrypted
}

func search_key_candidates(text string, key_size int) []string {
	key_cand := make([]string, key_size)

	for i := 0; i < key_size; i++ {

		for char := 0; char < 256; char++ {
			possible := check_all_blocks(text, key_size, i, byte(char))

			if possible == false {
				continue
			} else {
				// only add if is a valid char
				// we expect the key to be "readable"
				if is_valid_char(byte(char)) {
					key_cand[i] += string(byte(char))
				}
			}
		}
	}

	return key_cand
}

func comb(in []string) string {

	results := ""
	for j := 0; j < len(in[0]); j++ {
		tmp := string(in[0][j]) 
			
		if len(in[1:]) > 0 {
			tmp += comb(in[1:])
		}

		results += tmp + string(key_gen_sep)
	}

	return results
}

func generate(data []string) []string {
	in := comb(data)

	split := strings.Split(in, string(key_gen_sep))

//	fmt.Println(split)

	//is_group := false
	group_len := len(split[0])
//	fmt.Println(len(split[0]))
	group := split[0]
	result := ""
	for i := 0; i < len(split); i++ {
		if strings.Contains(split[i], string(key_gen_sep)) || len(split[i]) == 0 {
			continue
		}
//		fmt.Println("Split: ", split[i])
		if len(split[i]) == group_len {
			group = split[i]
//			fmt.Println(group)
			result += group + " "
			continue
		}
		
		ngroup_len := len(split[i])
		group = group[0:group_len - ngroup_len] + split[i] 
//		fmt.Println("Next: ", group)

		result += group + " "

	}

	result = result[:len(result) - 1]

	return strings.Split(result, " ")
}

func crack_key(text string, key_size int) (int, string) {
	var key string

	key_cand := search_key_candidates(text, key_size)
	for cand := range key_cand {
		fmt.Printf("Candidates #%d: '%s'\n", cand, key_cand[cand])
	}
	// key_cand looks like:
	// "HELLO"
	// "YOU"
	// "THERE"
	// -> possible keys: HYT, HYH, HYE, ... EYT, EYH, etc...

	// "H"
	// "IX"
	// "ABC"

	// 
	//key_cand = make([]string, 3)
	//key_cand[0] = "ABl"
	//key_cand[1] = "oX"
	//key_cand[2] = "ol"

	keys_arr := generate(key_cand)

	fmt.Printf("Possible keys (%d): %s\n", len(keys_arr), keys_arr)

	high_score := -100000
	high_score_text := ""
	for i := 0; i < len(keys_arr); i++ {
		pos_key := keys_arr[i]

		fmt.Println("Scoring...")
		sc, decr := score_text(text, pos_key)
		fmt.Println("Text scored: ", decr)
		if sc > high_score {
			high_score = sc
			high_score_text = decr
			key = pos_key
			fmt.Printf("Key: %s Score: %d => %s\n", pos_key, high_score, high_score_text)
		}
	}

	return high_score, key
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


	input_test := "hello, my dear friend!"
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

	// ACTIVATE LATER
	enc, _ := b64.StdEncoding.DecodeString(text)
	decoded := string(enc)

	if decoded == text {
		fmt.Println("LOL")
	}

//	key_size := detect_key_size(decoded, 2, 40)

//	fmt.Println("Keysize detected: ", key_size)

//	crack_key(decoded, key_size)

	key_sizes := detect_key_size(encrypted_test, 2, 40)
	for i := 0; i < len(key_sizes); i++ {
		fmt.Println("Checking for key size: ", key_sizes[i])
		key_score, key := crack_key(encrypted_test, key_sizes[i])

		fmt.Printf("Found key: '%s' with score: %d\n", key, key_score)
	}
}
