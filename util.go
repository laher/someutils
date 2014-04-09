package someutils

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







var (
	allCliUtils        = make(map[string]CliPipable)
	allPipables = make(map[string]NamedPipableFactory)
)

// Registers utils for use by 'some' command.
func Register(u CliPipable) {
	allCliUtils[u.Name()] = u
}
func RegisterSimple(somefunc CliPipableSimpleFactory) {
	RegisterPipable(func() NamedPipable { return WrapNamed(somefunc()) })
	Register(WrapUtil(somefunc()))
}

func RegisterPipable(somefunc NamedPipableFactory) {
	pipable := somefunc()
	name := pipable.Name()
	allPipables[name] = somefunc

	//register as CliUtil if possible
	pcu, ok := pipable.(CliPipable)
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
func Call(name string, args []string) (error, int) {
	ps := StdInvocation()
	util := allCliUtils[name]
	return CallUtil(util, args, ps)
}

func CallUtil(util CliPipable, args []string, invocation *Invocation) (error, int) {
	err, code := util.ParseFlags(args, invocation.ErrOutPipe)
	if err != nil {
		return err, code
	}
	return util.Invoke(invocation)
}

func PipableExists(name string) bool {
	_, exists := allPipables[name]
	return exists
}

func GetNamedPipableFactory(name string) NamedPipableFactory {
	return allPipables[name]
}

func GetCliPipableFactory(name string) CliPipableFactory {
	namedPipableFactory := GetNamedPipableFactory(name)
	pipableCliUtilFactory := func () CliPipable {
		return namedPipableFactory().(CliPipable)
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

// ArchiveItem is used by tar & zip
type ArchiveItem struct {
	//if FileSystemPath is empty, use Data instead
	FileSystemPath string
	ArchivePath    string
	Data           []byte
}
