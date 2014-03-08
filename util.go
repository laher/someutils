package someutils

const (
	VERSION = "0.5.0"
)

type Util struct {
	Name     string
	Function func([]string) error
}

var (
	allUtils = make(map[string]Util)
)

func Register(u Util) {
	allUtils[u.Name] = u
}

func Exists(name string) bool {
	_, exists := allUtils[name]
	return exists
}

func Call(name string, args []string) error {
	return allUtils[name].Function(args)
}

func List() []string {
	ret := []string{}
	for k, _ := range allUtils {
		ret = append(ret, k)
	}
	return ret
}
