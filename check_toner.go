package main

import (
	//"errors"
	"flag"
	"fmt"
	g "github.com/soniah/gosnmp"
	"log"
	"strconv"
	"strings"
	//"time"
)

var host = flag.String("H", "", "Printer to query")
var color = flag.String("C", "K", "Toner Color")
var brand = flag.String("B", "HP", "Printer Brand")

func tonerOutput(color string, maxValue string, lvlValue string) {
	max, errm := strconv.Atoi(maxValue)
	lvl, errl := strconv.Atoi(lvlValue)
	if errm == nil && errl == nil {
		if max != 0 {
			level := 100 * float64(float64(lvl)/float64(max))
			if level >= 10.0 {
				fmt.Printf("%s Toner OK -- Toner Levels: %s of %s Remaining | %s\n", color, lvlValue, maxValue, lvlValue)
			} else {
				fmt.Printf("%s Toner LOW -- Toner Levels: %s of %s Remaining | %s\n", color, lvlValue, maxValue, lvlValue)
			}
		}
	}
	return
}
func tonerLevel(color string, brand string) {
    color = strings.ToUpper(color)
	var tonerColor string
	if brand == "HP" {
		var C_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.2" // Max Cyan Toner
		var C_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.2" // Remaining Cyan Toner
		var M_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.3" // Max Magenta Toner
		var M_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.3" // Remaining Magenta Toner
		var Y_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.4" // Max Yellow Toner
		var Y_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.4" // Remaining Yellow Toner
		var K_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.1" // Max Black Toner
		var K_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.1" // Remaining Black Toner

		if color == "C" {
			tonerColor = "CYAN"
			tonerOutput(tonerColor, getSNMPValue(C_MAX), getSNMPValue(C_LVL))
		}
		if color == "M" {
			tonerColor = "MAGENTA"
			tonerOutput(tonerColor, getSNMPValue(M_MAX), getSNMPValue(M_LVL))
		}
		if color == "Y" {
			tonerColor = "YELLOW"
			tonerOutput(tonerColor, getSNMPValue(Y_MAX), getSNMPValue(Y_LVL))
		}
		if color == "K" {
			tonerColor = "BLACK"
			tonerOutput(tonerColor, getSNMPValue(K_MAX), getSNMPValue(K_LVL))
		}
	}
	if brand == "UTAX" || brand == "KYOCERA" {
		var C_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.1" // Max Cyan Toner
		var C_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.1" // Remaining Cyan Toner
		var M_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.2" // Max Magenta Toner
		var M_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.2" // Remaining Magenta Toner
		var Y_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.3" // Max Yellow Toner
		var Y_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.3" // Remaining Yellow Toner
		var K_MAX string = ".1.3.6.1.2.1.43.11.1.1.8.1.4" // Max Black Toner
		var K_LVL string = ".1.3.6.1.2.1.43.11.1.1.9.1.4" // Remaining Black Toner

		if color == "C" {
			tonerColor = "CYAN"
			tonerOutput(tonerColor, getSNMPValue(C_MAX), getSNMPValue(C_LVL))
		}
		if color == "M" {
			tonerColor = "MAGENTA"
			tonerOutput(tonerColor, getSNMPValue(M_MAX), getSNMPValue(M_LVL))
		}
		if color == "Y" {
			tonerColor = "YELLOW"
			tonerOutput(tonerColor, getSNMPValue(Y_MAX), getSNMPValue(Y_LVL))
		}
		if color == "K" {
			tonerColor = "BLACK"
			tonerOutput(tonerColor, getSNMPValue(K_MAX), getSNMPValue(K_LVL))
		}
	}
	return
}

func getSNMPValue(oid string) string {
	g.Default.Target = *host
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()
	oids := []string{oid}
	result, err2 := g.Default.Get(oids)
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	return fmt.Sprintf("%s", g.ToBigInt(result.Variables[0].Value))
}

func main() {
	// All the interesting pieces are with the variables declared above, but
	// to enable the flag package to see the flags defined there, one must
	// execute, typically at the start of main (not init!):
	flag.Parse()
	// We don't run it here because this is not a main function and
	// the testing suite has already parsed the flags.
	//var flg = ""
	if *host == "" {
		fmt.Println("Host not set")
		return
	} else {
		//fmt.Printf("Host:\t%s\nColor:\t%s\nBrand:\t%s\n", *host, *color, *brand)
		tonerLevel(*color, *brand)
	}
	return
}
