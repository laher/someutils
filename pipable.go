package someutils

import "io"

//a Pipable can be executed on a pipeline
type Pipable interface {
	Invoke(i *Invocation) (error, int)
}

//a PipableSimple can be executed on a pipeline when wrapped inside a PipableSimpleWrapper
type PipableSimple interface {
	Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) (error, int)
}

// a Named Pipable can be registered for use by e.g. xargs
type NamedPipable interface {
	Pipable
	Named
}

//PipableUtil represents a util which can be initialized by flags & executed on a Pipeline
type CliPipable interface {
	NamedPipable
	Cliable
}

type PipableWrapper struct {
	PipableSimple
}

type Named interface {
	Name() string
}
type Cliable interface {
	ParseFlags(call []string, errPipe io.Writer) (error, int)
}

type NamedPipableSimple interface {
	PipableSimple
	Named
}
type CliPipableSimple interface {
	NamedPipableSimple
	Cliable
}
type PipableSimpleWrapper struct {
	PipableSimple
}

type NamedPipableSimpleWrapper struct {
	NamedPipableSimple
}

type CliPipableSimpleWrapper struct {
	CliPipableSimple
}

func Wrap(ps PipableSimple) Pipable {
	return &PipableSimpleWrapper{ps}
}

func WrapNamed(ps NamedPipableSimple) NamedPipable {
	return &NamedPipableSimpleWrapper{ps}
}

func WrapCliPipable(ps CliPipableSimple) CliPipable {
	return &CliPipableSimpleWrapper{ps}
}

/*
func (w PipableWrapper) ExecFull(inPipe io.Reader, outPipe io.Writer, errInPipe io.Reader, errOutPipe io.Writer, signalChan chan int) (error, int) {
	go autoPipe(errMainPipe.Out, errMainPipe.In)
	return w.PipableSimple.Exec(inPipe, outPipe, errMainPipe.Out)
}
*/

func (npsw *NamedPipableSimpleWrapper) Invoke(i *Invocation) (error, int) {
	return invoke(npsw.NamedPipableSimple, i)
}

func (psw *PipableSimpleWrapper) Invoke(i *Invocation) (error, int) {
	return invoke(psw.PipableSimple, i)
}
func (cpsw *CliPipableSimpleWrapper) Invoke(i *Invocation) (error, int) {
	return invoke(cpsw.CliPipableSimple, i)
}

type PipableFactory func() Pipable

//type NamedPipableFactory func() NamedPipable
type CliPipableFactory func() CliPipable
type PipableSimpleFactory func() PipableSimple
type NamedPipableSimpleFactory func() NamedPipableSimple
type CliPipableSimpleFactory func() CliPipableSimple
