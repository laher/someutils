package someutils

import "errors"

//Deprecated. See someutils/some.CatCli
func Cat(call []string) error {
	return errors.New("Use someutils/some.CatCli(...)")
}

//Deprecated. See someutils/some.BasenameCli
func Basename(call []string) error {
	return errors.New("Use someutils/some.BasenameCli(...)")
}
