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

func Unpack(input string) (string, error) {
	ln := utf8.RuneCountInString(input)

	isFinished, output, err := preUnpack(input, ln)
	if err != nil {
		return "", err
	}

	if isFinished {
		return output, err
	}

	result := make([]byte, 0, ln)

	var letter rune

	var isWriteSlash bool

	for _, nextLetter := range input {
		if isSkipLetter(nextLetter) {
			continue
		}

		var cntRepeat int

		var errParse error

		cntRepeat, errParse = strconv.Atoi(string(nextLetter))

		switch {
		case isSlash(letter) && isSlash(nextLetter) && !isWriteSlash:
			isWriteSlash = true
			letter = nextLetter
		case errParse != nil:
			result, letter, isWriteSlash, err = singleAdd(result, letter, nextLetter, isWriteSlash)
		case errParse == nil:
			result, letter, isWriteSlash, err = multiAdd(result, letter, nextLetter, cntRepeat, isWriteSlash)
		}

		if err != nil {
			return "", err
		}
	}

	result, _, _, err = singleAdd(result, letter, 0, isWriteSlash)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func isSkipLetter(nextLetter rune) bool {
	return nextLetter == '\n' || nextLetter == '\t'
}

func isNotAvailableShielding(letter, nextLetter rune, isWriteSlash bool) bool {
	return isSlash(letter) && !isSlash(nextLetter) && !isWriteSlash
}

func multiAdd(
	result []byte,
	letter,
	nextLetter rune,
	cntRepeat int,
	isWriteSlash bool,
) ([]byte, rune, bool, error) {
	if letter == 0 {
		return nil, 0, isWriteSlash, ErrInvalidString
	}

	if isSlash(letter) && !isWriteSlash {
		return result, nextLetter, isWriteSlash, nil
	}

	result = append(result, repeatRune(letter, cntRepeat)...)

	return result, 0, false, nil
}

func singleAdd(result []byte, letter, nextLetter rune, isWriteSlash bool) ([]byte, rune, bool, error) {
	if isNotAvailableShielding(letter, nextLetter, isWriteSlash) {
		return nil, 0, false, ErrInvalidString
	}

	if letter != 0 {
		result = append(result, repeatRune(letter, 1)...)

		return result, nextLetter, false, nil
	}

	return result, nextLetter, isWriteSlash, nil
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
