package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type TDP struct {
	Distance float64 // in km
	Pace     time.Duration
	Time     time.Duration
}

// distance (D) * pace (P) = time (T)
// time (T) / distance (D) = pace (P)
// time (T) / pace (P) = distance (D)

// 10km * 4m30s = 45m
// 45m / 10km = 4m30s
// 45m / 4m30s = 10km

// in km
func ParseDistance(s string) (float64, bool) {
	s = strings.ReplaceAll(s, " ", "") // remove spaces

	if !strings.HasSuffix(s, "km") {
		return 0, false
	}

	s = strings.TrimSuffix(s, "km")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}

	return f, true
}

func ParseTime(s string) (time.Duration, bool) {
	s = strings.ReplaceAll(s, " ", "") // remove spaces

	t, err := time.ParseDuration(s)
	if err != nil {
		return 0, false
	}

	return t, true
}

func ParseExpression(s string) (*TDP, error) {
	op, o1, o2, err := ParseOperator(s)
	if err != nil {
		return nil, fmt.Errorf("ParseOperator: %w", err)
	}

	var tdp TDP
	switch op {
	case DivisionOperator:
		t, ok := ParseTime(o1)
		if !ok {
			return nil, fmt.Errorf("unable to parse time %q", o1)
		}
		tdp.Time = t

		d, ok := ParseDistance(o2)
		if ok {
			tdp.Distance = d
			tdp.Pace = time.Duration(int64(math.Round(float64(t) / d)))

			return &tdp, nil
		}

		p, ok := ParseTime(o2)
		if ok {
			tdp.Pace = p
			tdp.Distance = float64(t) / float64(p)

			return &tdp, nil
		}

		return nil, fmt.Errorf("unable to parse distance/pace %q", o2)

	case MultiplicationOperator:
		d, ok := ParseDistance(o1)
		if !ok {
			return nil, fmt.Errorf("unable to parse distance %q", o1)
		}
		tdp.Distance = d

		p, ok := ParseTime(o2)
		if !ok {
			return nil, fmt.Errorf("unable to parse pace %q", o2)
		}
		tdp.Pace = p

		tdp.Time = time.Duration(d * float64(p))

		return &tdp, nil
	}

	return nil, fmt.Errorf("unknown operator %q", op)
}

func main() {
	expr := os.Args[1]

	tdp, err := ParseExpression(expr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%.2f * %s = %s\n", tdp.Distance, tdp.Pace, tdp.Time)
}

const (
	DivisionOperatorSign       = "/"
	MultiplicationOperatorSign = "*"
)

type Operator uint

var (
	UnknownOperator        Operator = 0
	DivisionOperator       Operator = 1
	MultiplicationOperator Operator = 2
)

func ParseOperator(s string) (Operator, string, string, error) {
	if strings.Contains(s, DivisionOperatorSign) {
		operands := strings.Split(s, DivisionOperatorSign)
		if len(operands) > 2 {
			return UnknownOperator, "", "", fmt.Errorf("illegal number of operands: %d", len(strings.Split(s, DivisionOperatorSign)))
		}

		return DivisionOperator, operands[0], operands[1], nil
	}

	if strings.Contains(s, MultiplicationOperatorSign) {
		operands := strings.Split(s, MultiplicationOperatorSign)
		if len(operands) > 2 {
			return UnknownOperator, "", "", fmt.Errorf("illegal number of operands: %d", len(strings.Split(s, MultiplicationOperatorSign)))
		}

		return MultiplicationOperator, operands[0], operands[1], nil
	}

	return UnknownOperator, "", "", nil
}
