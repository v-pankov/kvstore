package item

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type KeyValidator interface {
	KeyMaxLen() int
	ForbiddenRunes() string
}

func (k Key) Validate(v KeyValidator) error {
	if len(k) > v.KeyMaxLen() {
		return ErrKeyLength{
			Max: v.KeyMaxLen(),
			Got: len(k),
		}
	}

	forbiddenCharIndex := strings.IndexAny(string(k), v.ForbiddenRunes())
	if forbiddenCharIndex >= 0 {
		forbiddenChar, _ := utf8.DecodeRuneInString(string(k)[forbiddenCharIndex:])
		return ErrKeyForbiddenChar{
			Rune:  forbiddenChar,
			Index: forbiddenCharIndex,
		}
	}

	return nil
}

type ErrKeyLength struct {
	Max int
	Got int
}

func (e ErrKeyLength) Error() string {
	return fmt.Sprintf(
		"key length [%d] is greater than the maximum [%d]",
		e.Got, e.Max,
	)
}

type ErrKeyForbiddenChar struct {
	Rune  rune
	Index int
}

func (e ErrKeyForbiddenChar) Error() string {
	return fmt.Sprintf(
		"key contains forbidden character [%s] at index [%d]",
		string(e.Rune), e.Index,
	)
}

type keyValidator struct {
	keyMaxLen      int
	forbiddenRunes string
}

func (v keyValidator) KeyMaxLen() int {
	return v.keyMaxLen
}

func (v keyValidator) ForbiddenRunes() string {
	return v.forbiddenRunes
}

func NewKeyValidator(keyMaxLen int, forbiddenRunes string) KeyValidator {
	return keyValidator{
		keyMaxLen:      keyMaxLen,
		forbiddenRunes: forbiddenRunes,
	}
}
