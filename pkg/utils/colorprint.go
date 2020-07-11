package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/zyedidia/highlight"
)

// PrintError prints out an error message in red
func PrintError(message string) {
	fmt.Printf("snap: ")
	color.Set(color.FgRed)
	fmt.Println(message)
	color.Unset()
}

// PrintErrorMessage prints out a message in red followed by an error on the next line
func PrintErrorMessage(message string, err error) {
	fmt.Printf("snap: ")
	color.Set(color.FgRed)
	fmt.Println(message)
	color.Unset()
	fmt.Printf("error: ")
	color.Set(color.FgRed)
	fmt.Println(err)
	color.Unset()
}

// PrintJSON prints out a byte slice as colorized JSON
func PrintJSON(input []byte) {
	f := colorjson.NewFormatter()
	f.Indent = 2

	var array []map[string]interface{}
	json.Unmarshal(input, &array)
	if len(array) > 0 {
		output := &bytes.Buffer{}
		err := json.Indent(output, input, "", "  ")
		if err != nil {
			PrintError("could not format response as json")
			fmt.Println(string(input))
			os.Exit(1)
		}
		fmt.Println(output.String())
	} else {
		var obj map[string]interface{}
		json.Unmarshal(input, &obj)
		s, _ := f.Marshal(obj)
		fmt.Println(string(s))
	}
}

// PrintMessage prints out a message in green
func PrintMessage(message string) {
	fmt.Printf("snap: ")
	color.Set(color.FgGreen)
	fmt.Println(message)
	color.Unset()
}

// PrintStatus prints out a status code and optional message
func PrintStatus(status string, message string) {
	fmt.Printf("snap: operation status: ")
	if status == "success" {
		color.Set(color.FgGreen)
		fmt.Println(status)
	} else {
		color.Set(color.FgRed)
		fmt.Println(status)
		color.Unset()
		fmt.Print("message: ")
		color.Set(color.FgRed)
		fmt.Println(message)
		color.Unset()
	}
}

// PrintYAML prints out a string as colorized YAML
func PrintYAML(inputString string) {

	// get yaml syntax file as a string
	syntaxFile := yamlSyntax()

	// Parse it into a `*highlight.Def`
	syntaxDef, err := highlight.ParseDef(syntaxFile)
	if err != nil {
		PrintError(fmt.Sprintf("error parsing definition\nerror: %s\n", err))
		os.Exit(1)
	}

	// Make a new highlighter from the definition
	h := highlight.NewHighlighter(syntaxDef)
	// Highlight the string
	// Matches is an array of maps which point to groups
	// matches[lineNum][colNum] will give you the change in group at that line and column number
	// Note that there is only a group at a line and column number if the syntax highlighting changed at that position
	matches := h.HighlightString(inputString)

	// We split the string into a bunch of lines
	// Now we will print the string
	lines := strings.Split(inputString, "\n")
	for lineN, l := range lines {
		for colN, c := range l {
			// Check if the group changed at the current position
			if group, ok := matches[lineN][colN]; ok {
				// Check the group name and set the color accordingly (the colors chosen are arbitrary)
				if group == highlight.Groups["statement"] {
					color.Set(color.FgGreen)
				} else if group == highlight.Groups["preproc"] {
					color.Set(color.FgHiRed)
				} else if group == highlight.Groups["special"] {
					color.Set(color.FgBlue)
				} else if group == highlight.Groups["constant.string"] {
					color.Set(color.FgCyan)
				} else if group == highlight.Groups["constant.specialChar"] {
					color.Set(color.FgHiMagenta)
				} else if group == highlight.Groups["type"] {
					color.Set(color.FgYellow)
				} else if group == highlight.Groups["constant.number"] {
					color.Set(color.FgCyan)
				} else if group == highlight.Groups["comment"] {
					color.Set(color.FgHiGreen)
				} else {
					color.Unset()
				}
			}
			// Print the character
			fmt.Print(string(c))
		}
		// This is at a newline, but highlighting might have been turned off at the very end of the line so we should check that.
		if group, ok := matches[lineN][len(l)]; ok {
			if group == highlight.Groups["default"] || group == highlight.Groups[""] {
				color.Unset()
			}
		}

		fmt.Print("\n")
	}
}

// embed the yaml syntax file
func yamlSyntax() []byte {
	return []byte(`filetype: yaml
detect:
  filename: "\\.ya?ml$"
  header: "%YAML"

rules:
  - type: "(^| )!!(binary|bool|float|int|map|null|omap|seq|set|str) "
  - constant:  "\\b(YES|yes|Y|y|ON|on|NO|no|N|n|OFF|off)\\b"
  - constant: "\\b(true|false)\\b"
  - statement: "(:[[:space:]]|\\[|\\]|:[[:space:]]+[|>]|^[[:space:]]*- )"
  - identifier: "[[:space:]][\\*&][A-Za-z0-9]+"
  - type: "[-.\\w]+:"
  - statement: ":"
  - special:  "(^---|^\\.\\.\\.|^%YAML|^%TAG)"

  - constant.string:
      start: "\""
      end: "\""
      skip: "\\\\."
      rules:
        - constant.specialChar: "\\\\."

  - constant.string:
      start: "'"
      end: "'"
      skip: "\\\\."
      rules:
        - constant.specialChar: "\\\\."

  - comment:
      start: "#"
      end: "$"
      rules: []`)
}
