// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains simple golden tests for various examples.
// Besides validating the results when the implementation changes,
// it provides a way to look at the generated code without having
// to execute the print statements in one's head.

package main

import (
	"strings"
	"testing"
)

// Golden represents a test case.
type Golden struct {
	name           string
	trimPrefix     string
	lineComment    bool
	withFromString bool
	input          string // input; the package clause is provided when running the test.
	output         string // exected output.
}

var golden = []Golden{
	{"day", "", false, false, day_in, day_out},
	{"offset", "", false, false, offset_in, offset_out},
	{"gap", "", false, false, gap_in, gap_out},
	{"num", "", false, false, num_in, num_out},
	{"unum", "", false, false, unum_in, unum_out},
	{"prime", "", false, false, prime_in, prime_out},
	{"prefix", "Type", false, false, prefix_in, prefix_out},
	{"tokens", "", true, false, tokens_in, tokens_out},
	{"month", "", true, true, month_in, month_out},
	{"mood", "", true, true, mood_in, mood_out},
	{"fib", "", false, true, fib_in, fib_out},
	{"hundreds", "", false, true, hundreds_in, hundreds_out},
}

// Each example starts with "type XXX [u]int", with a single space separating them.

// Simple test: enumeration of type int starting at 0.
const day_in = `type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
`

const day_out = `
import "strconv"

const _Day_name = "MondayTuesdayWednesdayThursdayFridaySaturdaySunday"

var _Day_index = [...]uint8{0, 6, 13, 22, 30, 36, 44, 50}

func (i Day) String() string {
	if i < 0 || i >= Day(len(_Day_index)-1) {
		return "Day(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Day_name[_Day_index[i]:_Day_index[i+1]]
}
`

// Enumeration with an offset.
// Also includes a duplicate.
const offset_in = `type Number int
const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One  // Duplicate; note that AnotherOne doesn't appear below.
)
`

const offset_out = `
import "strconv"

const _Number_name = "OneTwoThree"

var _Number_index = [...]uint8{0, 3, 6, 11}

func (i Number) String() string {
	i -= 1
	if i < 0 || i >= Number(len(_Number_index)-1) {
		return "Number(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Number_name[_Number_index[i]:_Number_index[i+1]]
}
`

// Gaps and an offset.
const gap_in = `type Gap int
const (
	Two Gap = 2
	Three Gap = 3
	Five Gap = 5
	Six Gap = 6
	Seven Gap = 7
	Eight Gap = 8
	Nine Gap = 9
	Eleven Gap = 11
)
`

const gap_out = `
import "strconv"

const (
	_Gap_name_0 = "TwoThree"
	_Gap_name_1 = "FiveSixSevenEightNine"
	_Gap_name_2 = "Eleven"
)

var (
	_Gap_index_0 = [...]uint8{0, 3, 8}
	_Gap_index_1 = [...]uint8{0, 4, 7, 12, 17, 21}
)

func (i Gap) String() string {
	switch {
	case 2 <= i && i <= 3:
		i -= 2
		return _Gap_name_0[_Gap_index_0[i]:_Gap_index_0[i+1]]
	case 5 <= i && i <= 9:
		i -= 5
		return _Gap_name_1[_Gap_index_1[i]:_Gap_index_1[i+1]]
	case i == 11:
		return _Gap_name_2
	default:
		return "Gap(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
`

// Signed integers spanning zero.
const num_in = `type Num int
const (
	m_2 Num = -2 + iota
	m_1
	m0
	m1
	m2
)
`

const num_out = `
import "strconv"

const _Num_name = "m_2m_1m0m1m2"

var _Num_index = [...]uint8{0, 3, 6, 8, 10, 12}

func (i Num) String() string {
	i -= -2
	if i < 0 || i >= Num(len(_Num_index)-1) {
		return "Num(" + strconv.FormatInt(int64(i+-2), 10) + ")"
	}
	return _Num_name[_Num_index[i]:_Num_index[i+1]]
}
`

// Unsigned integers spanning zero.
const unum_in = `type Unum uint
const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)
`

const unum_out = `
import "strconv"

const (
	_Unum_name_0 = "m0m1m2"
	_Unum_name_1 = "m_2m_1"
)

var (
	_Unum_index_0 = [...]uint8{0, 2, 4, 6}
	_Unum_index_1 = [...]uint8{0, 3, 6}
)

func (i Unum) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _Unum_name_0[_Unum_index_0[i]:_Unum_index_0[i+1]]
	case 253 <= i && i <= 254:
		i -= 253
		return _Unum_name_1[_Unum_index_1[i]:_Unum_index_1[i+1]]
	default:
		return "Unum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
`

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems
const prime_in = `type Prime int
const (
	p2 Prime = 2
	p3 Prime = 3
	p5 Prime = 5
	p7 Prime = 7
	p77 Prime = 7 // Duplicate; note that p77 doesn't appear below.
	p11 Prime = 11
	p13 Prime = 13
	p17 Prime = 17
	p19 Prime = 19
	p23 Prime = 23
	p29 Prime = 29
	p37 Prime = 31
	p41 Prime = 41
	p43 Prime = 43
)
`

const prime_out = `
import (
	"strconv"
	"sync"
)

const _Prime_name = "p2p3p5p7p11p13p17p19p23p29p37p41p43"

var _Prime_map map[Prime]string

var populatePrimeMapOnce sync.Once

func populatePrimeMap() {
	_Prime_map = map[Prime]string{
		2:  _Prime_name[0:2],
		3:  _Prime_name[2:4],
		5:  _Prime_name[4:6],
		7:  _Prime_name[6:8],
		11: _Prime_name[8:11],
		13: _Prime_name[11:14],
		17: _Prime_name[14:17],
		19: _Prime_name[17:20],
		23: _Prime_name[20:23],
		29: _Prime_name[23:26],
		31: _Prime_name[26:29],
		41: _Prime_name[29:32],
		43: _Prime_name[32:35],
	}
}

func (i Prime) String() string {
	populatePrimeMapOnce.Do(populatePrimeMap)
	if str, ok := _Prime_map[i]; ok {
		return str
	}
	return "Prime(" + strconv.FormatInt(int64(i), 10) + ")"
}
`

const prefix_in = `type Type int
const (
	TypeInt Type = iota
	TypeString
	TypeFloat
	TypeRune
	TypeByte
	TypeStruct
	TypeSlice
)
`

const prefix_out = `
import "strconv"

const _Type_name = "IntStringFloatRuneByteStructSlice"

var _Type_index = [...]uint8{0, 3, 9, 14, 18, 22, 28, 33}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
`

const tokens_in = `type Token int
const (
	And Token = iota // &
	Or               // |
	Add              // +
	Sub              // -
	Ident
	Period // .

	// not to be used
	SingleBefore
	// not to be used
	BeforeAndInline // inline
	InlineGeneral /* inline general */
)
`

const tokens_out = `
import "strconv"

const _Token_name = "&|+-Ident.SingleBeforeinlineinline general"

var _Token_index = [...]uint8{0, 1, 2, 3, 4, 9, 10, 22, 28, 42}

func (i Token) String() string {
	if i < 0 || i >= Token(len(_Token_index)-1) {
		return "Token(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Token_name[_Token_index[i]:_Token_index[i+1]]
}
`

// Simple test: enumeration of type int starting at 0.
const month_in = `type Month int
const (
	January Month = iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)
`

const month_out = `
import "strconv"

const _Month_name = "JanuaryFebruaryMarchAprilMayJuneJulyAugustSeptemberOctoberNovemberDecember"

var _Month_index = [...]uint8{0, 7, 15, 20, 25, 28, 32, 36, 42, 51, 58, 66, 74}

func (i Month) String() string {
	if i < 0 || i >= Month(len(_Month_index)-1) {
		return "Month(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Month_name[_Month_index[i]:_Month_index[i+1]]
}

func MonthFromString(s string) (i Month, ok bool) {
	for v := range _Month_index {
		v += 0
		if s == Month(v).String() {
			return Month(v), true
		}
	}
	return
}
`

// Enumeration with an negative offset.
const mood_in = `type Mood int
const (
	Negative Mood = iota -1
	Neutral
	Happy
)
`

const mood_out = `
import "strconv"

const _Mood_name = "NegativeNeutralHappy"

var _Mood_index = [...]uint8{0, 8, 15, 20}

func (i Mood) String() string {
	i -= -1
	if i < 0 || i >= Mood(len(_Mood_index)-1) {
		return "Mood(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _Mood_name[_Mood_index[i]:_Mood_index[i+1]]
}

func MoodFromString(s string) (i Mood, ok bool) {
	for v := range _Mood_index {
		v += -1
		if s == Mood(v).String() {
			return Mood(v), true
		}
	}
	return
}
`

// Gaps and an offset with duplicates.
const fib_in = `type Fib int
const (
	Zero       Fib = 0
	One        Fib = 1
	AnotherOne Fib = 1
	Two        Fib = 2
	Three      Fib = 3
	Five       Fib = 5
	Eight      Fib = 8
	Thirteen   Fib = 13
)
`

const fib_out = `
import "strconv"

const (
	_Fib_name_0 = "ZeroOneTwoThree"
	_Fib_name_1 = "Five"
	_Fib_name_2 = "Eight"
	_Fib_name_3 = "Thirteen"
)

var (
	_Fib_index_0 = [...]uint8{0, 4, 7, 10, 15}
)

func (i Fib) String() string {
	switch {
	case 0 <= i && i <= 3:
		return _Fib_name_0[_Fib_index_0[i]:_Fib_index_0[i+1]]
	case i == 5:
		return _Fib_name_1
	case i == 8:
		return _Fib_name_2
	case i == 13:
		return _Fib_name_3
	default:
		return "Fib(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

func FibFromString(s string) (i Fib, ok bool) {
	switch s {
	case "Zero":
		return Fib(0), true
	case "One":
		return Fib(1), true
	case "Two":
		return Fib(2), true
	case "Three":
		return Fib(3), true
	case "Five":
		return Fib(5), true
	case "Eight":
		return Fib(8), true
	case "Thirteen":
		return Fib(13), true
	}
	return
}
`

// Enough gaps to trigger a map implementation of the method.
// Also includes a duplicate to test that it doesn't cause problems
const hundreds_in = `type Hundred int
const (
	One Hundred = 100
	Two Hundred = 200
	Three Hundred = 300
	Four Hundred = 400
	Five Hundred = 500
	HighFive Hundred = 500 // duplicate
	Six Hundred = 600
	Seven Hundred = 700
	Eight Hundred = 800
	Nine Hundred = 900
	Ten Hundred = 1000
	Eleven Hundred = 1100
	Twelve Hundred = 1200
)
`

const hundreds_out = `
import (
	"strconv"
	"sync"
)

const _Hundred_name = "OneTwoThreeFourFiveSixSevenEightNineTenElevenTwelve"

var _Hundred_map map[Hundred]string

var populateHundredMapOnce sync.Once

func populateHundredMap() {
	_Hundred_map = map[Hundred]string{
		100:  _Hundred_name[0:3],
		200:  _Hundred_name[3:6],
		300:  _Hundred_name[6:11],
		400:  _Hundred_name[11:15],
		500:  _Hundred_name[15:19],
		600:  _Hundred_name[19:22],
		700:  _Hundred_name[22:27],
		800:  _Hundred_name[27:32],
		900:  _Hundred_name[32:36],
		1000: _Hundred_name[36:39],
		1100: _Hundred_name[39:45],
		1200: _Hundred_name[45:51],
	}
}

func (i Hundred) String() string {
	populateHundredMapOnce.Do(populateHundredMap)
	if str, ok := _Hundred_map[i]; ok {
		return str
	}
	return "Hundred(" + strconv.FormatInt(int64(i), 10) + ")"
}

func HundredFromString(s string) (i Hundred, ok bool) {
	populateHundredMapOnce.Do(populateHundredMap)
	for k, v := range _Hundred_map {
		if s == v {
			return k, true
		}
	}
	return
}
`

func TestGolden(t *testing.T) {
	for _, test := range golden {
		g := Generator{
			trimPrefix:     test.trimPrefix,
			lineComment:    test.lineComment,
			withFromString: test.withFromString,
		}
		input := "package test\n" + test.input
		file := test.name + ".go"
		g.parsePackage(".", []string{file}, input)
		// Extract the name and type of the constant from the first line.
		tokens := strings.SplitN(test.input, " ", 3)
		if len(tokens) != 3 {
			t.Fatalf("%s: need type declaration on first line", test.name)
		}
		g.generate(tokens[1])
		got := string(g.format())
		if got != test.output {
			t.Errorf("%s: got\n====\n%s====\nexpected\n====%s", test.name, got, test.output)
		}
	}
}
