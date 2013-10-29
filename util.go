package someutils

import (
	"github.com/laher/uggo"
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
	return uggo.Gnuify(call)
}
