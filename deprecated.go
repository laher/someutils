package someutils

import "errors"

//Deprecated. See someutils/some.CatCli
func Cat(call []string) error {
	return errors.New("Deprecated function. Please see package github.com/someutils/some")
}

//Deprecated. See someutils/some.BasenameCli
func Basename(call []string) error {
	return errors.New("Deprecated function. Please see package github.com/someutils/some")
}


//Deprecated. See someutils/some.CpCli
func Cp(call []string) error {
	return errors.New("Deprecated function. Please see package github.com/someutils/some")
}
