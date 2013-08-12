package someutils

import (
	"strings"
)

func splitSingleHyphenOpts(call []string) []string {
	splut := []string{}
	for _, item := range call {
		if strings.HasPrefix(item, "-") && !strings.HasPrefix(item, "--") {
			for _, letter := range item[1:] {
				splut = append(splut, "-"+string(letter))
			}
		} else {
			splut = append(splut, item)
		}
	}
	return splut
}
