package someutils

import "io"

const (
	VERSION = "0.5.0"
)

//The Util is the simplest form of utility. By itself it can be registered as a commandline tool, but NOT a pipeline-able function
type CliUtil struct {
	Name     string
	Function func([]string) error
}

//a Pipable can be executed on a pipeline
type Pipable interface {
	Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error
}

//PipableUtil represents a util which can be initialized by flags & executed on a Pipeline
type PipableCliUtil interface {
	Pipable
	ParseFlags(call []string, errOut io.Writer) error
	Name() string
}

type PipableFunc func() PipableCliUtil

type ArchiveItem struct {
	//if FileSystemPath is empty, use Data instead
	FileSystemPath string
	ArchivePath    string
	Data           []byte
}

var (
	allCliUtils = make(map[string]CliUtil)
)

//Registers utils for use by 'some' command
func Register(u CliUtil) {
	allCliUtils[u.Name] = u
}

func RegisterPipable(somefunc PipableFunc) {

	inPipe, outPipe, errPipe := StdPipes()

	Register(CliUtil{somefunc().Name(), func(call []string) error {
		someutil := somefunc()
		err := someutil.ParseFlags(call, errPipe)
		if err != nil {
			return err
		}
		err = someutil.Exec(inPipe, outPipe, errPipe)
		return err
	}})

}

func Exists(name string) bool {
	_, exists := allCliUtils[name]
	return exists
}

func Call(name string, args []string) error {
	return allCliUtils[name].Function(args)
}

func List() []string {
	ret := []string{}
	for k, _ := range allCliUtils {
		ret = append(ret, k)
	}
	return ret
}



