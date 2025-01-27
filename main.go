package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
	"encoding/hex"
	"log"
	"regexp"
)

var (
	singleWord     = flag.String("w", "", "Single word to munge")
	inputFile      = flag.String("i", "", "Input wordlist file")
	outputFile     = flag.String("o", "", "Output wordlist file")
	all            = flag.Bool("all", false, "Enable all munging techniques")
	capitalizeFlag = flag.Bool("c", false, "Capitalizization")
	substituteFlag = flag.Bool("cs", false, "Use l33t substitutions")
	prependFlag    = flag.Bool("p", false, "Prepend special chars")
	appendFlag     = flag.Bool("a", false, "Append special chars")
	duplicateFlag  = flag.Bool("d", false, "Duplicate word after munging")
	insaneFlag     = flag.Bool("1ns4n3", false, "Generate maximum capitalization and l33t variations (single word only)")
	wordSwapFlag   = flag.Bool("ws", false, "Generate word swaps for multi-word inputs")
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

func generateCapitalizationVariations(word string) []string {
	n := len(word)
	var results []string
	for i := 0; i < (1 << n); i++ {
		var variation []rune
		for j, char := range word {
			if (i>>j)&1 == 1 {
				variation = append(variation, unicode.ToUpper(char))
			} else {
				variation = append(variation, unicode.ToLower(char))
			}
		}
		results = append(results, string(variation))
	}
	return results
}

func l33t(word string) []string {
	var leetMap = map[rune][]string{
		'a': {"@", "4",},
		'A': {"@", "4"},
		'e': {"3"},
		'E': {"3"},
		'i': {"!", "1"},
		'I': {"!", "1"},
		'o': {"0"},
		'O': {"0"},
		's': {"$", "5"},
		'S': {"$", "5"},
	}

	expansions := []string{""}
	for _, ch := range word {
		possibleReplacements, found := leetMap[ch]
		if !found {
			possibleReplacements = []string{string(ch)}
		}
		var newExpansions []string
		for _, partial := range expansions {
			for _, rep := range possibleReplacements {
				newExpansions = append(newExpansions, partial+rep)
			}
		}
		expansions = newExpansions
	}
	return expansions
}

func applyMunging(word string) []string {
	var results []string

	if *insaneFlag {
		if *singleWord != "" {
			capitalizationVariations := generateCapitalizationVariations(word)
			for _, variation := range capitalizationVariations {
				results = append(results, l33t(variation)...) // Apply l33t to each capitalization variation
			}
		} else {
			fmt.Println("ERROR: --1ns4n3 flag can only be used with single-word inputs.")
			os.Exit(1)
		}
	} else {
		if *capitalizeFlag {
			results = append(results, strings.Title(word)) // Standard capitalization
			results = append(results, swapCase(word)) // Swap case
		}

		if *substituteFlag  {
			results = append(results, l33t(word)...) // Leet substitutions only if capitalization variations are skipped
		}
	}

	if *duplicateFlag {
		results = append(results, word+word)
	}

	if strings.Contains(word, " ") {
		results = append(results, strings.ReplaceAll(word, " ", ""))
		results = append(results, strings.ReplaceAll(word, " ", "^"))
		results = append(results, strings.ReplaceAll(word, " ", "."))
	}

	return results
}

func generateWordSwaps(word string) []string {
	fields := strings.Fields(word)
	if len(fields) < 2 {
		return nil // No swap possible for single-word inputs
	}
	var swaps []string
	for i := 0; i < len(fields); i++ {
		for j := i + 1; j < len(fields); j++ {
			swapped := make([]string, len(fields))
			copy(swapped, fields)
			swapped[i], swapped[j] = swapped[j], swapped[i]
			swaps = append(swaps, strings.Join(swapped, " "))
		}
	}
	return swaps
}

func appendPrepend(base string, writer *bufio.Writer){
	numbers := []string{
		"20", "21", "22", "23", "24", "25", "26",
		"1", "12", "21", "123", "321", "1234", "4321", "12345", "54321", "123456", "654321", "1234567", "7654321", "12345678", "87654321", "123456789", "987654321",
		"2018", "2019", "2020", "2021", "2022", "2023", "2024", "2025", "2026",
	}
	specialChars := []string{
		" ", "_", ".", "*", "&", "&&",
		"!", "!!", "!!!", "!!!!", "!!!!!",
		"@", "@@", "@@@",
		"#", "##", "###",
		"$", "$$", "$$$", "$$$$", "$$$$$",
		"!@", "!@#", "!@#$", "!@#$%", "$#@!", "#@!", "@!",
	}

	for _, num := range numbers {
		for _, char := range specialChars {
			if *prependFlag {
				writer.WriteString(num + base + "\n")
				writer.WriteString(num + char + base + "\n")
				writer.WriteString(char + num + base + "\n")
			}
			if *appendFlag {
				writer.WriteString(base + num + "\n")
				writer.WriteString(base + num + char + "\n")
				writer.WriteString(base + char + num + "\n")
			}
			if *prependFlag && *appendFlag {
				writer.WriteString(char + num + base + num + char + "\n")
				writer.WriteString(char + num + base + char + num + "\n")
				writer.WriteString(num + char + base + char + num + "\n")
				writer.WriteString(num + char + base + num + char + "\n")
			}
		}
	}
}

func munge(word string, writer *bufio.Writer) {
	if strings.HasPrefix(word, "$HEX[") {
		// Decode the hex string to bytes
		re := regexp.MustCompile(`\[(.*?)\]`)
		match := re.FindStringSubmatch(word)
		bytes, err := hex.DecodeString(match[1])
		if err != nil {
			log.Print(err)
		}
		word = string(bytes)
	}
	wordsToProcess := []string{word}
	if *wordSwapFlag {
		wordSwaps := generateWordSwaps(word)
		if wordSwaps != nil {
			wordsToProcess = append(wordsToProcess, wordSwaps...)
		}
	}
	writer.WriteString(word + "\n")
	for _, baseWord := range wordsToProcess {
		for _, base := range applyMunging(baseWord) {
			writer.WriteString(base + "\n")
			if (*appendFlag || *prependFlag){
				appendPrepend(base, writer)
			}
		}
	}
}

func main() {
	flag.Parse()

	if *all {
		*capitalizeFlag = true
		*substituteFlag = true
		*prependFlag = true
		*appendFlag = true
		*duplicateFlag = true
		*wordSwapFlag = true
	} else if !*capitalizeFlag && !*substituteFlag && !*prependFlag && !*appendFlag && !*duplicateFlag && !*insaneFlag && !*wordSwapFlag {
		fmt.Println("ERROR: Please specify at least one option for modification (or use -all).")
		flag.Usage()
		os.Exit(1)
	}

	if *outputFile == "" {
		fmt.Println("ERROR: No output file specified. Use -o <file>")
		os.Exit(1)
	}

	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	if *singleWord != "" {
		munge(*singleWord, writer)
	}

	if *inputFile != "" {
		inFile, err := os.Open(*inputFile)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer inFile.Close()

		scanner := bufio.NewScanner(inFile)
		for scanner.Scan() {
			munge(scanner.Text(), writer)
		}
		if err = scanner.Err(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	if *singleWord == "" && *inputFile == "" {
		fmt.Println("ERROR: Provide a single word with -w or an input file with -i.")
		flag.Usage()
		os.Exit(1)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Munged wordlist saved to", *outputFile)
}
