package main

import "fmt"
import "math/bits"
import "os"
import "bufio"
import b64 "encoding/base64"
import "strings"
import "regexp"

var key_gen_sep byte = 1
var common_english_words []string 
var english_word_regex *regexp.Regexp

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

	common_chars := "zqxjkvbpygfwmucldrhsnioate "

	for _, char := range in {
		
		if is_valid_char(byte(char)) == false {
			panic("Key should only produce valid characters!")
		}

		worth := strings.Index(common_chars, string(char))
		if worth >= 0 {
			sc += (worth + 1) * 12 // otherwise z would be worth 0
		} else {
			matched, _ := regexp.MatchString(`[A-Z0-9,.:!?]`, string(char))
			if matched {
				sc += 5
			}
		}
	}

	// do some more regex scoring
	matched := english_word_regex.FindAllString(in, -1)
	if len(matched) > 0	{
//		fmt.Println(matched)
//		fmt.Println(in)
		sc += 1000 * len(matched)
	}

	return sc
}

func score_text(text string, key string) (int, string) {
//	fmt.Println("Decrypting with key ", key, " text: ", text)
	decrypted := crypt(text, key)

//	fmt.Println("Decrypted")
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
					if len(string(char)) != 1 {
						panic("Char should be len 1!")
					}
					key_cand[i] += string(byte(char))
				}
			}
		}
	}

	return key_cand
}

func comb(in []string) string {

	var results strings.Builder
	for j := 0; j < len(in[0]); j++ {
		tmp := string(in[0][j]) 
			
		if len(in[1:]) > 0 {
			tmp += comb(in[1:])
		}

		results.WriteString(tmp + string(key_gen_sep))
	}

	return results.String()
}

func generate(data []string) []string {
	expected := 1
	for i := 0; i < len(data); i++ {
		expected *= len(data[i])
	}

	in := comb(data)

	split := strings.Split(in, string(key_gen_sep))

	fmt.Printf("Generating %d key combinations!\n", expected)
//	fmt.Println("Split: ", split)

	group_len := len(split[0])
	group := split[0]
	var result strings.Builder
	for i := 0; i < len(split); i++ {
//		if len(split[i]) == 0 || strings.Contains(split[i], string(key_gen_sep)) {
		split_len := len(split[i])
		if split_len == 0 {
			continue
		}
//		fmt.Println("Split: ", split[i])
		if split_len == group_len {
			group = split[i]
//			fmt.Println(group)
			result.WriteString(group + string(key_gen_sep))
			continue
		}
		
		group = group[0:group_len - split_len] + split[i] 
//		fmt.Println("Next: ", group)

		result.WriteString(group + string(key_gen_sep))
	}

	ret := strings.Split(result.String(), string(key_gen_sep))
	ret = ret[:len(ret) - 1]

	if len(ret) != expected {
		panic("Key generator broken - not expected number of keys!")
	}

	return ret
}

func crack_key(text string, key_size int) (int, string, string) {
	var key string

	key_cand := search_key_candidates(text, key_size)
	for cand := range key_cand {
		fmt.Printf("Candidates #%d: '%s'\n", cand, key_cand[cand])
	}

	// For debugging
	//key_cand = make([]string, 3)
	//key_cand[0] = "lo"
	//key_cand[1] = "ol"
	//key_cand[2] = "lo"

//	fmt.Println("Generating key combinations...")
	keys_arr := generate(key_cand)

	fmt.Printf("Possible keys (%d): %s\n", len(keys_arr), keys_arr)

	high_score := -100000
	high_score_text := ""
	for i := 0; i < len(keys_arr); i++ {
		pos_key := keys_arr[i]

		if len(pos_key) != key_size {
			panic("Size of possible key does not match key_size")
		}

//		fmt.Printf("Scoring with key: '%s'\n", pos_key)
		sc, decr := score_text(text, pos_key)
//		fmt.Printf("Text scored (%d) with key '%s': '%s'\n", sc, pos_key, decr)
		if sc > high_score {
			high_score = sc
			high_score_text = decr
			key = pos_key
			fmt.Printf("Key: '%s' Score: %d => %s\n", pos_key, high_score, high_score_text)
		}
	}

	return high_score, key, high_score_text
}

func load_english_dict() {

	fmt.Println("Loading most commong english words into RAM...")
	file, _ := os.Open("engl_words.txt")
	defer file.Close()

	var str strings.Builder
	str.WriteString(`(^|\s)(`)
	scanner := bufio.NewScanner(file)

	first := true
	for scanner.Scan() {
		common_english_words = append(common_english_words, []string{scanner.Text()}...)
		if first == false {
			str.WriteString("|")
		}
		first = false
		str.WriteString(scanner.Text())
	}

	str.WriteString(`)($|\s)`)
	fmt.Println("Regexp = ", str.String())
	fmt.Println(common_english_words)
	english_word_regex, _ = regexp.Compile(str.String())
}

func main() {

	load_english_dict()

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

	enc, _ := b64.StdEncoding.DecodeString(text)
	decoded := string(enc)

	// For debugging
	decoded = encrypted_test

	key_sizes := detect_key_size(decoded, 2, 40)
	for i := 0; i < len(key_sizes); i++ {
		fmt.Println("Checking for key size: ", key_sizes[i])
		key_score, key, decr := crack_key(decoded, key_sizes[i])

		fmt.Printf("Found key: '%s' with score: %d -> '%s'\n", key, key_score, decr)
	}
}
