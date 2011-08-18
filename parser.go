// TODO: Bug. Parameter can start quoted with one char and closed with another.
//       For instance: 'PARAM" or "PARAM'
//
// TODO: Quote chars (", ') are ignored/processed_wrong if they are appears
//       in the middle of the word. For instance: "PA'RAM", P'AR"AM
package cue

import (
	"os"
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"unicode"
)

// parseCommand retrive string line and parses it with the following algorythm:
// * first word in the line is command name (cmd return value)
// * all rest words are command's parameters
// * if parameter includes more than one word it should be wrapped with ' or "
func parseCommand(line string) (cmd string, params []string, err os.Error) {
	line = strings.TrimSpace(line)
	params = make([]string, 0)

	// Find cmd.
	i := strings.IndexFunc(line, unicode.IsSpace)
	if i < 0 { // We have only command without any parameters.
		cmd = line
		return
	}
	cmd = line[:i]
	line = strings.TrimSpace(line[i:])

	// Split parameters.
	l := len(line)
	quoted := false
	param := bytes.NewBufferString("")
	for i = 0; i < l; i++ {
		c := line[i]

		if !quoted && unicode.IsSpace(int(c)) { // Start new parameter.
			params = append(params, param.String())
			param = bytes.NewBufferString("")

			// Jump over any spaces.
			for ; i+1 < l && unicode.IsSpace(int(line[i+1])); i++ {

			}
		} else {
			if c == '\\' { // Escape sequence in the text.
				if i+1 >= l {
					err = fmt.Errorf("Unfinished escape sequence")
					return
				}

				s, e := parseEscapeSequence(line[i : i+2])
				if e != nil {
					err = e
					return
				}
				param.WriteByte(s)
				i++
			} else if c == '\'' || c == '"' { // Start/end quoted parameter.
				quoted = !quoted
			} else {
				param.WriteByte(c)
			}
		}
	}

	params = append(params, param.String())

	return
}

// parseEscapeSequence returns escape character by it's string "source code" equivalent.
func parseEscapeSequence(seq string) (char byte, err os.Error) {
	var m = map[string]byte{
		"\\\"": '"',
		"\\'":  '\'',
		"\\\\": '\\',
		"\\n":  '\n',
		"\\t":  '\t',
	}

	char, ok := m[seq]
	if !ok {
		err = fmt.Errorf("Usupported escape sequence '%s'", seq)
	}

	return
}

// parserTime parses time string and returns separate values.
// Input string format: mm:ss:ff
func parseTime(length string) (min int, sec int, frames int, err os.Error) {
	parts := strings.Split(length, ":")
	if len(parts) != 3 {
		err = os.NewError("Illegal time format. mm:ss:ff should be.")
		return
	}

	min, err = strconv.Atoi(parts[0])
	if err != nil {
		err = os.NewError("Failed to parse minutes. " + err.String())
		return
	}

	sec, err = strconv.Atoi(parts[1])
	if err != nil {
		err = os.NewError("Failed to parse seconds. " + err.String())
		return
	}
	if sec > 59 {
		err = os.NewError("Failed to parse seconds. Seconds value can't be more than 59.")
		return
	}

	frames, err = strconv.Atoi(parts[2])
	if err != nil {
		err = os.NewError("Failed to parse frames value. " + err.String())
		return
	}
	if frames > 74 {
		err = os.NewError("Failed to parse frames. Frames value can't be more than 74.")
		return
	}

	return
}
