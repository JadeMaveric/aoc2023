package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func ReadInput(filename string) ([]string, error) {
	bytes, e := os.ReadFile(filename)
	if e != nil {
		return nil, e
	}

  chars := string(bytes)
	chars = strings.TrimSuffix(chars, "\n")

	content := strings.Split(chars, "\n")
	return content, nil
}

func SolveA(input []string) (string, error) {
	acc := 0

	for i, line := range input {
		calibration_val := -1
		last_digit := -1

		for _, ch := range line {
			if unicode.IsDigit(ch) {
				if calibration_val < 0 {
					calibration_val = int(ch - '0')
				}
				last_digit = int(ch - '0')
			}
		}

		if last_digit >= 0 {
			calibration_val = 10*calibration_val + last_digit
		}

		log.Printf("Line #%d: %d", i, calibration_val)

		if calibration_val >= 0 {
			acc += calibration_val
		}
	}

	return fmt.Sprint(acc), nil
}

func getNextDigit(line string, start int) (int, int) {
	digit_name := []string{"?", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for i, ch := range line[start:] {
		if unicode.IsDigit(ch) {
			return int(ch - '0'), i + 1
		}

		for d, name := range digit_name {
			if strings.HasPrefix(line[start+i:], name) {
				return d, i + len(name)
			}
		}
	}

	return 0, 0
}

func SolveB(input []string) (string, error) {
	acc := 0

	for _, line := range input {
		firstDigit := -1
		lastDigit := -1
		cursor := 0
		for {
			if cursor > len(line) {
				break
			}

			digit, _ := getNextDigit(line, cursor)

			// log.Print(line[cursor:], "\t", digit)
			if digit <= 0 {
				break
			}

			if firstDigit < 0 {
				firstDigit = digit
			}

			lastDigit = digit
			cursor += 1
		}

		num := 10*firstDigit + lastDigit

		if num > 0 {
			acc += num
		}
		log.Println(line, "\t", num)
	}
	return fmt.Sprint(acc), nil
}

func main() {
  log.SetOutput(os.Stdout)
  input, _ := ReadInput("in.txt")

	ans, _ := SolveB(input)

	log.Print(ans)
}
