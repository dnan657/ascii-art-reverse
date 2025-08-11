package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const charHeight = 8

func main() {
	var outputFileName, inputString, bannerName, reverseFile, bannerFlag string
	var positionalArgs []string

	bannerName = "standard" // Default

	args := os.Args[1:]

	// 1. Разбор аргументов
	for _, arg := range args {
		if strings.HasPrefix(arg, "--output=") {
			outputFileName = strings.TrimPrefix(arg, "--output=")
		} else if strings.HasPrefix(arg, "--reverse=") {
			reverseFile = strings.TrimPrefix(arg, "--reverse=")
		} else if strings.HasPrefix(arg, "--banner=") {
			bannerFlag = strings.TrimPrefix(arg, "--banner=")
		} else if strings.HasPrefix(arg, "--") {
			fmt.Printf("Error: Unknown flag %s\n", arg)
			printUsage()
			return
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	// 2. Установка баннера (флаг имеет приоритет)
	if bannerFlag != "" {
		bannerName = bannerFlag
	}

	// 3. Выбор режима работы
	isReverseMode := reverseFile != ""

	if isReverseMode {
		if len(positionalArgs) > 0 {
			if bannerFlag == "" && len(positionalArgs) == 1 {
				bannerName = positionalArgs[0]
			} else {
				fmt.Println("Warning: Positional arguments are ignored when --reverse is used.")
			}
		}
		if reverseFile == "" {
			printUsage()
			return
		}
		result, err := reverseAsciiArt(reverseFile, bannerName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Print(result)
		return
	}

	// Режим генерации
	if len(positionalArgs) > 0 {
		inputString = positionalArgs[0]
	}
	if len(positionalArgs) > 1 {
		if bannerFlag == "" {
			bannerName = positionalArgs[1]
		}
	}
	if len(positionalArgs) > 2 {
		printUsage()
		return
	}
	if inputString == "" {
		printUsage()
		return
	}

	// Генерация
	asciiArtMap, err := loadBanner(bannerName)
	if err != nil {
		fmt.Println("Error loading banner:", err)
		return
	}

	output := generateAsciiArt(inputString, asciiArtMap)

	if outputFileName != "" {
		err := os.WriteFile(outputFileName, []byte(output), 0644)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
		}
	} else {
		fmt.Print(output)
	}
}

func printUsage() {
	fmt.Println("Usage: go run . [OPTION]... [STRING] [BANNER]")
	fmt.Println("\nExample: go run . --output=banner.txt 'Hello' standard")
	fmt.Println("\nOPTIONS:")
	fmt.Println("  --output=<fileName>      Save the output to a file.")
	fmt.Println("  --reverse=<fileName>     Convert ASCII art from a file back to text.")
	fmt.Println("  --banner=<bannerName>    Specify the banner font to use (standard, shadow, thinkertoy).")
}

func reverseAsciiArt(fileName string, bannerName string) (string, error) {
	// Читаем файл с ASCII-арт
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", fileName, err)
	}

	// Загружаем баннер для сравнения
	asciiArtMap, err := loadBanner(bannerName)
	if err != nil {
		return "", fmt.Errorf("failed to load banner '%s': %w", bannerName, err)
	}

	// Разбиваем содержимое файла на строки
	lines := strings.Split(string(data), "\n")

	// Удаляем пустые строки в конце
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return "", nil
	}

	// Группируем строки по блокам высотой charHeight
	var result strings.Builder

	for i := 0; i < len(lines); i += charHeight {
		// Берем блок строк высотой charHeight
		block := make([]string, charHeight)
		for j := 0; j < charHeight && i+j < len(lines); j++ {
			block[j] = lines[i+j]
		}

		// Дополняем блок пустыми строками если нужно
		for j := len(block); j < charHeight; j++ {
			block = append(block, "")
		}

		// Декодируем блок в текст
		text := decodeAsciiBlock(block, asciiArtMap)
		if text != "" {
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(text)
		}
	}

	return result.String(), nil
}

func decodeAsciiBlock(block []string, asciiArtMap map[rune][]string) string {
	if len(block) != charHeight {
		return ""
	}

	// Найдем максимальную ширину блока
	maxWidth := 0
	for _, line := range block {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	if maxWidth == 0 {
		return ""
	}

	// Дополняем все строки до максимальной ширины
	for i := range block {
		for len(block[i]) < maxWidth {
			block[i] += " "
		}
	}

	var result strings.Builder
	pos := 0

	for pos < maxWidth {
		foundChar := false

		// Проверяем каждый символ из ASCII таблицы
		for char := rune(32); char <= rune(126); char++ {
			if charArt, exists := asciiArtMap[char]; exists {
				charWidth := len(charArt[0])

				// Проверяем, помещается ли символ в оставшееся пространство
				if pos+charWidth <= maxWidth {
					// Проверяем соответствие всех строк символа
					matches := true
					for i := 0; i < charHeight; i++ {
						if block[i][pos:pos+charWidth] != charArt[i] {
							matches = false
							break
						}
					}

					if matches {
						result.WriteRune(char)
						pos += charWidth
						foundChar = true
						break
					}
				}
			}
		}

		if !foundChar {
			// Если символ не найден, пропускаем одну позицию
			pos++
		}
	}

	return result.String()
}

// Ваши существующие функции остаются без изменений
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
			for i := 0; i < charHeight; i++ {
				spaceWidth := 1
				if spaceArt, ok := asciiArtMap[' ']; ok && len(spaceArt) > 0 {
					spaceWidth = len(spaceArt[0])
				}
				charLines[i] += strings.Repeat("?", spaceWidth)
			}
		}
	}
	return charLines
}



func generateAsciiArt(input string, asciiArtMap map[rune][]string) string {
	var output strings.Builder
	lines := regexp.MustCompile(`\\n|\n`).Split(input, -1)

	for _, line := range lines {
		if line == "" {
			output.WriteString("\n")
			continue
		}

		asciiLines := generateCharLines(line, asciiArtMap)

		for _, asciiLine := range asciiLines {
			output.WriteString(asciiLine + "\n")
		}
	}

	// We remove the final newline to prevent an extra empty line at the end of the file
	return strings.TrimSuffix(output.String(), "\n")
}
