package builder

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type sourceParser struct {
	path string
	ext  string

	vars map[string]string
	body string
}

func newSourceParser(path string) *sourceParser {
	p := &sourceParser{path: path}
	p.ext = filepath.Ext(path)
	p.vars = make(map[string]string)
	return p
}

func (sp *sourceParser) parse() error {
	file, err := os.Open(sp.path)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := make([]string, 0, 35)
	divided := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		divider, err := regexp.MatchString(VAR_CONTENT_DIVIDER, line)
		if err != nil {
			return err
		}
		if divider && !divided {
			divided = true
			vars, err := sp.parseVars(lines)
			if err != nil {
				return errors.New(sp.path + " " + err.Error())
			}
			sp.vars = vars
			lines = make([]string, 0, 35)
		} else {
			lines = append(lines, line)
		}
	}
	sp.body = strings.Join(lines, "\n")
	return nil
}

func (sp *sourceParser) parseVars(lines []string) (map[string]string, error) {
	vars := make(map[string]string)
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := strings.SplitN(line, NAME_VAL_DIVIDER, 2)
		if len(parts) != 2 {
			return vars, errors.New(string(i) + ": error parsing variables.")
		}

		name := strings.Trim(parts[0], " ")
		value := strings.Trim(parts[1], " ")
		vars[name] = value
	}
	return vars, nil
}
