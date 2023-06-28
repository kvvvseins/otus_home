package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrMaxString     = errors.New("invalid string, max length is 100")
)

const maxLength = 100

type unpack struct {
	result       []byte
	letter       rune
	nextLetter   rune
	cntRepeat    int
	isWriteSlash bool
}

func Unpack(input string) (string, error) {
	ln := utf8.RuneCountInString(input)

	isFinished, output, err := preUnpack(input, ln)
	if err != nil {
		return "", err
	}

	if isFinished {
		return output, err
	}

	process := &unpack{
		result: make([]byte, 0, ln),
	}

	for _, process.nextLetter = range input {
		if isSkipLetter(process.nextLetter) {
			continue
		}
		var errParse error

		process.cntRepeat, errParse = strconv.Atoi(string(process.nextLetter))

		switch {
		case isSlash(process.letter) && isSlash(process.nextLetter) && !process.isWriteSlash:
			process.isWriteSlash = true
			process.letter = process.nextLetter
		case errParse != nil:
			process.result, process.letter, process.isWriteSlash, err = singleAdd(process)
		case errParse == nil:
			process.result, process.letter, process.isWriteSlash, err = multiAdd(process)
		}

		if err != nil {
			return "", err
		}
	}

	process.nextLetter = 0
	process.result, _, _, err = singleAdd(process)
	if err != nil {
		return "", err
	}

	return string(process.result), nil
}

func isSkipLetter(nextLetter rune) bool {
	return nextLetter == '\n' || nextLetter == '\t'
}

func isNotAvailableShielding(letter, nextLetter rune, isWriteSlash bool) bool {
	return isSlash(letter) && !isSlash(nextLetter) && !isWriteSlash
}

func multiAdd(process *unpack) ([]byte, rune, bool, error) {
	if process.letter == 0 {
		return nil, 0, process.isWriteSlash, ErrInvalidString
	}

	if isSlash(process.letter) && !process.isWriteSlash {
		return process.result, process.nextLetter, process.isWriteSlash, nil
	}

	process.result = append(process.result, repeatRune(process.letter, process.cntRepeat)...)

	return process.result, 0, false, nil
}

func singleAdd(process *unpack) ([]byte, rune, bool, error) {
	if isNotAvailableShielding(process.letter, process.nextLetter, process.isWriteSlash) {
		return nil, 0, false, ErrInvalidString
	}

	if process.letter != 0 {
		process.result = append(process.result, repeatRune(process.letter, 1)...)

		return process.result, process.nextLetter, false, nil
	}

	return process.result, process.nextLetter, process.isWriteSlash, nil
}

func isSlash(letter rune) bool {
	return letter == '\\'
}

func preUnpack(input string, ln int) (isFinished bool, output string, err error) {
	if input == "" {
		return true, "", nil
	}

	if ln == 1 {
		return true, input, nil
	}

	if ln > maxLength {
		return true, "", ErrMaxString
	}

	if _, err = strconv.Atoi(string(input[0])); err == nil {
		return true, "", ErrInvalidString
	}

	return false, "", nil
}

func repeatRune(symbolRune rune, cnt int) []byte {
	if cnt == 0 {
		return nil
	}

	var result []byte

	for i := 0; i < cnt; i++ {
		result = append(result, []byte(string(symbolRune))...)
	}

	return result
}
