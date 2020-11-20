package objects

import (
	"fmt"
	"regexp"
)

const idRegexString = "[a-f0-9]{40}"

func init() {
	var err error
	regex, err = regexp.Compile(idRegexString)
	if err != nil {
		panic("couldn't compile id regex")
	}
}

type ID string

var regex *regexp.Regexp

func IdFromSum(sum [20]byte) ID {
	return ID(fmt.Sprintf("%x", sum))
}

func IdFromString(s string) (ID, error) {
	if !regex.MatchString(s) {
		return "", fmt.Errorf("%s is not an id", s)
	}
	return ID(s), nil
}
