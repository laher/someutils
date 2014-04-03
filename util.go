package someutils

import "io"

const (
	VERSION = "0.5.1-snapshot"
)
/*
// The CliUtil can be registered as a commandline tool, but NOT a pipeline-able function
type CliUtil interface {
	Name() string
	Function func([]string) error
}
*/

//a Pipable can be executed on a pipeline
type Pipable interface {
	Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error
}

// a Named Pipable can be registered for use by e.g. xargs
type NamedPipable interface {
	Pipable
	Name() string
}

//PipableUtil represents a util which can be initialized by flags & executed on a Pipeline
type PipableCliUtil interface {
	NamedPipable
	ParseFlags(call []string, errOut io.Writer) error
}

type PipableFactory func() Pipable
type NamedPipableFactory func() NamedPipable
type PipableCliUtilFactory func() PipableCliUtil

type ArchiveItem struct {
	//if FileSystemPath is empty, use Data instead
	FileSystemPath string
	ArchivePath    string
	Data           []byte
}

var (
	allCliUtils        = make(map[string]PipableCliUtil)
	allPipables = make(map[string]NamedPipableFactory)
)

// Registers utils for use by 'some' command.
func Register(u PipableCliUtil) {
	allCliUtils[u.Name()] = u
}

func RegisterPipable(somefunc NamedPipableFactory) {
	pipable := somefunc()
	name := pipable.Name()
	allPipables[name] = somefunc

	//register as CliUtil if possible
	pcu, ok := pipable.(PipableCliUtil)
	if ok {
		//inPipe, outPipe, errPipe := StdPipes()
		Register(pcu)
			/*CliUtil{ func () string { return name }, func(call []string) error {
			someutil := somefunc()
			err := pcu.ParseFlags(call, errPipe)
			if err != nil {
				return err
			}
			err = pcu.Exec(inPipe, outPipe, errPipe)
			return err*/
	}
}

// deprecated. Use CliExists instead.
func Exists(name string) bool {
	return CliExists(name)
}

// Has a CLI function been registered?
func CliExists(name string) bool {
	_, exists := allCliUtils[name]
	return exists
}

// deprecated. Use GetCliUtil, ParseFlags & Exec instead.
func Call(name string, args []string) error {
	inPipe, outPipe, errPipe := StdPipes()
	util := allCliUtils[name]
	err := util.ParseFlags(args, errPipe)
	if err != nil {
		return err
	}
	err = util.Exec(inPipe, outPipe, errPipe)
	return err
}

func PipableExists(name string) bool {
	_, exists := allPipables[name]
	return exists
}

func GetNamedPipableFactory(name string) NamedPipableFactory {
	return allPipables[name]
}

func GetPipableCliUtilFactory(name string) PipableCliUtilFactory {
	namedPipableFactory := GetNamedPipableFactory(name)
	pipableCliUtilFactory := func () PipableCliUtil {
		return namedPipableFactory().(PipableCliUtil)
	}
	return pipableCliUtilFactory

}
func GetPipableFactory(name string) PipableFactory {
	namedPipableFactory := GetNamedPipableFactory(name)
	pipableFactory := func () Pipable {
		return namedPipableFactory()
	}
	return pipableFactory
}

func List() []string {
	ret := []string{}
	for k, _ := range allCliUtils {
		ret = append(ret, k)
	}
	return ret
}
