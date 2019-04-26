package main

import "fmt"
import "os"
import "log"
import "bufio"
import hex "encoding/hex"

func score(in string) int {
	sc := 0

	common_chars := "zqxjkvbpygfwmucldrhsnioate"

	for _, char := range in {
		
		for worth, comm := range common_chars {
			if char == comm || byte(char) == byte(comm) - 22 {
				sc += worth
			}
		}
	}

	return sc
}

func xor(in string) (int, string) {
	high_score := 0
	var high_res string
	for char := 0; char < 256; char++ {
		res := ""
		for i, _ := range in {
			res += string(in[i] ^ byte(char))
		}
		
		sc := score(res)
		if sc > high_score {
			high_score = sc
			high_res = res
		}
	}

	fmt.Println("Best score (", high_score, ") => ", high_res)

	return high_score, high_res
}

func main() {
	file, err := os.Open("4.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	master_score := 0
	var master_res string
	line := -1
	for scanner.Scan() {
//		fmt.Println(scanner.Text())

		line++
		dec, _ := hex.DecodeString(scanner.Text())
		sc, res := xor(string(dec))

		if sc > master_score {
			master_score = sc
			master_res = res
		}

	}

	fmt.Println("Winner winner (", master_score, ") => ", master_res)
	fmt.Println("Line: ", line)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
