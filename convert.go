package utils

import (
	"errors"
	"math"
	"strings"
)

type ConversionFunctions struct{}

func (c *ConversionFunctions) Length(value float64, from, to string) float64 {
	units := map[string]float64{
		"m":  100,
		"km": 1000,
		"cm": 1,
		"mm": 0.01,
		"um": 1e-3,
		"nm": 1e-6,
		"dm": 10,
	}
	_, fromValid := units[from]
	_, toValid := units[to]
	if fromValid && toValid {
		return (value * units[from]) / units[to]
	} else {
		return -1
	}
}

func (c *ConversionFunctions) Weight(value float64, from, to string) (float64, error) {
	units := map[string]float64{
		"g":  1,
		"kg": 1000,
		"mg": 0.001,
		"t":  1000000,
	}

	if _, exists := units[from]; !exists {
		return 0, errors.New("Invalid weight unit")
	}
	if _, exists := units[to]; !exists {
		return 0, errors.New("Invalid weight unit")
	}

	return (value * units[from]) / units[to], nil
}

func (c *ConversionFunctions) Size(value float64, from, to string) float64 {
	units := map[string]float64{
		"k": 1000,
		"w": 10000,
		"m": 1000000,
	}
	_, fromValid := units[from]
	_, toValid := units[to]
	if fromValid && toValid {

		return (value * units[from]) / units[to]
	} else {
		return -1
	}
}

func (c *ConversionFunctions) Volume(value float64, from, to string) float64 {
	units := map[string]float64{
		"ml": 1,
		"cl": 10,
		"dl": 100,
		"l":  1000,
	}
	_, fromValid := units[from]
	_, toValid := units[to]
	if fromValid && toValid {
		return (value * units[from]) / units[to]
	} else {
		return -1
	}
}

func (c *ConversionFunctions) Storage(value float64, from, to string) (float64, error) {
	units := map[string]int{
		"B":  0,
		"KB": 1,
		"MB": 2,
		"GB": 3,
		"TB": 4,
		"PB": 5,
	}
	fromIndex, existsFrom := units[from]
	toIndex, existsTo := units[to]

	if !existsFrom || !existsTo {
		return 0, errors.New("Invalid storage unit")
	}

	diff := toIndex - fromIndex

	if diff == 0 {
		return value, nil
	}

	factor := math.Pow(1024, math.Abs(float64(diff)))
	if diff > 0 {
		return value / factor, nil
	}
	return value * factor, nil
}

func (c *ConversionFunctions) NetSpeed(speed float64, from, to string) float64 {
	networkUnitTable := map[string]float64{
		"bps":  1,
		"Kbps": 1000,
		"Mbps": 1000000,
		"Gbps": 1000000000,
	}

	fromFactor, existsFrom := networkUnitTable[from]
	toFactor, existsTo := networkUnitTable[to]

	if !existsFrom || !existsTo {
		return -1
	}

	return (speed * fromFactor) / toFactor
}

func (c *ConversionFunctions) WeightEN(weight float64, from, to string) float64 {
	unitTable := map[string]float64{
		"pounds": 1,
		"grams":  453.592,           // 1 磅 ≈ 453.592 克
		"kgs":    453.592 / 1000,    // 1 磅 ≈ 0.453592 千克
		"tons":   453.592 / 1000000, // 1 磅 ≈ 0.000453592 吨
		"ans":    28.349523125,
		"kg":     1000,    // 1 千克 = 1000 克
		"ton":    1000000, // 1 吨 = 1,000,000 克
	}
	fromFactor, existsFrom := unitTable[from]
	toFactor, existsTo := unitTable[to]

	if !existsFrom || !existsTo {
		return -1
	}

	// 如果 fromUnit 和 toUnit 相同，则直接返回 weight
	if from == to {
		return weight
	}

	return (weight * fromFactor) / toFactor
}

func indexOf(slice []string, ele string) int {
	for k, v := range slice {
		if v == ele {
			return k
		}
	}
	return -1
}

func stringToNumberCallback(digits string,base int) int {
    slice := strings.Split("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-$","")
    index := indexOf(slice,digits)
    if index > base {
        return -1
    }
    return index
}

func numberToStringCallback(remainder int, unit int) (string, string, int) {
    slice := strings.Split("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-$","")
    digits := slice[remainder]
    _unit := ""
    return digits,_unit,len(slice)
}

func NumberToString(numberic float64, base int, callback func(remainder int, unit int) (string, string, int)) string {
	if callback == nil {
		callback = numberToStringCallback
	}
	if base == 0 {
		base = 10
	}
	quotient := math.Abs(numberic)
	result := ""
	unitIndex := 0
	_, _, length := callback(0, 0)
	if base > length {
		return result
	}
	for quotient > 0 {
		remainder := int(quotient) % base
		digits, unit, _ := callback(remainder, unitIndex)
		result = digits + unit + result
		quotient = math.Floor(float64(quotient) / float64(base))
		unitIndex++
	}
	if numberic < 0 {
		return "-" + result
	}
	return result

}

func StringToNumber(numberic string, base int, callback func(digits string,base int) int) int {
	if callback == nil {
		callback = stringToNumberCallback
	}
	if base == 0 {
		base = 10
	}
	result := 0
	power := 1
	nums := strings.Split(numberic, "")
	for i := len(nums); i > 0; i-- {
		val := callback(nums[i-1],base)
        if val > base {
            return -0
        }
		result += (val * power)
		power *= base

	}
	return result
}
