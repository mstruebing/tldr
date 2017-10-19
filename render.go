package tldr

import (
	"bufio"
	"io"
	"strings"
)

// Output terms
const (
	BLUE  = "\x1b[34;1m"
	GREEN = "\x1b[32;1m"
	RED   = "\x1b[31;1m"
	RESET = "\x1b[30;1m"
)

// Render takes the given input and renders it for a prettier output.
func Render(markdown io.Reader) (string, error) {
	var rendered string
	var renderingExample bool
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if renderingExample {
			// Skip the empty line
			scanner.Scan()
			line = scanner.Text()

			line = strings.Replace(line, "{{", BLUE, -1)
			line = strings.Replace(line, "}}", RED, -1)
			rendered += "\t" + RED + strings.Trim(line, "`") + RESET + "\n"

			renderingExample = false
		} else if strings.HasPrefix(line, "#") {
			// Heading
			rendered += line[2:] + "\n"
		} else if strings.HasPrefix(line, ">") {
			// Quote
			rendered += line[2:] + "\n"
		} else if strings.HasPrefix(line, "-") {
			// Example
			rendered += GREEN + line + RESET + "\n"
			renderingExample = true
		} else {
			rendered += line + "\n"
		}
	}
	return rendered, scanner.Err()
}

// Write is a convenience function that calls Render and writes the output
// to the destination.
func Write(markdown io.Reader, dest io.Writer) error {
	out, err := Render(markdown)
	if err != nil {
		return err
	}
	_, err = io.WriteString(dest, out)
	return err
}
