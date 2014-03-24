package someutils

import "io"

const (
	VERSION = "0.5.0"
)

//The Util is the simplest form of utility. By itself it can be registered as a commandline tool, but NOT a pipeline-able function
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

//an Execable can be executed on a pipeline
type Execable interface {
	Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error
}

//Someutil represents a util which can be initialized by flags & executed on a Pipeline
type SomeUtil interface {
	Execable
	ParseFlags(call []string, errOut io.Writer) error
	Name() string
}

type SomeFunc func() SomeUtil

func RegisterSome(somefunc SomeFunc) {

	inPipe, outPipe, errPipe := StdPipes()

	Register(Util{somefunc().Name(), func(call []string) error {
		someutil := somefunc()
		err := someutil.ParseFlags(call, errPipe)
		if err != nil {
			return err
		}
		err = someutil.Exec(inPipe, outPipe, errPipe)
		return err
	}})

}


type ArchiveItem struct {
	//if FileSystemPath is empty, use Data instead
	FileSystemPath string
	ArchivePath    string
	Data           []byte
}
