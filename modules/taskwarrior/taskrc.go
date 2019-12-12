/**
 * This package was copied directly from https://github.com/jubnzv/go-taskwarrior/blob/master/taskrc.go
 **/
package taskwarrior

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// Default configuration path.
var TASKRC = PathExpandTilda("~/.taskrc")

// Describes configuration file entries that currently supported by this library.
type TaskRC struct {
	ConfigPath   string // Location of this .taskrc
	DataLocation string `taskwarrior:"data.location"`
}

// Regular expressions that describes parser rules.
var reEntry = regexp.MustCompile(`^\s*([a-zA-Z\.]+)\s*=\s*(.*)\s*$`)
var reInclude = regexp.MustCompile(`^\s*include\s*(.*)\s*$`)

// Expand tilda in filepath as $HOME of current user.
func PathExpandTilda(path string) string {
	if len(path) < 2 {
		return path
	}

	fixedPath := path
	if fixedPath[:2] == "~/" {
		userDir, _ := user.Current()
		homeDir := userDir.HomeDir
		fixedPath = filepath.Join(homeDir, fixedPath[2:])
	}

	return fixedPath
}

// Return new TaskRC struct instance that contains fields with given configuration file values.
func ParseTaskRC(configPath string) (*TaskRC, error) {
	// Fix '~' in a path
	configPath = PathExpandTilda(configPath)

	// Use default configuration file as we need
	if configPath == "" {
		configPath = TASKRC
	} else if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	// Read the given configuration file content in temporary buffer
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Map content in new TaskRC instance
	task := TaskRC{ConfigPath: configPath}
	task.MapTaskRC(string(buf[:]))

	return &task, nil
}

// Map buffer values to given TaskRC struct.
func (c *TaskRC) MapTaskRC(buf string) error {
	// Since we need a little part of all available configuration values we can just traverse line-by-line and check
	// that key from this line represents in out structure. Otherwise skip this line and continue.
	avaialbleKeys := GetAvailableKeys()
	lines := strings.Split(buf, `\n`)
	for _, line := range lines {
		// Remove comments
		line = StripComments(line)

		// Here is an empty line: continue
		if len(line) == 0 {
			continue
		}

		// Is there include pattern?
		res := reInclude.FindStringSubmatch(line)
		if len(res) >= 2 {
			continue // TODO: Not implemented
		}

		// Is there regular configuration entry?
		res = reEntry.FindStringSubmatch(line)
		if len(res) >= 3 {
			// Fill the structure
			keyTag, val := res[1], res[2]
			for _, k := range avaialbleKeys {
				// Check field tag
				field, _ := reflect.TypeOf(c).Elem().FieldByName(k)
				tag := field.Tag
				if tag.Get("taskwarrior") != keyTag {
					continue
				}

				// Set the value
				ps := reflect.ValueOf(c)
				s := ps.Elem()
				if s.Kind() == reflect.Struct {
					f := s.FieldByName(k)
					if f.IsValid() {
						if f.CanSet() {
							if f.Kind() == reflect.String {
								f.SetString(val)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// Return list of available configuration options represented by TaskRC structure fields.
func GetAvailableKeys() []string {
	var availableKeys []string
	t := &TaskRC{}
	s := reflect.ValueOf(t).Elem()
	typeOf := s.Type()
	for i := 0; i < s.NumField(); i++ {
		availableKeys = append(availableKeys, typeOf.Field(i).Name)
	}
	return availableKeys
}

// Remove commented part of input string.
func StripComments(line string) string {
	newLine := line
	for i, a := range line {
		if a == '#' {
			newLine = newLine[:i]
			break
		}
	}
	return newLine
}
