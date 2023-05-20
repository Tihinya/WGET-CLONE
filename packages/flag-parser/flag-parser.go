package flag_parser

import (
	"fmt"
	"regexp"
	"wget/packages/utils"
)

type Flag struct {
	value       string
	description string
	isBool      bool
}

type flagParser struct {
	whitelist map[string]Flag
}

type flagStorage struct {
	flags map[string]Flag
	tags  []string
}

var re *regexp.Regexp

func CreateParser() *flagParser {
	return &flagParser{
		whitelist: make(map[string]Flag),
	}
}

func (f *flagParser) Add(flagName, description string, isBool bool) *flagParser {
	if _, exist := f.whitelist[flagName]; !exist {
		f.whitelist[flagName] = Flag{
			description: description,
			isBool:      isBool,
		}
	}

	return f
}

func (p *flagParser) Parse(args []string) (*flagStorage, error) {
	storage := &flagStorage{
		flags: make(map[string]Flag),
	}

	for _, value := range args {
		match := re.FindStringSubmatch(value)

		if len(match) < 3 {
			storage.tags = append(storage.tags, value)
			continue
		}

		flagName := ""

		if match[2] != "" {
			flagName = match[2]
		} else if match[1] != "" {
			flagName = match[1]
		}

		if flag, exist := p.whitelist[flagName]; exist {
			if flag.isBool && match[3] != "" {
				return nil, fmt.Errorf(`bool flag "%s" cannot have value`, flagName)
			}

			if len(match) > 3 && match[3] != "" {
				flag.value = match[3]
			}
			storage.flags[flagName] = flag
			continue
		}

		return nil, fmt.Errorf("flag %s not found", flagName)
	}

	return storage, nil
}

func (storage flagStorage) HasFlag(flagName string) bool {
	_, exist := storage.flags[flagName]

	return exist
}

func (storage flagStorage) GetFlag(flagName string) (*Flag, error) {
	if flag, exists := storage.flags[flagName]; exists {
		return &flag, nil
	}
	return nil, fmt.Errorf("no such flag %s", flagName)
}

func (storate flagStorage) GetTags() []string {
	return storate.tags
}

func (storate flagStorage) ArgsExcluded(names ...string) []string {
	args := make([]string, 0)
	for name, flag := range storate.flags {
		if utils.IsContains(names, name) {
			continue
		}

		temp := name
		if len(name) > 1 {
			temp = "-" + temp
		}
		temp = "-" + temp

		if !flag.isBool {
			temp += "=" + flag.value
		}

		args = append(args, temp)
	}

	args = append(args, storate.tags...)

	return args
}

func (flag Flag) GetValue() string {
	return flag.value
}

func init() {
	re = regexp.MustCompile(`^-(?:(?:-(?P<long>[a-zA-Z\-]{2,}))|(?P<short>[a-zA-Z]{1}))\s*(?:=\s*(?P<value>[a-zA-Z0-9/.,]{1,}))?$`)
}
