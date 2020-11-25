package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"
	"b56"
)

var b56UnitFloor   [12]uint64
var b56UnitCeiling [12]uint64
var b56toB10Values = make(map[string] string)
var b10toB56Values = make(map[string] string)
var debug bool
var debugLevel int
var valueStart uint64
var valueEnd uint64
var performanceMetric bool
var zeroPadding bool
var domainPrefix string

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable Debugging")
	flag.BoolVar(&performanceMetric, "performanceMetric", false, "Performance Metric")
	flag.BoolVar(&zeroPadding, "zeroPadding", false, "Enable Padded b56 Values; this flag has no effect when using a prefix.")
	flag.IntVar(&debugLevel, "debugLevel", 1, "Debug Verbose Level")
	flag.Uint64Var(&valueStart, "valueStart", 0, "Starting Value")
	flag.Uint64Var(&valueEnd, "valueEnd", 100000, "Ending Value")
	flag.StringVar(&domainPrefix, "domainPrefix", "", "Append a specific prefix to a b56 value.")
	flag.Parse()

	// the largest number uint64 can handle: 18446744073709551614 == 14PSsSsyWste
	start :=time.Now()
	for base10 := uint64(valueStart); base10 <= valueEnd; base10++ {
		var curBase56 string
		if domainPrefix != "" {
			//values with domain prefix are always zeroPadding = false
			curBase56 = b56.Base56EncodeWithDomainPrefix(base10, domainPrefix)
		} else {
			curBase56 = b56.Base56Encode(base10, zeroPadding)
		}

		fmt.Printf("%v => %v\n", base10, curBase56)
		var curBase10 uint64
		if domainPrefix != "" {
			curBase10 = b56.Base10EncodeWithDomainPrefix(curBase56)
		} else {
			curBase10 = b56.Base10Encode(curBase56)
		}

		fmt.Printf("%v <= %v\n", curBase10, curBase56)
		checkResult := base10 - curBase10
		if (debug && debugLevel >= 1) || checkResult > 0 {
			fmt.Printf("difference in b10 is [%v]\n", base10-curBase10)
			if checkResult > 0 {
				os.Exit(0)
			}
		}
	}
	stop := time.Now()
	elapsed := stop.Sub(start)
	total := valueEnd - valueStart
	rate := math.Round(float64(total) / float64(elapsed.Seconds()))
	if ((debug && debugLevel >= 1) && total > 9999) || performanceMetric {
		fmt.Printf("Generated %v Base56 IDs in %v at a rate of %v per second\n", total, elapsed, rate)
	}
}