package utils

import (
	"strings"
)

type BigNumberValue interface {
    ~string | BigNumber
}

// BigNumberConfig holds configuration for BigNumber formatting.
type BigNumberConfig struct {
	alphabet         string
	groupSeparator   string
	decimalSeparator string
	base             int
	maxDecimal       int
}

// BigNumber represents a large number with potential decimals as a string.
type BigNumber struct {
	value  string
	format *BigNumberConfig
	sign   int
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// NewBigNumber creates a new BigNumber instance.
func NewBigNumber(value string) *BigNumber {
	sign := 1
	if strings.HasPrefix(value, "-") {
		sign = -1
		value = value[1:]
	}
	value = strings.TrimLeft(value, "0")
	if value == "" || value == "." {
		value = "0"
	}
	if strings.HasPrefix(value, ".") {
		value = "0" + value
	}
	bignumber := &BigNumber{
		value: value,
		format: &BigNumberConfig{
			alphabet:         "0123456789abcdefghijklmnopqrstuvwxyz",
			groupSeparator:   ",",
			decimalSeparator: ".",
			base:             10,
            maxDecimal:       6,
		},
		sign: sign,
	}
	return bignumber
}

func (a *BigNumber) SetConfig(config *BigNumberConfig) {
	if config == nil {
		a.format = &BigNumberConfig{
			alphabet:         "0123456789abcdefghijklmnopqrstuvwxyz",
			groupSeparator:   ",",
			decimalSeparator: ".",
			base:             10,
			maxDecimal:       6,
		}
		return
	}
	a.format = config
}

// subtractStrings subtracts two number strings and returns the result.
func (a *BigNumber) subtractStrings(num1 string, num2 string) string {
	// Assume num1 >= num2 for simplicity
	i, j := len(num1)-1, len(num2)-1
	borrow := 0
	result := ""

	for i >= 0 || j >= 0 {
		diff := borrow
		if i >= 0 {
			diff += a.toInteger(string(num1[i]))
			i--
		}
		if j >= 0 {
			diff -= a.toInteger(string(num2[j]))
			j--
		}
		if diff < 0 {
			diff += a.format.base
			borrow = -1
		} else {
			borrow = 0
		}
		result = a.toNumber(diff) + result
	}

	return strings.TrimLeft(result, "0")
}

// multiplyStrings multiplies two number strings and returns the result.
func (a *BigNumber) multiplyStrings(num1 string, num2 string) string {
	intPart1, decPart1 := splitDecimal(num1, a.format.decimalSeparator)
	intPart2, decPart2 := splitDecimal(num2, a.format.decimalSeparator)

	num1 = intPart1 + decPart1
	num2 = intPart2 + decPart2

	result := NewBigNumber("0")
	for i := len(num2) - 1; i >= 0; i-- {
		digit := a.toInteger(string(num2[i]))
		tempResult := strings.Repeat("0", len(num2)-1-i)
		carry := 0
		for j := len(num1) - 1; j >= 0; j-- {
			mul := digit*a.toInteger(string(num1[j])) + carry
			carry = mul / a.format.base
			tempResult = a.toNumber(mul%a.format.base) + tempResult
		}
		if carry > 0 {
			tempResult = a.toNumber(carry) + tempResult
		}
		result = result.add(NewBigNumber(tempResult))
	}

	decPlaces := len(decPart1) + len(decPart2)
	if decPlaces > 0 {
		result.value = insertDecimalPoint(result.value, decPlaces, a.format.decimalSeparator)
	}

	return result.value
}

// divideStrings divides two number strings and returns the quotient and remainder.
func (a *BigNumber) divideStrings(num1 string, num2 string) (*BigNumber, *BigNumber) {
	if num2 == "0" {
		panic("division by zero")
	}

	num1, num2 = normalizeDecimal(num1, num2, a.format.decimalSeparator)
	intPart1, decPart1 := splitDecimal(num1, a.format.decimalSeparator)
	intPart2, decPart2 := splitDecimal(num2, a.format.decimalSeparator)

	lenDec1 := len(decPart1)
	lenDec2 := len(decPart2)

	num1 = intPart1 + decPart1
	num2 = intPart2 + decPart2

	shift := lenDec2 - lenDec1
	if shift > 0 {
		num1 += strings.Repeat("0", shift)
	} else if shift < 0 {
		num2 += strings.Repeat("0", -shift)
	}

	quotient := ""
	remainder := 0
	for i := 0; i < len(num1)+a.format.maxDecimal; i++ {
		var digit1 int
		if i < len(num1) {
			digit1 = a.toInteger(string(num1[i]))
		}
		remainder = remainder*a.format.base + digit1
		digitResult := remainder / a.toInteger(string(num2[0]))
		remainder = remainder % a.toInteger(string(num2[0]))
		quotient += a.toNumber(digitResult)
	}

	quotient = insertDecimalPoint(quotient, a.format.maxDecimal, a.format.decimalSeparator)
	quotient = strings.TrimRight(quotient, "0")
	if strings.HasSuffix(quotient, a.format.decimalSeparator) {
		quotient = strings.TrimSuffix(quotient, a.format.decimalSeparator)
	}

	quotientBigNumber := NewBigNumber(quotient)

	remainderBigNumber := NewBigNumber(a.toNumber(remainder))

	return quotientBigNumber, remainderBigNumber
}

// AbsoluteValue returns the absolute value of the BigNumber.
func (a *BigNumber) AbsoluteValue() *BigNumber {
	if a.sign < 0 {
		return NewBigNumber(a.value)
	}
	return a
}

func (a *BigNumber) toInteger(number string) int {
	index := strings.Index(a.format.alphabet, number)
	if index == -1 {
		return 0
	}
	return index
}

func (a *BigNumber) toNumber(index int) string {
	if index > len(a.format.alphabet) {
		panic("out of range")
	}
	return string(a.format.alphabet[index])
}

// ComparedTo compares two BigNumber instances.
func (a *BigNumber) ComparedTo(b *BigNumber) int {
	if a.sign != b.sign {
		if a.sign > b.sign {
			return a.sign
		}
		return b.sign
	}
	intPart1, decPart1 := splitDecimal(a.value, a.format.decimalSeparator)
	intPart2, decPart2 := splitDecimal(b.value, b.format.decimalSeparator)

	if intPart1 != intPart2 {
		if len(intPart1) != len(intPart2) {
			return len(intPart1) - len(intPart2)
		}
		for i := 0; i < len(intPart1); i++ {
			if intPart1[i] != intPart2[i] {
				return a.toInteger(string(intPart1)) - b.toInteger(string(intPart2[i]))
			}
		}
	}

	for i := 0; i < max(len(decPart1), len(decPart2)); i++ {
		d1, d2 := 0, 0
		if i < len(decPart1) {
			d1 = a.toInteger(string(decPart1[i]))
		}
		if i < len(decPart2) {
			d2 = b.toInteger(string(decPart2[i]))
		}
		if d1 != d2 {
			return d1 - d2
		}
	}

	return 0
}

// DecimalPlaces returns the number of decimal places.
func (a *BigNumber) DecimalPlaces() int {
	parts := strings.Split(a.value, a.format.decimalSeparator)
	if len(parts) == 2 {
		return len(parts[1])
	}
	return 0
}

func (a *BigNumber) SetBase(base int) {
	a.format.base = base
}

// DividedBy divides the BigNumber by another BigNumber.
func (a *BigNumber) DividedBy(b *BigNumber) *BigNumber {
	quotient := a.Divide(b)
	return quotient
}

// DividedToIntegerBy divides the BigNumber by another BigNumber and returns the integer quotient.
func (a *BigNumber) DividedToIntegerBy(b *BigNumber) *BigNumber {
	quotient := a.Divide(b)
	return quotient.IntegerValue()
}

// Add adds two BigNumber instances and returns a new BigNumber instance.
func (a *BigNumber) add(b *BigNumber) *BigNumber {
	num1, num2 := normalizeDecimal(a.value, b.value, a.format.decimalSeparator)
	i, j := len(num1)-1, len(num2)-1
	carry := 0
	result := ""

	decPlaces := 0
	hasDecimal := strings.Contains(a.value, a.format.decimalSeparator) || strings.Contains(b.value, a.format.decimalSeparator)
	if hasDecimal {
		_, decPart1 := splitDecimal(a.value, a.format.decimalSeparator)
		_, decPart2 := splitDecimal(b.value, a.format.decimalSeparator)
		decPlaces = max(len(decPart1), len(decPart2))
	}

	for i >= 0 || j >= 0 || carry > 0 {
		var digit1, digit2 int
		if i >= 0 {
			if string(num1[i]) == a.format.decimalSeparator {
				i--
				j--
				continue
			}
			digit1 = a.toInteger(string(num1[i]))
			i--
		}
		if j >= 0 {
			digit2 = b.toInteger(string(num2[j]))
			j--
		}
		sum := digit1 + digit2 + carry
		carry = sum / a.format.base
		result = a.toNumber(sum%a.format.base) + result
	}

	if hasDecimal {
		result = insertDecimalPoint(result, decPlaces, a.format.decimalSeparator)
	}

	return NewBigNumber(strings.TrimLeft(result, "0"))
}

// Sub subtracts another BigNumber from the BigNumber.
func (a *BigNumber) Sub(b *BigNumber) *BigNumber {
	num1, num2 := normalizeDecimal(a.value, b.value, a.format.decimalSeparator)
	i, j := len(num1)-1, len(num2)-1
	borrow := 0
	result := ""

	decPlaces := 0
	hasDecimal := strings.Contains(a.value, a.format.decimalSeparator) || strings.Contains(b.value, b.format.decimalSeparator)
	if hasDecimal {
		_, decPart1 := splitDecimal(a.value, a.format.decimalSeparator)
		_, decPart2 := splitDecimal(b.value, b.format.decimalSeparator)
		decPlaces = max(len(decPart1), len(decPart2))
	}

	for i >= 0 || j >= 0 {
		var digit1, digit2 int
		if i >= 0 {
			if string(num1[i]) == a.format.decimalSeparator {
				i--
				j--
				continue
			}
			digit1 = a.toInteger(string(num1[i]))
			i--
		}
		if j >= 0 {
			digit2 = a.toInteger(string(num2[j]))
			j--
		}
		digit1 -= borrow
		if digit1 < digit2 {
			digit1 += a.format.base
			borrow = 1
		} else {
			borrow = 0
		}
		result = a.toNumber(digit1-digit2) + result
	}

	// Trim leading zeros
	result = strings.TrimLeft(result, "0")

	if hasDecimal {
		result = insertDecimalPoint(result, decPlaces, a.format.decimalSeparator)
	}

	// Trim trailing zeros after decimal point
	if hasDecimal {
		result = strings.TrimRight(result, "0")
		if strings.HasSuffix(result, a.format.decimalSeparator) {
			result = strings.TrimSuffix(result, a.format.decimalSeparator)
		}
	}

	if result == "" {
		result = "0"
	}

	return NewBigNumber(result)
}

// Multiply multiplies two BigNumber instances and returns a new BigNumber instance.
func (a *BigNumber) Multiply(b *BigNumber) *BigNumber {
    num1, num2 := a.value, b.value
    if strings.Contains(num1, a.format.decimalSeparator) || strings.Contains(num2, a.format.decimalSeparator) {
        num1, num2 = normalizeDecimal(num1, num2, a.format.decimalSeparator)
    }

    intPart1, decPart1 := splitDecimal(num1, a.format.decimalSeparator)
    intPart2, decPart2 := splitDecimal(num2, a.format.decimalSeparator)

    intPart1 += decPart1
    intPart2 += decPart2

    m := len(decPart1) + len(decPart2)
    len1 := len(intPart1)
    len2 := len(intPart2)

    result := make([]int, len1+len2)
    for i := len1 - 1; i >= 0; i-- {
        for j := len2 - 1; j >= 0; j-- {
            mul := a.toInteger(string(intPart1[i])) * a.toInteger(string(intPart2[j]))
            p1 := i + j
            p2 := i + j + 1
            sum := mul + result[p2]
            result[p1] += sum / a.format.base
            result[p2] = sum % a.format.base
        }
    }

    resultStr := ""
    for _, digit := range result {
        resultStr += a.toNumber(digit)
    }

    if m > 0 {
        resultStr = insertDecimalPoint(resultStr, m, a.format.decimalSeparator)
    }

    // 修正前导零和小数点处理逻辑
    resultStr = strings.TrimLeft(resultStr, "0")
    if strings.HasPrefix(resultStr, a.format.decimalSeparator) {
        resultStr = "0" + resultStr
    }
    if strings.Contains(resultStr, a.format.decimalSeparator) {
        parts := strings.Split(resultStr, a.format.decimalSeparator)
        if len(parts[1]) < a.format.maxDecimal {
            parts[1] = PadEndString(parts[1], a.format.maxDecimal, "0")
        }
        resultStr = strings.Join(parts, a.format.decimalSeparator)
    }

    if a.format.maxDecimal <= 0 {
        resultStr = strings.TrimRight(resultStr, "0")
    }
    if strings.HasSuffix(resultStr, a.format.decimalSeparator) {
        resultStr = strings.TrimSuffix(resultStr, a.format.decimalSeparator)
    }
    if resultStr == "" {
        resultStr = "0"
    }

    resultBigNumber := NewBigNumber(resultStr)
    resultBigNumber.sign = a.sign * b.sign
    return resultBigNumber
}

// Divide divides two BigNumber instances and returns the quotient and remainder.
func (a *BigNumber) Divide(b *BigNumber) *BigNumber {
	if b.value == "0" {
        panic("division by zero")
    }

    num1, num2 := a.value, b.value
    intPart1, decPart1 := splitDecimal(num1, a.format.decimalSeparator)
    intPart2, decPart2 := splitDecimal(num2, a.format.decimalSeparator)

    lenDec1 := len(decPart1)
    lenDec2 := len(decPart2)

    num1 = intPart1 + decPart1
    num2 = intPart2 + decPart2

    shift := lenDec2 - lenDec1
    if shift > 0 {
        num1 += strings.Repeat("0", shift)
    } else if shift < 0 {
        num2 += strings.Repeat("0", -shift)
    }

    result := ""
    remainder := 0
    for i := 0; i < len(num1)+a.format.maxDecimal; i++ {
        var digit1 int
        if i < len(num1) {
            digit1 = a.toInteger(string(num1[i]))
        }
        remainder = remainder*a.format.base + digit1
        digitResult := remainder / a.toInteger(string(num2[0]))
        remainder = remainder % a.toInteger(string(num2[0]))
        result += a.toNumber(digitResult)
    }

    result = insertDecimalPoint(result, a.format.maxDecimal, a.format.decimalSeparator)
    integerPart,decimalPart := splitDecimal(result,a.format.decimalSeparator)

    result = strings.Join([]string{integerPart,PadEndString(decimalPart,a.format.maxDecimal,"0")},a.format.decimalSeparator)

    if(a.format.maxDecimal <= 0){
        result = strings.TrimRight(result, "0")
    }
    if strings.HasSuffix(result, a.format.decimalSeparator) {
        result = strings.TrimSuffix(result, a.format.decimalSeparator)
    }

    resultBigNumber := NewBigNumber(result)
    resultBigNumber.sign = a.sign * b.sign
    return resultBigNumber
}

// ExponentiatedBy raises the BigNumber to the power of an integer exponent.
func (a *BigNumber) ExponentiatedBy(exponent int) *BigNumber {
	result := NewBigNumber("1")
	for i := 0; i < exponent; i++ {
		result = result.Multiply(a)
	}
	return result
}

func (a *BigNumber) From(other string) *BigNumber {
	return NewBigNumber(other)
}

// IntegerValue returns the integer part of the BigNumber.
func (a *BigNumber) IntegerValue() *BigNumber {
	parts := strings.Split(a.value, a.format.decimalSeparator)
	return NewBigNumber(parts[0])
}

// IsEqualTo checks if the BigNumber is equal to another BigNumber.
func (a *BigNumber) IsEqualTo(b *BigNumber) bool {
	return a.ComparedTo(b) == 0
}

// IsFinite checks if the BigNumber is finite (not infinite or NaN).
func (a *BigNumber) IsFinite() bool {
	return true // For simplicity, as we're not handling infinities or NaN.
}

// IsGreaterThan checks if the BigNumber is greater than another BigNumber.
func (a *BigNumber) IsGreaterThan(b *BigNumber) bool {
	return a.ComparedTo(b) > 0
}

// IsGreaterThanOrEqualTo checks if the BigNumber is greater than or equal to another BigNumber.
func (a *BigNumber) IsGreaterThanOrEqualTo(b *BigNumber) bool {
	return a.ComparedTo(b) >= 0
}

// IsInteger checks if the BigNumber is an integer.
func (a *BigNumber) IsInteger() bool {
	return a.DecimalPlaces() == 0
}

// IsLessThan checks if the BigNumber is less than another BigNumber.
func (a *BigNumber) IsLessThan(b *BigNumber) bool {
	return a.ComparedTo(b) < 0
}

// IsLessThanOrEqualTo checks if the BigNumber is less than or equal to another BigNumber.
func (a *BigNumber) IsLessThanOrEqualTo(b *BigNumber) bool {
	return a.ComparedTo(b) <= 0
}

func (a *BigNumber) Sum(values ...any) *BigNumber {
	total := a
	for _, val := range values {
        var num *BigNumber
        switch val.(type) {
        case string:
            num = NewBigNumber(val.(string))
        case BigNumber:
            num = val.(*BigNumber)
        default:
            panic("type error")
        }
		total = total.Plus(num)
	}
	return total
}


// Minus subtracts another BigNumber from the BigNumber.
func (a *BigNumber) Minus(b *BigNumber) *BigNumber {
	if a.sign == b.sign {
		if a.IsGreaterThan(b) {
			result := a.Sub(b)
			result.sign = a.sign
			return result
		}
		result := b.Sub(a)
		result.sign = -a.sign
		return result
	}
	result := a.add(b)
	result.sign = a.sign
	return result
}

// Mod calculates the modulus of two BigNumber instances.
func (a *BigNumber) Mod(b *BigNumber) *BigNumber {
	if b.value == "0" {
		panic("division by zero")
	}

	_, remainder := a.divideStrings(a.value, b.value)
	return remainder
}

// MultipliedBy multiplies the BigNumber by another BigNumber.
func (a *BigNumber) MultipliedBy(b *BigNumber) *BigNumber {
	return a.Multiply(b)
}

// Plus adds another BigNumber to the BigNumber.
func (a *BigNumber) Plus(b *BigNumber) *BigNumber {
	if a.sign == b.sign {
		result := a.add(b)
		result.sign = a.sign
		return result
	}
	if a.IsGreaterThan(b) {
		result := a.Sub(b)
		result.sign = a.sign
		return result
	}
	result := b.Sub(a)
	result.sign = b.sign
	return result
}

// Precision returns the precision of the BigNumber.
func (a *BigNumber) Precision() int {
	return len(strings.ReplaceAll(a.value, a.format.decimalSeparator, ""))
}

// ShiftedBy shifts the decimal point by a given number of places.
func (a *BigNumber) ShiftedBy(places int) *BigNumber {
	if places == 0 {
		return a
	}

	parts := strings.Split(a.value, a.format.decimalSeparator)
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	if places > 0 {
		for places > len(decPart) {
			decPart += "0"
		}
		return NewBigNumber(intPart + decPart[:places] + a.format.decimalSeparator + decPart[places:])
	} else {
		places = -places
		for places > len(intPart) {
			intPart = "0" + intPart
		}
		return NewBigNumber(intPart[:len(intPart)-places] + a.format.decimalSeparator + intPart[len(intPart)-places:] + decPart)
	}
}

// String returns the string representation of a BigNumber.
func (a *BigNumber) String() string {
	if a.sign == -1 && a.value != "0" {
		return "-" + a.value
	}
	return a.value
}

// compareIntegerParts compares the integer parts of two numbers.
func compareIntegerParts(intPart1, intPart2 string) int {
	if len(intPart1) > len(intPart2) {
		return 1
	} else if len(intPart1) < len(intPart2) {
		return -1
	}
	for i := 0; i < len(intPart1); i++ {
		if intPart1[i] > intPart2[i] {
			return 1
		} else if intPart1[i] < intPart2[i] {
			return -1
		}
	}
	return 0
}

// splitDecimal splits a number into integer and decimal parts
func splitDecimal(value string, split string) (string, string) {
	parts := strings.Split(value, split)
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}
	return intPart, decPart
}

// insertDecimalPoint inserts the decimal point into the number string at the correct position.
func insertDecimalPoint(num string, decPlaces int, decimalSeparator string) string {
	if decPlaces <= 0 {
		return num
	}

	if decPlaces >= len(num) {
		num = strings.Repeat("0", decPlaces-len(num)+1) + num
	}
	return num[:len(num)-decPlaces] + decimalSeparator + num[len(num)-decPlaces:]
}

func normalizeDecimal(num1, num2, decimalSeparator string) (string, string) {
	intPart1, decPart1 := splitDecimal(num1, decimalSeparator)
	intPart2, decPart2 := splitDecimal(num2, decimalSeparator)

	if len(decPart1) > len(decPart2) {
		decPart2 += strings.Repeat("0", len(decPart1)-len(decPart2))
	} else {
		decPart1 += strings.Repeat("0", len(decPart2)-len(decPart1))
	}

	return intPart1 + decimalSeparator + decPart1, intPart2 + decimalSeparator + decPart2
}
