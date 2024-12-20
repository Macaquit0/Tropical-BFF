package helper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	cpfFirstDigitTable   = []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	cpfSecondDigitTable  = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjFirstDigitTable  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondDigitTable = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

func sumDigit(s string, table []int) int {
	if len(s) != len(table) {
		return 0
	}

	sum := 0

	for i, v := range table {
		c := string(s[i])
		d, err := strconv.Atoi(c)
		if err == nil {
			sum += v * d
		}
	}

	return sum
}

func ValidateCNPJ(cnpj string) bool {
	blacklisted := []string{
		"00000000000000",
		"11111111111111",
		"22222222222222",
		"33333333333333",
		"44444444444444",
		"55555555555555",
		"66666666666666",
		"77777777777777",
		"88888888888888",
		"99999999999999",
	}

	cleanCnpj := OnlyNumbers(cnpj)

	if len(cleanCnpj) != 14 {
		return false
	}

	for _, blacklist := range blacklisted {
		if blacklist == cleanCnpj {
			return false
		}
	}

	firstPart := cleanCnpj[:12]
	sum1 := sumDigit(firstPart, cnpjFirstDigitTable)
	rest1 := sum1 % 11
	d1 := 0

	if rest1 >= 2 {
		d1 = 11 - rest1
	}

	secondPart := fmt.Sprintf("%s%d", firstPart, d1)
	sum2 := sumDigit(secondPart, cnpjSecondDigitTable)
	rest2 := sum2 % 11
	d2 := 0

	if rest2 >= 2 {
		d2 = 11 - rest2
	}

	finalPart := fmt.Sprintf("%s%d", secondPart, d2)
	return finalPart == cleanCnpj
}

func ValidateCPF(cpf string) bool {
	blacklisted := []string{
		"00000000000",
		"11111111111",
		"22222222222",
		"33333333333",
		"44444444444",
		"55555555555",
		"66666666666",
		"77777777777",
		"88888888888",
		"99999999999",
	}

	cleanCPF := OnlyNumbers(cpf)

	if len(cleanCPF) != 11 {
		return false
	}

	for _, blacklist := range blacklisted {
		if blacklist == cleanCPF {
			return false
		}
	}

	firstPart := cleanCPF[0:9]
	sum := sumDigit(firstPart, cpfFirstDigitTable)

	r1 := sum % 11
	d1 := 0

	if r1 >= 2 {
		d1 = 11 - r1
	}

	secondPart := firstPart + strconv.Itoa(d1)

	dsum := sumDigit(secondPart, cpfSecondDigitTable)

	r2 := dsum % 11
	d2 := 0

	if r2 >= 2 {
		d2 = 11 - r2
	}

	finalPart := fmt.Sprintf("%s%d%d", firstPart, d1, d2)
	return finalPart == cleanCPF
}

func FormatDocumentNumber(documentNumber string) string {
	if strings.Contains(documentNumber, "X") {
		return strings.ReplaceAll(documentNumber, "X", "*")
	}

	if strings.Contains(documentNumber, "*") {
		return documentNumber
	}

	documentNumberClean := OnlyNumbers(documentNumber)

	if len(documentNumberClean) > 11 {
		return formatCNPJ(documentNumberClean)
	} else {
		return formatCPF(documentNumberClean)
	}
}

func formatCPF(cpf string) string {
	regExp := regexp.MustCompile(`^(\d{3})(\d{3})(\d{3})(\d{2})$`)

	return regExp.ReplaceAllString(OnlyNumbers(cpf), "***.$2.$3-**")
}

func formatCNPJ(cnpj string) string {
	regExp := regexp.MustCompile(`^(\d{2})(\d{3})(\d{3})(\d{4})(\d{2})$`)

	return regExp.ReplaceAllString(OnlyNumbers(cnpj), "$1.$2.$3/$4-$5")
}

func ValidateTaxId(taxID string) bool {
	cleanDocument := OnlyNumbers(taxID)
	documentLen := len(cleanDocument)
	switch documentLen {
	case 11:
		return ValidateCPF(taxID)
	case 14:
		return ValidateCNPJ(taxID)
	default:
		return false
	}
}
