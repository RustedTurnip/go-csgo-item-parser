package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-collections/collections/stack"
)

type token int

const (
	unknown token = iota

	// root represents the opening position of the file outside of
	// any sections
	root

	// section represents a named section of data
	section

	// opener represents the opening token of a section
	opener

	// closer represents the closing token of a section
	closer

	// data represents a line of data
	data

	// openData represents a line of data that hasn't been concluded
	// (spills over to another line)
	openData

	// empty lines (e.g. comments or whitespace) are to be ignored
	empty
)

const (
	whitespaceCutset = " \t"
)

var (
	language = map[token]map[token]interface{}{
		root: {
			section: struct{}{},
		},
		section: {
			opener: struct{}{},
		},
		opener: {
			section:  struct{}{},
			data:     struct{}{},
			openData: struct{}{},
		},
		closer: {
			section:  struct{}{},
			data:     struct{}{},
			openData: struct{}{},
			closer:   struct{}{},
		},
		data: {
			section:  struct{}{},
			data:     struct{}{},
			openData: struct{}{},
			opener:   struct{}{},
			closer:   struct{}{},
		},
		openData: {
			data:     struct{}{},
			openData: struct{}{},
		},
	}
)

func parse(fileLocation string) (map[string]interface{}, error) {

	// initialise/reset
	dataTree := make(map[string]interface{})
	openSections := stack.New()
	openSections.Push(dataTree)

	lastToken := root

	fi, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(fi)

	lineCount := 0
	currentLine := ""

	// process lines
	for s.Scan() {

		lineCount++

		if s.Err() != nil {
			return nil, fmt.Errorf("failed to scan line %d with error: %s", lineCount, s.Err())
		}

		if lastToken == openData {
			currentLine += s.Text()
		} else {
			currentLine = s.Text()
		}

		t, lineData, err := getLineType(currentLine)
		if err != nil {
			fmt.Println(s.Text())
			return nil, fmt.Errorf("unable to parse line %d with error: %s", lineCount, err.Error())
		}

		// ignore empty lines
		if t == empty {
			currentLine = ""
			continue
		}

		// if token is unexpected
		if _, ok := language[lastToken][t]; !ok {
			return nil, fmt.Errorf("unexpected token on line %d", lineCount)
		}

		lastToken = t

		// process token
		switch t {

		// section: add new section to current section and add to stack
		case section:
			currentSection := openSections.Peek().(map[string]interface{})
			subsection := lineData[0]

			// if subsection already exists in currentSection
			if m, ok := currentSection[subsection]; ok {
				openSections.Push(m)
				continue
			}

			// else add it and push onto stack
			currentSection[subsection] = make(map[string]interface{})
			openSections.Push(currentSection[subsection])
			continue

		// closer: close currently open section
		case closer:
			if openSections.Len() == 0 {
				return nil, errors.New("unexpected token")
			}

			openSections.Pop()
			continue

		// data: add data to currently open section
		case data:
			currentSection := openSections.Peek().(map[string]interface{})
			currentSection[lineData[0]] = lineData[1]
			continue

		// opener, empty: ignore line as doesn't contain data
		case opener, openData, empty:
			continue
		}
	}

	return dataTree, nil
}

func getLineType(line string) (token, []string, error) {

	line = strings.Trim(line, whitespaceCutset)

	// if line is 0 chars after trim (or is a comment), it is a blank line to ignore
	if len(line) == 0 || strings.HasPrefix(line, "//") {
		return empty, nil, nil
	}

	// opener/closer is a single { or } (respectively)
	if len(line) == 1 {

		if line == "{" {
			return opener, nil, nil
		}

		if line == "}" {
			return closer, nil, nil
		}
	}

	split, open := parseDataLine(line)
	if open {
		return openData, nil, nil
	}

	if len(split) == 1 {
		return section, split, nil
	}

	if len(split) == 2 {
		return data, split, nil
	}

	return unknown, nil, errors.New("unrecognised line type")
}

// parseLineData returns the string sub-elements of a data line as a slice
// of strings, and a boolean to highlight whether the line is open ended.
func parseDataLine(line string) ([]string, bool) {

	subStrings := make([]string, 0)
	currentString := ""

	quoted := false

	for i, c := range line {

		// break if string is comment (comments begin with '//')
		if !quoted && c == '/' {
			if i < len(line)-1 {
				if line[i+1] == '/' {
					break
				}
			}
		}

		// if we reach a quote mark
		if c == '"' {

			// and the quote isn't escaped
			if i == 0 || line[i-1] != '\\' {

				// add current string if necessary and reset
				if quoted {
					subStrings = append(subStrings, currentString)
					currentString = ""
				}

				// flip quoted
				quoted = !quoted
				continue
			}
		}

		// ignore anything outside of quotes
		if !quoted {
			continue
		}

		currentString += string(c)
	}

	return subStrings, quoted
}
