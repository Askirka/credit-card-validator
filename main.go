package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ===== Структура =====
type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func loadBankData(filename string) ([]Bank, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var banks []Bank

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}

		binFrom, _ := strconv.Atoi(parts[1])
		binTo, _ := strconv.Atoi(parts[2])

		bank := Bank{
			Name:    parts[0],
			BinFrom: binFrom,
			BinTo:   binTo,
		}

		banks = append(banks, bank)
	}

	return banks, nil
}

func validateInput(input string) bool {
	if len(input) < 13 || len(input) > 19 {
		return false
	}

	for _, ch := range input {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}

func validateLuhn(a string) bool {
	sum := 0
	counter := 0

	for i := len(a) - 1; i >= 0; i-- {
		digit := int(a[i] - '0')

		if counter%2 == 1 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		counter++
	}

	return sum%10 == 0
}

func extractBIN(card string) int {
	bin, _ := strconv.Atoi(card[:6])
	return bin
}

func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Unknown bank"
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите номер карты (или Enter для выхода): ")

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {

	fmt.Println("=== Валидатор кредитных карт ===")

	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка загрузки данных банков:", err)
		return
	}

	fmt.Println("Загружено банков:", len(banks))

	for {
		input := getUserInput()

		if input == "" {
			fmt.Println("До свидания!")
			break
		}

		if !validateInput(input) {
			if len(input) < 13 || len(input) > 19 {
				fmt.Println("✗ Ошибка: номер должен содержать от 13 до 19 цифр")
			} else {
				fmt.Println("✗ Ошибка: номер должен содержать только цифры")
			}
			continue
		}
		fmt.Println("✓ Формат корректен")

		if !validateLuhn(input) {
			fmt.Println("✗ Номер невалиден (не прошёл проверку Луна)")
			continue
		}
		fmt.Println("✓ Номер валиден (алгоритм Луна)")

		bin := extractBIN(input)
		bank := identifyBank(bin, banks)

		fmt.Println("✓ Банк:", bank)
	}
}
