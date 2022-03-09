package main

import (
	"bufio"
	"flag"
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

func matcher(answers []string, mustContainLetters []rune, excluded string,
	hasMismatches bool, mismatchRegexList []string,
	hasMatches bool, matches string) []string {
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
	return testWords
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
				if !strings.Contains(string(mustContainLetters), string(result.GuessRune[i])) {
					mustContainLetters = append(mustContainLetters, result.GuessRune[i])
				}
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
					if !strings.Contains(string(mustContainLetters), string(result.GuessRune[i])) {
						mustContainLetters = append(mustContainLetters, result.GuessRune[i])
					}
					hasMismatches = true
				} else if !strings.Contains(excluded, string(result.GuessRune[i])) {
					excluded += string(result.GuessRune[i])
				}
			}
		}
		matchedWords = matcher(answers, mustContainLetters, excluded,
			hasMismatches, mismatchRegexList,
			hasMatches, matches)
	}
	return matchedWords
}

func removeLetter(letters []rune, letter rune)  []rune {
	for i,l := range letters {
		if l == letter {
			letters[i] = letters[len(letters)-1]
			return append(letters[:i], letters[i+1:]...)
		}
	}
	return letters
}

func displayResult(result GuessResult) {
	var codes string
	if result.IsMatch {
		codes = "\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9\U0001F7E9"
	} else {
		codes = ""
		availableLetters := result.TargetRune
		for i := 0; i < 5; i++ {
			if result.Matches[i] {
				availableLetters = removeLetter(availableLetters, result.GuessRune[i])
			}
		}
		for i := 0; i < 5; i++ {
			if result.Matches[i] {
				codes += "\U0001F7E9"
			} else if result.MisMatches[i] && strings.Contains(string(availableLetters), string(result.GuessRune[i])) {
				codes += "\U0001F7E8"
				availableLetters = removeLetter(availableLetters, result.GuessRune[i])
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
	var vocabGuesses []string

	var easyMode = flag.Bool("e", false, "easy mode")

	flag.Parse()

	// special case helper mode (3 args)
	if flag.NArg() == 3 {
		answers, _ = readLines(flag.Arg(2))
		guessWord := flag.Arg(0)
		guessRune := []rune(guessWord)
		clue := flag.Arg(1)
		clueRune := []rune(clue)
		matches := "^"
		var mismatchRegexList []string
		excluded := ""
		var mustContainLetters []rune
		hasMatches := false
		hasMismatches := false
		for i, c := range clueRune {
			//fmt.Fprintf(os.Stderr, "clue[%d]: 0x%x\n", i, c)
			// match
			if c == 0x1f7e9 || c == 0x67 {
				matches += string(guessRune[i])
				if !strings.Contains(string(mustContainLetters), string(guessRune[i])) {
					mustContainLetters = append(mustContainLetters, guessRune[i])
				}
				hasMatches = true
			} else {
				matches += "."
				// mismatch
				if c == 0x1f7e8 || c == 0x79 {
					mismatchRegex := "^"
					for j := 0; j < 5; j++ {
						if j == i {
							mismatchRegex += string(guessRune[i])
						} else {
							mismatchRegex += "."
						}
					}
					mismatchRegexList = append(mismatchRegexList, mismatchRegex)
					if !strings.Contains(string(mustContainLetters), string(guessRune[i])) {
						mustContainLetters = append(mustContainLetters, guessRune[i])
					}
					hasMismatches = true
				}
			}
		}
		for i, c := range clueRune {
			// not present
			//fmt.Fprintf(os.Stderr, " test excluded: clue[%d]: 0x%x\n", i, c)
			/*if c != 0x1f7e9 && c != 0x1f7e8 && c != 0x67 && c != 0x79 {
				fmt.Fprintf(os.Stderr, "possible excluded: %s\n", string(guessRune[i]))
			}*/
			if c != 0x1f7e9 && c != 0x1f7e8 && c != 0x67 && c != 0x79 &&
				!strings.Contains(string(mustContainLetters), string(guessRune[i])) &&
				!strings.Contains(excluded, string(guessRune[i])) {
				excluded += string(guessRune[i])
			}
		}
		/*fmt.Fprintf(os.Stderr, "excluded: %s mustInclude: %s matches: %s\n",
			excluded, string(mustContainLetters), matches)*/
		answers = matcher(answers, mustContainLetters, excluded,
			hasMismatches, mismatchRegexList,
			hasMatches, matches)
		if *easyMode {
			vocabGuesses, _ = readLines("words.txt")
		} else {
			vocabGuesses = answers
		}
	} else {
		if flag.NArg() > 1 {
			guessWord := flag.Arg(0)
			target := flag.Arg(1)
			result, err := guess(guessWord, target)
			if err == nil {
				displayResult(result)
				if flag.NArg() > 3 {
					answers, _ = readLines(flag.Arg(2))
					f, err := os.Create(flag.Arg(3))
					if err != nil {
						panic(err)
					}
					defer f.Close()
					w := bufio.NewWriter(f)
					matches := findMatches(result, answers)
					n := len(matches)
					fmt.Fprintf(os.Stderr, "Possibilities: %d\n", n)
					for _, match := range matches {
						_, err := w.WriteString(match + "\n")
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
		if flag.NArg() > 0 {
			answers, _ = readLines(flag.Arg(0))
			if *easyMode {
				vocabGuesses, _ = readLines("guesses.txt")
			} else {
				vocabGuesses = answers
			}
		} else {
			answers, _ = readLines("answers.txt")
			vocabGuesses, _ = readLines("words.txt")
		}
	}
	fmt.Fprintf(os.Stderr, "Answers: %d\n", len(answers))
	fmt.Fprintf(os.Stderr, "Vocabulary: %d\n", len(vocabGuesses))

	/*result, err := guess("coate", "those")
	if err == nil {
		matches := findMatches(result, answers)
		n := len(matches)
		fmt.Fprintf(os.Stderr, "Matches: %d\n", n)
		for _, match := range matches {
			fmt.Fprintln(os.Stderr, match)
		}
	}*/


	for _, guessWord := range vocabGuesses {
		total := 0
		high := 0
		low := len(answers)+1
		lowWord := guessWord
		highWord := guessWord
		for _, target := range answers {
			result, err := guess(guessWord, target)
			if err == nil {
				matches := findMatches(result, answers)
				n := len(matches)
				total += n
				if n < low {
					low = n
					lowWord = target
				}
				if n > high {
					high = n
					highWord = target
				}
			}
		}
		avg := float32(total) / float32(len(answers))
		fmt.Fprintf(os.Stderr, "Guess: %s Averge: %.2f low: %d high: %d (%s/%s)\n", guessWord, avg, low, high, lowWord, highWord)
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

