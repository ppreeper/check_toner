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

func tonerOutput(color string, maxValue string, lvlValue string) string {
	color = strings.ToUpper(color)
	max, errm := strconv.Atoi(maxValue)
	lvl, errl := strconv.Atoi(lvlValue)
	var output string

	if errm == nil && errl == nil {
		if max != 0 {
			level := 100 * float64(float64(lvl)/float64(max))
			tLevels := "-- Toner Levels: " + lvlValue + " of " + maxValue + " Remaining | "
			tLevels += strconv.FormatFloat(level, 'f', 0, 64) + "\n"
			if level >= 10.0 {
				output = color + " Toner OK " + tLevels
			} else {
				output = color + " Toner LOW " + tLevels
			}
		}
	}
	return output
}

func tonerLevel(color string, brand string) string {
	color = strings.ToUpper(color)
	var tonerColor string
	var t toners
	var output string

	switch brand {
	case "HP":
		t = *hp
	case "UTAX":
		t = *utax
	case "KYOCERA":
		t = *utax
	}

	if color == "C" {
		tonerColor = "CYAN"
		output = tonerOutput(tonerColor, getSNMPValue(t.cyanMax), getSNMPValue(t.cyanLvl))
	}
	if color == "M" {
		tonerColor = "MAGENTA"
		output = tonerOutput(tonerColor, getSNMPValue(t.magentaMax), getSNMPValue(t.magentaLvl))
	}
	if color == "Y" {
		tonerColor = "YELLOW"
		output = tonerOutput(tonerColor, getSNMPValue(t.yellowMax), getSNMPValue(t.yellowLvl))
	}
	if color == "K" {
		tonerColor = "BLACK"
		output = tonerOutput(tonerColor, getSNMPValue(t.kromaMax), getSNMPValue(t.kromaLvl))
	}
	return output
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

// main function
func main() {
	flag.Parse()
	if *host == "" {
		fmt.Println("Host not set")
	} else {
		//fmt.Printf("Host:\t%s\nColor:\t%s\nBrand:\t%s\n", *host, *color, *brand)
		r := tonerLevel(*color, *brand)
		fmt.Printf("%v", r)
	}
	return
}
