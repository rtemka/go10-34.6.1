package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stdout, "Usage: %s <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}

	re := regexp.MustCompile(`(\d+)([\+\-\*\/:])(\d+)\=\?`)

	outName := func() string {
		if len(os.Args) == 3 {
			return filepath.Clean(os.Args[2])
		} else {
			return "out.txt"
		}
	}()

	in, err := os.Open(filepath.Clean(os.Args[1]))
	if err != nil {
		errorExit(err)
	}
	out, err := os.Create(outName)
	if err != nil {
		errorExit(err)
	}
	defer func() {
		fmt.Printf("Processed %s >> %s\n", in.Name(), out.Name())
		_ = in.Close()
		_ = out.Close()
	}()

	w := bufio.NewWriter(out)
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {

		str := scanner.Text()

		m := re.FindAllStringSubmatch(str, -1)
		if len(m) == 0 {
			continue
		}

		l, err := strconv.ParseInt(m[0][1], 10, 0)
		if err != nil {
			errorExit(err)
		}
		r, err := strconv.ParseInt(m[0][3], 10, 0)
		if err != nil {
			errorExit(err)
		}

		res := calc(l, r, m[0][2])

		w.WriteString(fmt.Sprintf("%s%s%s=%s\n", m[0][1], m[0][2], m[0][3], res))
	}

	w.Flush()

	if scanner.Err() != nil {
		errorExit(scanner.Err())
	}
}

func errorExit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func calc(left, right int64, operator string) string {
	switch {
	case operator == "+":
		return fmt.Sprint(left + right)
	case operator == "-":
		return fmt.Sprint(left - right)
	case operator == "*":
		return fmt.Sprint(left * right)
	case (operator == "/" || operator == ":") && right != 0:
		return fmt.Sprint(left / right)
	default:
		return "NaN"
	}
}
