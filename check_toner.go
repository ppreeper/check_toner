package main

import (
	//"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	g "github.com/soniah/gosnmp"
	//"time"
)

var host = flag.String("H", "", "Printer to query")
var color = flag.String("C", "K", "Toner Color")
var brand = flag.String("B", "HP", "Printer Brand")

// toners struct
type toners struct {
	cyanMax    string
	cyanLvl    string
	magentaMax string
	magentaLvl string
	yellowMax  string
	yellowLvl  string
	kromaMax   string
	kromaLvl   string
}

var hp = &toners{
	".1.3.6.1.2.1.43.11.1.1.8.1.2", ".1.3.6.1.2.1.43.11.1.1.9.1.2",
	".1.3.6.1.2.1.43.11.1.1.8.1.3", ".1.3.6.1.2.1.43.11.1.1.9.1.3",
	".1.3.6.1.2.1.43.11.1.1.8.1.4", ".1.3.6.1.2.1.43.11.1.1.9.1.4",
	".1.3.6.1.2.1.43.11.1.1.8.1.1", ".1.3.6.1.2.1.43.11.1.1.9.1.1",
}

var utax = &toners{
	".1.3.6.1.2.1.43.11.1.1.8.1.1", ".1.3.6.1.2.1.43.11.1.1.9.1.1",
	".1.3.6.1.2.1.43.11.1.1.8.1.2", ".1.3.6.1.2.1.43.11.1.1.9.1.2",
	".1.3.6.1.2.1.43.11.1.1.8.1.3", ".1.3.6.1.2.1.43.11.1.1.9.1.3",
	".1.3.6.1.2.1.43.11.1.1.8.1.4", ".1.3.6.1.2.1.43.11.1.1.9.1.4",
}

func tonerOutput(color string, maxValue string, lvlValue string) {
	max, errm := strconv.Atoi(maxValue)
	lvl, errl := strconv.Atoi(lvlValue)
	// var output string
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
	var t toners

	switch brand {
	case "HP":
		t = *hp
	case "UTAX":
		t = *utax
	case "KYOCERA":
		t = *utax
	}

	fmt.Printf("%v", t.cyanMax)

	if color == "C" {
		tonerColor = "CYAN"
		tonerOutput(tonerColor, getSNMPValue(t.cyanMax), getSNMPValue(t.cyanLvl))
	}
	if color == "M" {
		tonerColor = "MAGENTA"
		tonerOutput(tonerColor, getSNMPValue(t.magentaMax), getSNMPValue(t.magentaLvl))
	}
	if color == "Y" {
		tonerColor = "YELLOW"
		tonerOutput(tonerColor, getSNMPValue(t.yellowMax), getSNMPValue(t.yellowLvl))
	}
	if color == "K" {
		tonerColor = "BLACK"
		tonerOutput(tonerColor, getSNMPValue(t.kromaMax), getSNMPValue(t.kromaLvl))
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
	} else {
		//fmt.Printf("Host:\t%s\nColor:\t%s\nBrand:\t%s\n", *host, *color, *brand)
		tonerLevel(*color, *brand)
	}
	return
}
