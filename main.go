package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type GuessResult struct {
	IsMatch bool
	Target string
	Guess string
	GuessRune []rune
	TargetRune []rune
	Matches []bool
	MisMatches []bool
}

func guess(guess string, target string) (result GuessResult, err error) {
	result = GuessResult{}
	result.Target = target
	result.Guess = guess
	result.GuessRune = []rune(guess)
	result.TargetRune = []rune(target)
	result.Matches = make([]bool, len(guess))
	result.MisMatches = make([]bool, len(guess))

	result.IsMatch = guess == target
	if !result.IsMatch {
		for i := 0; i < len(result.GuessRune); i++ {
			result.Matches[i] = result.GuessRune[i] == result.TargetRune[i]
			result.MisMatches[i] = false
			if !result.Matches[i] {
				for j := 0; j < len(result.TargetRune); j++ {
					if result.GuessRune[i] == result.TargetRune[j] {
						result.MisMatches[i] = true
						break
					}
				}
			}
		}
	}
	return result, nil

}

func findMatches(result GuessResult, answers []string) []string {
	var matchedWords []string
	if result.IsMatch {
		matchedWords = append(matchedWords, result.Target)
	} else {
		matches := "^"
		var mismatchRegexList []string
		excluded := ""
		var mustContainLetters []rune
		hasMatches := false
		hasMismatches := false
		for i := 0; i < 5; i++ {
			if result.Matches[i] {
				matches += string(result.GuessRune[i])
				mustContainLetters = append(mustContainLetters, result.GuessRune[i])
				hasMatches = true
			} else {
				matches += "."
				if result.MisMatches[i] {
					mismatchRegex := "^"
					for j := 0; j < 5; j++ {
						if j == i {
							mismatchRegex += string(result.GuessRune[i])
						} else {
							mismatchRegex += "."
						}
					}
					mismatchRegexList = append(mismatchRegexList, mismatchRegex)
					mustContainLetters = append(mustContainLetters, result.GuessRune[i])
					hasMismatches = true
				} else {
					excluded += string(result.GuessRune[i])
				}
			}
		}
		var testWords []string

		// word has these letters but maybe in the wrong place,
		// so exclude words that don't have all these letters
		if len(mustContainLetters) > 0 {
			for _, word := range answers {
				isExcluded := false
				for _, include := range mustContainLetters {
					if !strings.Contains(word, string(include)) {
						isExcluded = true
						break
					}
				}
				if !isExcluded {
					testWords = append(testWords, word)
				}
			}
		} else {
			testWords = answers
		}
		// word doesn't have ANY of these letters
		if len(excluded) > 0 {
			var newTestWords []string
			for _, word := range testWords {
				isExcluded := false
				for _, exclude := range excluded {
					if strings.Contains(word, string(exclude)) {
						isExcluded = true
						break
					}
				}
				if !isExcluded {
					newTestWords = append(newTestWords, word)
				}
			}
			testWords = newTestWords
		}

		// exclude words that have the matched letters in the wrong places
		if hasMismatches {
			var newTestWords []string
			for _, word := range testWords {
				isExcluded := false
				for _, mismatchRegex := range mismatchRegexList {
					match, _ := regexp.MatchString(mismatchRegex, word)
					if match {
						isExcluded = true
						break
					}
				}
				if !isExcluded {
					newTestWords = append(newTestWords, word)
				}
			}
			testWords = newTestWords
		}
		// filter down to words with the matched letters in the correct places
		if hasMatches {
			var newTestWords []string
			for _, word := range testWords {
				match, _ := regexp.MatchString(matches, word)
				if match {
					newTestWords = append(newTestWords, word)
				}
			}
			testWords = newTestWords
		}
		matchedWords = testWords
	}
	return matchedWords
}

func displayResult(result GuessResult) {
	var codes string
	if result.IsMatch {
		codes = "\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9"
	} else {
		codes = ""
		for i := 0; i < 5; i++ {
			if result.Matches[i] {
				codes += "\U0001F7E9"
			} else if result.MisMatches[i] {
				codes += "\U0001F7E8"
			} else {
				codes += "⬛"
			}
		}
	}
	fmt.Fprintf(os.Stderr, "%s\n%s\n", strings.ToUpper(result.Guess), codes)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {

	var answers []string
	var acceptedGuesses []string
	if len(os.Args) > 2 {
		guessWord := os.Args[1]
		target := os.Args[2]
		result, err := guess(guessWord, target)
		if err == nil {
			displayResult(result)
			if len(os.Args) > 4 {
				answers, _ = readLines(os.Args[3])
				f, err := os.Create(os.Args[4])
				if err != nil {
					panic(err)
				}
				defer f.Close()
				w := bufio.NewWriter(f)
				matches := findMatches(result, answers)
				n := len(matches)
				fmt.Fprintf(os.Stderr, "Possibilities: %d\n", n)
				for _, match := range matches {
					_, err := w.WriteString(match+"\n")
					if err != nil {
						panic(err)
					}
				}
				w.Flush()
			}
		} else {
			panic(err)
		}
		os.Exit(0)
	}


	if len(os.Args) > 1 {
		answers, _ = readLines(os.Args[1])
		acceptedGuesses = answers
	} else {
		answers, _ = readLines("answers.txt")
		acceptedGuesses, _ = readLines("guesses.txt")
	}
	fmt.Fprintf(os.Stderr, "Answers: %d\n", len(answers))
	fmt.Fprintf(os.Stderr, "Accepted Guesses: %d\n", len(acceptedGuesses))

	/*result, err := guess("coate", "those")
	if err == nil {
		matches := findMatches(result, answers)
		n := len(matches)
		fmt.Fprintf(os.Stderr, "Matches: %d\n", n)
		for _, match := range matches {
			fmt.Fprintln(os.Stderr, match)
		}
	}*/


	for _, guessWord := range acceptedGuesses {
		total := 0
		high := 1
		low := len(answers)
		for _, target := range answers {
			result, err := guess(guessWord, target)
			if err == nil {
				matches := findMatches(result, answers)
				n := len(matches)
				total += n
				if n < low {
					low = n
				}
				if n > high {
					high = n
				}
			}
		}
		avg := float32(total) / float32(len(answers))
		fmt.Fprintf(os.Stderr, "Guess: %s Averge: %.2f low: %d high: %d\n", guessWord, avg, low, high)
	}

	/*
	result, err := guess("grist", "prism")
	if err == nil {
		displayResult(result)
		matches := findMatches(result, answers)
		fmt.Fprintf(os.Stderr, "Matches: %d\n", len(matches))
		for i, match := range matches {
			fmt.Fprintf(os.Stderr, "%d: %s\n", i+1, match)
		}
	}
	 */
}

