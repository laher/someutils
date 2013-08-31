package someutils

import (
	"strings"
)

type Util struct {
	Name            string
	Function        func([]string) error
}

var (
	allUtils = make(map[string]Util)
)

func Register(u Util) {
	allUtils[u.Name] = u
}

func Call(name string, args []string) error {
	return allUtils[name].Function(args)
}

//except '-help'	
func splitSingleHyphenOpts(call []string) []string {
	splut := []string{}
	for _, item := range call {
		if strings.HasPrefix(item, "-") && !strings.HasPrefix(item, "--") &&
			item != "-help" {
			for _, letter := range item[1:] {
				splut = append(splut, "-"+string(letter))
			}
		} else {
			splut = append(splut, item)
		}
	}
	return splut
}
