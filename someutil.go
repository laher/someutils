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

	pipes := StdPipes()

	Register(Util{somefunc().Name(), func(call []string) error {
		someutil := somefunc()
		err := someutil.ParseFlags(call, pipes.Err())
		if err != nil {
			return err
		}
		err = someutil.Exec(pipes)
		return err
	}})

}
