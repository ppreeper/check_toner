package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

// mulColorTests use values that you know are right
var mulColorTests = []struct {
	color    string
	maxValue string
	lvlValue string
	expected string
}{
	{"cyan", "100", "85", "CYAN Toner OK -- Toner Levels: 85 of 100 Remaining | 85\n"},
	{"cyan", "10000", "8500", "CYAN Toner OK -- Toner Levels: 8500 of 10000 Remaining | 85\n"},
	{"cyan", "10000", "900", "CYAN Toner LOW -- Toner Levels: 900 of 10000 Remaining | 9\n"},
	{"cyan", "100", "9", "CYAN Toner LOW -- Toner Levels: 9 of 100 Remaining | 9\n"},
	{"MAGENTA", "100", "85", "MAGENTA Toner OK -- Toner Levels: 85 of 100 Remaining | 85\n"},
	{"MAGENTA", "10000", "8500", "MAGENTA Toner OK -- Toner Levels: 8500 of 10000 Remaining | 85\n"},
	{"MAGENTA", "10000", "900", "MAGENTA Toner LOW -- Toner Levels: 900 of 10000 Remaining | 9\n"},
	{"MAGENTA", "100", "9", "MAGENTA Toner LOW -- Toner Levels: 9 of 100 Remaining | 9\n"},
	{"YELLOW", "100", "85", "YELLOW Toner OK -- Toner Levels: 85 of 100 Remaining | 85\n"},
	{"YELLOW", "10000", "8500", "YELLOW Toner OK -- Toner Levels: 8500 of 10000 Remaining | 85\n"},
	{"YELLOW", "10000", "900", "YELLOW Toner LOW -- Toner Levels: 900 of 10000 Remaining | 9\n"},
	{"YELLOW", "100", "9", "YELLOW Toner LOW -- Toner Levels: 9 of 100 Remaining | 9\n"},
	{"BLACK", "100", "85", "BLACK Toner OK -- Toner Levels: 85 of 100 Remaining | 85\n"},
	{"BLACK", "10000", "8500", "BLACK Toner OK -- Toner Levels: 8500 of 10000 Remaining | 85\n"},
	{"BLACK", "10000", "900", "BLACK Toner LOW -- Toner Levels: 900 of 10000 Remaining | 9\n"},
	{"BLACK", "100", "9", "BLACK Toner LOW -- Toner Levels: 9 of 100 Remaining | 9\n"},
}

// TestTonerOutput test
func TestTonerOutput(t *testing.T) {
	// var op string
	for _, mt := range mulColorTests {
		r := tonerOutput(mt.color, mt.maxValue, mt.lvlValue)
		if r != mt.expected {
			t.Errorf("%s %s", mt.color, r)
		}
	}
}

// BenchmarkPattern
func BenchmarkTonerOutput(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tonerOutput("CYAN", "900", "10000")
	}
}

// mulColorTests use values that you know are right
var mulColorBrand = []struct {
	color string
	brand string
}{
	{"C", "HP"},
	{"M", "HP"},
	{"Y", "HP"},
	{"K", "HP"},
	{"C", "UTAX"},
	{"M", "UTAX"},
	{"Y", "UTAX"},
	{"K", "UTAX"},
	{"C", "KYOCERA"},
	{"M", "KYOCERA"},
	{"Y", "KYOCERA"},
	{"K", "KYOCERA"},
}

// TestTonerLevel test
func TestTonerLevel(t *testing.T) {
	sPattern := regexp.MustCompile(`^(\w+) Toner (OK|LOW) -- Toner Levels: (\d+) of (\d+) Remaining \| (\d+)`)
	hostList := []string{"192.168.101.80", "192.168.101.51", "192.168.101.28"}
	for h := 0; h < len(hostList); h++ {
		*host = hostList[h]
		for _, mt := range mulColorBrand {
			r := tonerLevel(mt.color, mt.brand)
			sc := sPattern.FindSubmatch([]byte(r))
			if len(sc) > 0 {
				if !sPattern.Match([]byte(r)) {
					t.Errorf("TestTonerLevel output not correct\n")
				}
			}
		}
	}
}

// BenchmarkTonerLevel
func BenchmarkTonerLevel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tonerOutput("CYAN", "900", "10000")
	}
}

// mulSNMPTests use values that you know are right
var mulSNMPTests = []struct {
	host     string
	expected string
}{
	{"127.0.0.1", "Get() err"},
	{"0.0.0.1", "Connect() err"},
}

// TestGetSNMPValue testing
func TestGetSNMPValue(t *testing.T) {
	oid := ".1.3.6.1.2.1.43.11.1.1.8.1.1"
	for _, mt := range mulSNMPTests {
		*host = mt.host
		_, err := getSNMPValue(oid)
		if err != nil {
			if !strings.HasPrefix(err.Error(), mt.expected) {
				t.Errorf("getSNMPValue error does not Match")
			}
		}
	}
}

// TestMain testing
func TestMain(t *testing.T) {
	sPattern := regexp.MustCompile(`^(\w+) Toner (OK|LOW) -- Toner Levels: (\d+) of (\d+) Remaining \| (\d+)`)
	hostListBlack := []string{"192.168.101.80"}
	for h := 0; h < len(hostListBlack); h++ {
		*host = hostListBlack[h]
		for _, mt := range mulColorBrand {
			if mt.color == "K" {
				*color = mt.color
				*brand = mt.brand
				message := captureStdout(main)
				sc := sPattern.FindSubmatch([]byte(message))
				if len(sc) > 0 {
					if !sPattern.Match([]byte(message)) {
						t.Errorf("Main output not correct: %s %s\n", mt.brand, mt.color)
					}
				}
			}
		}
	}
	hostListColor := []string{"192.168.101.28"}
	for h := 0; h < len(hostListColor); h++ {
		*host = hostListColor[h]
		for _, mt := range mulColorBrand {
			*color = mt.color
			*brand = mt.brand
			message := captureStdout(main)
			sc := sPattern.FindSubmatch([]byte(message))
			if len(sc) > 0 {
				if !sPattern.Match([]byte(message)) {
					t.Errorf("Main output not correct: %s %s\n", mt.brand, mt.color)
				}
			}
		}
	}
	hostListBlank := []string{""}
	for h := 0; h < len(hostListBlank); h++ {
		*host = hostListBlank[h]
		for _, mt := range mulColorBrand {
			*color = mt.color
			*brand = mt.brand
			message := captureStdout(main)
			if message != "Host not set\n" {
				t.Errorf("Not matching: Host not set")
			}
		}
	}
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
