package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"flag"
	"unicode"
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

func swapCase(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsUpper(r) {
			result = append(result, unicode.ToLower(r))
		} else if unicode.IsLower(r) {
			result = append(result, unicode.ToUpper(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

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

func applyAllMungeTechniques(word string) []string {
	techniques := []string{word}
	if *capitalize {
		techniques = append(techniques, strings.Title(word))
	}
	if *substitute {
		techniques = append(techniques, l33t(word))
	}
	
	techniques = append(techniques, swapCase(word))

	if *duplicate {
		duplicatedWord := word + word
		techniques = append(techniques, duplicatedWord)
	}

	return techniques
}

func munge(word string, writer *bufio.Writer) {
	specialChars := []string{"!", "@", "#", "$", "%"}
	suffixes := []string{"123", "1", "69", "21", "22", "23", "2019", "2020", "2021", "2022", "2023", "1984", "1985", "1986","1987","1988"}

	for _, baseWord := range applyAllMungeTechniques(word) {
		writer.WriteString(baseWord + "\n")

		for _, suffix := range suffixes {
			writer.WriteString(baseWord + suffix + "\n")
		}

		for _, char := range specialChars {
			if *prependSpecial {
				writer.WriteString(char + baseWord + "\n")
			}
			if *appendSpecial {
				writer.WriteString(baseWord + char + "\n")
			}
			if *prependSpecial && *appendSpecial {
				writer.WriteString(char + baseWord + char + "\n")
			}
		}
	}
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
		word = strings.ToLower(word)
		munge(word, writer)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Munged wordlist saved to", *outputFile)
}

