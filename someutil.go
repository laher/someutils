package someutils

import "io"

//Interface representing a util which can be invoked independently on in a Pipeline
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
