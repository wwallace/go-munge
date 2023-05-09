package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"flag"
)

var (
	inputFile      = flag.String("i", "", "input wordlist file")
	outputFile     = flag.String("o", "", "output wordlist file")
	all            = flag.Bool("all", true, "apply all munging techniques (default)")
	capitalize     = flag.Bool("c", false, "capitalize the first letter")
	substitute     = flag.Bool("cs", false, "substitute characters with l33t speak")
	prependSpecial = flag.Bool("p", false, "prepend special characters")
	appendSpecial  = flag.Bool("a", false, "append special characters")
	duplicate      = flag.Bool("d", false, "duplicate the word and apply munging techniques")
)

func l33t(word string) string {
	substitutions := map[string]string{
		"a": "@", "A": "4",
		"e": "3", "E": "3",
		"i": "!", "I": "1",
		"o": "0", "O": "0",
		"s": "$", "S": "5",
	}

	for k, v := range substitutions {
		word = strings.Replace(word, k, v, -1)
	}
	return word
}

func munge(word string) []string {
	var result []string

	// Add the original word to the result
	result = append(result, word)

	if *capitalize {
		word = strings.Title(word)
	}

	// Handle the munging techniques without duplication
	result = applyMungingTechniques(word, result)

	// If duplication is enabled, apply the munging techniques to the duplicated word
	if *duplicate {
		duplicatedWord := word + word
		result = applyMungingTechniques(duplicatedWord, result)
	}

	return result
}

func applyMungingTechniques(word string, result []string) []string {
	var tempResults []string

	if *capitalize {
		tempResults = append(tempResults, strings.Title(word))
	}

	if *substitute {
		tempResults = append(tempResults, l33t(word))
	}

	for _, tempResult := range tempResults {
		if *prependSpecial || *appendSpecial {
			specialChars := []string{"!", "@", "#", "$", "%"}

			for _, char := range specialChars {
				if *prependSpecial {
					result = append(result, char+tempResult)
				}

				if *appendSpecial {
					result = append(result, tempResult+char)
				}
			}
		} else {
			result = append(result, tempResult)
		}
	}

	return result
}



func main() {
	flag.Parse()

	if *all {
		*capitalize = true
		*substitute = true
		*prependSpecial = true
		*appendSpecial = true
		*duplicate = true
	}
	

	if *inputFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	inFile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer inFile.Close()

	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		word := scanner.Text()
		mungedWords := munge(word)

		for _, mungedWord := range mungedWords {
			_, err := writer.WriteString(mungedWord + "\n")
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Munged wordlist saved to", *outputFile)
}
