package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const charHeight = 8

func main() {
	var inputString string
	var bannerName string = "standard"
	var align string = "left" // Default alignment

	args := os.Args[1:]

	// Обработка аргументов командной строки
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--align=") {
			align = strings.TrimPrefix(args[i], "--align=")
			if align != "center" && align != "left" && align != "right" && align != "justify" {
				fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
				return
			}
		} else if inputString == "" {
			inputString = args[i]
		} else if bannerName == "standard" {
			bannerName = args[i]
		} else {
			fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
			return
		}
	}

	if inputString == "" {
		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		return
	}

	asciiArtMap, err := loadBanner(bannerName)
	if err != nil {
		fmt.Println("Error loading banner:", err)
		return
	}

	terminalWidth := getTerminalWidth()
	if terminalWidth <= 0 {
		fmt.Println("Could not determine terminal width, defaulting to no wrapping.")
	}

	output := generateAsciiArt(inputString, asciiArtMap, align, terminalWidth)
	fmt.Print(output)
}

func loadBanner(bannerName string) (map[rune][]string, error) {
	asciiArtMap := make(map[rune][]string)
	bannerFileName := bannerName + ".txt"
	data, err := os.ReadFile(bannerFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read banner file '%s': %w", bannerFileName, err)
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 95*(charHeight+1) {
		return nil, fmt.Errorf("invalid banner file format")
	}
	for i := 0; i < 95; i++ {
		startLine := i*(charHeight+1) + 1
		charRune := rune(32 + i)
		asciiArtMap[charRune] = lines[startLine : startLine+charHeight]
	}
	return asciiArtMap, nil
}

func generateCharLines(text string, asciiArtMap map[rune][]string) []string {
	charLines := make([]string, charHeight)
	for _, char := range text {
		if art, ok := asciiArtMap[char]; ok {
			for i := 0; i < charHeight; i++ {
				charLines[i] += art[i]
			}
		} else {
			// Handle characters not found in the banner
			for i := 0; i < charHeight; i++ {
				spaceWidth := 1 // Default to 1 if space not found (should not happen with proper banner)
				if spaceArt, ok := asciiArtMap[' ']; ok && len(spaceArt) > 0 {
					spaceWidth = len(spaceArt[0])
				}
				charLines[i] += strings.Repeat("?", spaceWidth)
			}
		}
	}
	return charLines
}

func alignText(textLines []string, align string, terminalWidth int) []string {
	if terminalWidth <= 0 {
		return textLines // No alignment if terminal width is not available
	}

	alignedLines := make([]string, len(textLines))
	for i, line := range textLines {
		lineWidth := utf8.RuneCountInString(line)
		padding := terminalWidth - lineWidth

		switch align {
		case "center":
			if padding > 0 {
				leftPadding := padding / 2
				rightPadding := padding - leftPadding
				alignedLines[i] = strings.Repeat(" ", leftPadding) + line + strings.Repeat(" ", rightPadding)
			} else {
				alignedLines[i] = line
			}
		case "right":
			if padding > 0 {
				alignedLines[i] = strings.Repeat(" ", padding) + line
			} else {
				alignedLines[i] = line
			}
		case "justify":
			// Justify is handled at the word level for ASCII art
			alignedLines[i] = line
		default: // "left"
			if padding > 0 {
				alignedLines[i] = line + strings.Repeat(" ", padding)
			} else {
				alignedLines[i] = line
			}
		}
	}
	return alignedLines
}

// and adds a fixed value of -1 to the width.
func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0 // Unable to get terminal width
	}
	s := strings.TrimSpace(string(out))
	fields := strings.Split(s, " ")
	if len(fields) != 2 {
		return 0 // Unexpected output format
	}
	width, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0 // Unable to parse width
	}
	width += -1 // Add a fixed value of -1
	if width < 0 {
		width = 0 // Ensure the width is non-negative
	}
	return width
}

func justifyWordArts(wordArts [][]string, terminalWidth int) []string {
	if len(wordArts) <= 1 {
		return wordArts[0]
	}

	// Calculate total width of words
	totalWordWidth := 0
	for _, wordArt := range wordArts {
		if len(wordArt) > 0 {
			totalWordWidth += utf8.RuneCountInString(wordArt[0])
		}
	}

	spaceAvailable := terminalWidth - totalWordWidth
	if spaceAvailable <= 0 {
		// Not enough space to justify, return combined words
		combined := make([]string, charHeight)
		for i := 0; i < charHeight; i++ {
			for _, wordArt := range wordArts {
				if len(wordArt) > i {
					combined[i] += wordArt[i]
				}
			}
		}
		return combined
	}

	gaps := len(wordArts) - 1
	spacesPerGap := spaceAvailable / gaps
	//extraSpaces := spaceAvailable % gaps

	justifiedLines := make([]string, charHeight)
	for i := 0; i < charHeight; i++ {
		var line strings.Builder
		for wordIndex, wordArt := range wordArts {
			if len(wordArt) > i {
				line.WriteString(wordArt[i])
			}
			if wordIndex < len(wordArts)-1 {
				spaces := spacesPerGap
				// if extraSpaces > 0 {
				// 	spaces++
				// 	extraSpaces--
				// }

				// fmt.Println(spaces)
				line.WriteString(strings.Repeat(" ", spaces))
			}
		}
		justifiedLines[i] = line.String()
	}
	return justifiedLines
}

func generateAsciiArt(input string, asciiArtMap map[rune][]string, align string, terminalWidth int) string {
	var output strings.Builder
	lines := regexp.MustCompile(`\\n|\n`).Split(input, -1)

	spaceArt := asciiArtMap[' ']

	// Обрабатываем все строки, включая первую, с выравниванием
	for _, line := range lines {
		if line == "" {
			// Пропускаем пустые строки
			continue
		}

		// Разбиваем строку на слова
		words := strings.Fields(line)
		wordArts := make([][]string, len(words))

		// Формируем ASCII-арт для слов
		for i, word := range words {
			wordArts[i] = generateCharLines(word, asciiArtMap)
		}

		// Генерация строк ASCII-арта
		asciiLines := make([]string, charHeight)
		for i := 0; i < charHeight; i++ {
			var currentLine strings.Builder
			for wordIndex, wordArt := range wordArts {
				if len(wordArt) > i {
					currentLine.WriteString(wordArt[i])
				}
				if wordIndex < len(wordArts)-1 {
					currentLine.WriteString(spaceArt[i])
				}
			}
			asciiLines[i] = currentLine.String()
		}

		// Применяем выравнивание для первой строки, если оно включено
		if align == "justify" && terminalWidth > 0 {
			// Выравнивание применяется к первой и последующим строкам одинаково
			asciiLines = justifyWordArts(wordArts, terminalWidth)
		} else {
			asciiLines = alignText(asciiLines, align, terminalWidth)
		}

		// Добавляем строки ASCII-арта в вывод
		for _, asciiLine := range asciiLines {
			output.WriteString(asciiLine + "\n")
		}
	}

	// Удаляем лишние переносы строк в конце
	result := strings.TrimRight(output.String(), "\n")
	return result
}
