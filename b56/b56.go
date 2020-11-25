package b56

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var b56UnitFloor   [12]uint64
var b56UnitCeiling [12]uint64
var b56toB10Values = map[string] string{
	"0": "0",
	"1": "1",
	"2": "2",
	"3": "3",
	"4": "4",
	"5": "5",
	"6": "6",
	"7": "7",
	"8": "8",
	"9": "9",
	"a": "10",
	"b": "11",
	"c": "12",
	"d": "13",
	"e": "14",
	"f": "15",
	"g": "16",
	"h": "17",
	"j": "18",
	"k": "19",
	"m": "20",
	"n": "21",
	"p": "22",
	"q": "23",
	"r": "24",
	"s": "25",
	"t": "26",
	"u": "27",
	"v": "28",
	"w": "29",
	"x": "30",
	"y": "31",
	"z": "32",
	"A": "33",
	"B": "34",
	"C": "35",
	"D": "36",
	"E": "37",
	"F": "38",
	"G": "39",
	"H": "40",
	"J": "41",
	"K": "42",
	"M": "43",
	"N": "44",
	"P": "45",
	"Q": "46",
	"R": "47",
	"S": "48",
	"T": "49",
	"U": "50",
	"V": "51",
	"W": "52",
	"X": "53",
	"Y": "54",
	"Z": "55",
}
var b10toB56Values = map[string] string{
	"0":  "0",
	"1":  "1",
	"2":  "2",
	"3":  "3",
	"4":  "4",
	"5":  "5",
	"6":  "6",
	"7":  "7",
	"8":  "8",
	"9":  "9",
	"10": "a",
	"11": "b",
	"12": "c",
	"13": "d",
	"14": "e",
	"15": "f",
	"16": "g",
	"17": "h",
	"18": "j",
	"19": "k",
	"20": "m",
	"21": "n",
	"22": "p",
	"23": "q",
	"24": "r",
	"25": "s",
	"26": "t",
	"27": "u",
	"28": "v",
	"29": "w",
	"30": "x",
	"31": "y",
	"32": "z",
	"33": "A",
	"34": "B",
	"35": "C",
	"36": "D",
	"37": "E",
	"38": "F",
	"39": "G",
	"40": "H",
	"41": "J",
	"42": "K",
	"43": "M",
	"44": "N",
	"45": "P",
	"46": "Q",
	"47": "R",
	"48": "S",
	"49": "T",
	"50": "U",
	"51": "V",
	"52": "W",
	"53": "X",
	"54": "Y",
	"55": "Z",
}
var debug = false
var debugLevel = 0

func Base56Encode (base10 uint64, zeroPadding bool) string {

	// we know that b56 can store the max uint64 value in 12 units or less: initialize the globals
	for unit := int(0); unit <= 11; unit++ {
		b56UnitFloor[unit] = uint64(math.Pow(56, float64(unit)))
		if unit != 11 {
			//unit == 11 will be out of range for a uint64
			b56UnitCeiling[unit] = b56UnitFloor[unit]*55
		} else {
			//this should be uint64 limit + 1
			b56UnitCeiling[unit] = 18446744073709551615
		}
		if (debug && debugLevel >= 1) {
			fmt.Printf("Unit: %v, Floor: %v, Ceiling: %v\n", unit, b56UnitFloor[unit], b56UnitCeiling[unit])
		}

	}

	// there are other ways to handle the 0 edge case, favoring brevity
	b56FinalValue  := [...]string{"0","0","0","0","0","0","0","0","0","0","0","0"}

	baseTenRemainder := base10
	if base10 != 0 {
		for arrayPosition := 11; arrayPosition >= 0; arrayPosition-- {
			/*
			 * arrayPosition / digitPosition is not intuitive at first.  The array reads left to right, but
			 * the number itself reads right to left in terms of position.  This swaps them so the number is correct
			 * even though the array is backwards.
			 */
			digitPosition := 11 - arrayPosition
			b56FinalValue[digitPosition], baseTenRemainder = b56Digit(uint64(arrayPosition), baseTenRemainder)
			if debug && debugLevel >= 2 {
				fmt.Printf("Digit Position: %v, Remainder: %v, b56Value: %v\n", arrayPosition, baseTenRemainder, b56FinalValue)
			}
		}
	}
	if zeroPadding == false {
		// if the zeroPadding flag is set to false, remove leading zeros
		previousDigit := "0"
		foundNonZero := false
		for digitPosition := 0; digitPosition <= 11; digitPosition++ {
			curDigit := b56FinalValue[digitPosition]
			if debug && debugLevel >= 1 {
				fmt.Printf("Previous Digit: %v, Current Digit: %v\n", previousDigit, curDigit)
			}
			if previousDigit == "0" && digitPosition != 0 && !foundNonZero {
				b56FinalValue[digitPosition - 1] = ""
			}
			if curDigit != "0" {
				foundNonZero = true
			}
			previousDigit = curDigit
		}
	}

	b56ReturnValue := strings.Join(b56FinalValue[:],"")
	return b56ReturnValue
}

func Base56EncodeWithDomainPrefix (base10 uint64, domainPrefix string) string {
	//always padded if there's a domain prefix
	zeroPadding := false
	base56 := Base56Encode(base10, zeroPadding)
	base56WithDomainPrefix := domainPrefix + "-" + base56
	return base56WithDomainPrefix
}

func Base10Encode (base56 string) uint64 {
	base56Length := len(base56)
	if base56Length < 12 {
		difference := uint64(12 - base56Length)
		var padding string
		for pad := uint64(0); pad < difference; pad++ {
			padding += "0"
		}
		base56 = padding + base56
		if debug && debugLevel >= 2 {
			fmt.Printf("Base10 needed padded with %v digits resulting in %v, with a length of %v\n", difference, base56, len(base56))
		}
	}
	b10FinalValue := uint64(0)
	b56AsUnit := strings.Split(base56, "")
	for digitPosition := 11; digitPosition >= 0; digitPosition-- {
		digitPositionAdjusted := 11 - digitPosition
		// the unit value of the b56, as a string
		currentB56Character := b56AsUnit[digitPosition]
		// this is the decimal representation, also a string
		currentB10Character := b56toB10Values[currentB56Character]
		currentB10Decimal, _ := strconv.ParseUint(currentB10Character, 10, 64)
		if debug && debugLevel >= 2 {
			fmt.Printf("current decimal value is: %v\n", currentB10Decimal)
		}
		computedDecimalValue := math.Pow(56, float64(digitPositionAdjusted))*float64(currentB10Decimal)
		b10FinalValue = b10FinalValue + uint64(computedDecimalValue)
		if debug && debugLevel >= 2 {
			fmt.Printf("b10FinalValue position[%v]: %v for decimal value: %v\n", digitPositionAdjusted, uint64(computedDecimalValue), currentB10Decimal)
		}
	}

	return b10FinalValue
}

func Base10EncodeWithDomainPrefix (base56 string) uint64 {
	base56Parts := strings.Split(base56, "-")
	base56NoDomain := base56Parts[1]
	b10FinalValue := Base10Encode(base56NoDomain)
	return b10FinalValue
}

func b56Digit (digitPosition uint64, baseTen uint64) (string, uint64) {

	b56UnitValue := "0"
	baseTenRemainder := baseTen

	if debug && debugLevel >= 2 {
		fmt.Printf("Digit Position %v\n", digitPosition)
	}
	if baseTen >= b56UnitFloor[digitPosition] {
		if digitPosition == 11 {
			// there is only one value that can satisfy in digit 11 without breaking uint64 limits
			b56UnitValue = b10toB56Values["1"]
			baseTenRemainder = baseTen - b56UnitFloor[digitPosition]
		} else {
			for valueSelection := uint64(55); valueSelection > 0; valueSelection-- {
				currentCeiling := b56UnitFloor[digitPosition] * valueSelection
				if debug && debugLevel >= 2 {
					fmt.Printf("Ceiling of %v for value %v\n", currentCeiling, baseTen)
				}
				if baseTen >= currentCeiling {
					baseTenRemainder = baseTen - currentCeiling
					b56UnitDecimal := strconv.FormatUint(valueSelection, 10)
					b56UnitValue = b10toB56Values[b56UnitDecimal]
					break
				}
			}
		}

	}

	if debug && debugLevel >= 2 {
		fmt.Printf("remainder %v, unit value %v\n", baseTenRemainder, b56UnitValue)
	}
	return b56UnitValue, baseTenRemainder
}