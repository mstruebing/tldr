package tldr

import (
	"bufio"
	"io"
	"strings"
)

const (
	BLUE  = "\x1b[34;1m"
	GREEN = "\x1b[32;1m"
	RED   = "\x1b[31;1m"
	RESET = "\x1b[30;1m"
)

func Render(markdown io.Reader) (string, error) {
	var rendered string
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Heading
			rendered += line[2:] + "\n"
		} else if strings.HasPrefix(line, ">") {
			// Quote
			rendered += line[2:] + "\n"
		} else if strings.HasPrefix(line, "-") {
			// List
			rendered += GREEN + line + RESET + "\n"
			// rendered += "\t" + RED + convertExample()
			// fmt.Printf("    %s%s%s\n", RED, convertExample(lines[i+2]), RESET)
			// if i < len(lines)-3 {
			// 	fmt.Println()
			// }
		} else if strings.HasPrefix(line, "`") {
			// Code
			line = strings.Replace(line, "{{", BLUE, -1)
			line = strings.Replace(line, "}}", RED, -1)
			rendered += "\t" + strings.Trim(line, "`") + RESET + "\n"
		} else {
			rendered += line + "\n"
		}
	}
	return rendered, scanner.Err()
}
