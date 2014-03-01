package someutils

import (
	"github.com/laher/uggo"
	"io"
	"os"
)

type TeeOptions struct {
	isAppend bool
}

func init() {
	Register(Util{
		"tee",
		Tee})
}

func Tee(call []string) error {
	options := TeeOptions{}
	flagSet := uggo.NewFlagSetDefault("tee", "", VERSION)
	flagSet.AliasedBoolVar(&options.isAppend, []string{"a", "append"}, false, "Append instead of overwrite")

	err := flagSet.Parse(call[1:])
	if err != nil {
		return err
	}
	if flagSet.ProcessHelpOrVersion() {
		return nil
	}
	flag := os.O_CREATE
	if options.isAppend {
		flag = flag | os.O_APPEND
	}
	writeables := flagSet.ArgsAsWriteables(flag, 0666)
	files, err := uggo.OpenAll(writeables)
	if err != nil {
		return err
	}
	writers := []io.Writer{os.Stdout}
	for _, file := range files {
		writers = append(writers, file)
	}
	multiwriter := io.MultiWriter(writers...)
	_, err = io.Copy(multiwriter, os.Stdin)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
