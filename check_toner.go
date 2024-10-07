package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	g "github.com/gosnmp/gosnmp"
)

type toner struct {
	Max   string
	Level string
}

type toners struct {
	Cyan    toner
	Magenta toner
	Yellow  toner
	Kroma   toner
}

var brands = map[string]toners{
	"HP": {
		Cyan:    toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.2", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.2"},
		Magenta: toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.3", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.3"},
		Yellow:  toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.4", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.4"},
		Kroma:   toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.1", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.1"},
	},
	"KYOCERA": {
		Cyan:    toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.1", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.1"},
		Magenta: toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.2", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.2"},
		Yellow:  toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.3", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.3"},
		Kroma:   toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.4", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.4"},
	},
	"UTAX": {
		Cyan:    toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.1", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.1"},
		Magenta: toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.2", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.2"},
		Yellow:  toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.3", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.3"},
		Kroma:   toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.4", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.4"},
	},
	"HP8710": {
		Cyan:    toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.2", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.2"},
		Magenta: toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.3", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.3"},
		Yellow:  toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.1", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.1"},
		Kroma:   toner{Max: ".1.3.6.1.2.1.43.11.1.1.8.1.4", Level: ".1.3.6.1.2.1.43.11.1.1.9.1.4"},
	},
}

func tonerOutput(color string, maxValue string, lvlValue string) string {
	color = strings.ToUpper(color)
	max, errm := strconv.Atoi(maxValue)
	lvl, errl := strconv.Atoi(lvlValue)

	if errm == nil && errl == nil {
		if max != 0 {
			level := 100 * float64(float64(lvl)/float64(max))
			ok := "OK"
			if level <= 10.0 {
				ok = "LOW"
			}
			return fmt.Sprintf("%s Toner %s -- Toner Levels: %s of %s Remaining | %d", color, ok, lvlValue, maxValue, int(level))
		}
	}
	return ""
}

func tonerLevel(host, color, brand string) string {
	color = strings.ToUpper(color)
	var tonerColor string
	var t toners
	var output string
	var max, lvl string

	t = brands[brand]

	switch color {
	case "C":
		tonerColor = "Cyan"
		max, _ = getSNMPValue(host, t.Cyan.Max)
		lvl, _ = getSNMPValue(host, t.Cyan.Level)
	case "M":
		tonerColor = "Magenta"
		max, _ = getSNMPValue(host, t.Magenta.Max)
		lvl, _ = getSNMPValue(host, t.Magenta.Level)
	case "Y":
		tonerColor = "Yellow"
		max, _ = getSNMPValue(host, t.Yellow.Max)
		lvl, _ = getSNMPValue(host, t.Yellow.Level)
	case "K":
		tonerColor = "Black"
		max, _ = getSNMPValue(host, t.Kroma.Max)
		lvl, _ = getSNMPValue(host, t.Kroma.Level)
	}

	output = tonerOutput(tonerColor, max, lvl)
	return output
}

func getSNMPValue(host, oid string) (string, error) {
	g.Default.Target = host
	err := g.Default.Connect()
	if err != nil {
		return "", fmt.Errorf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()
	oids := []string{oid}
	result, err := g.Default.Get(oids)
	if err != nil {
		return "", fmt.Errorf("Get() err: %v", err)
	}
	return fmt.Sprintf("%s", g.ToBigInt(result.Variables[0].Value)), err
}

// main function
func main() {
	var host string
	var color string
	var brand string
	flag.StringVar(&host, "H", "", "Printer to query")
	flag.StringVar(&color, "C", "K", "Toner Color")
	flag.StringVar(&brand, "B", "HP", "Printer Brand")

	flag.Parse()
	if host == "" {
		fmt.Fprintf(os.Stdout, "Host not set\n")
	} else {
		r := tonerLevel(host, color, brand)
		fmt.Fprintf(os.Stdout, "%s", r)
	}
}
