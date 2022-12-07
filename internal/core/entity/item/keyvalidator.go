package item

import (
	"fmt"
	"strings"
)

type KeyValidator interface {
	KeyMaxLen() int
	KeyForbiddenChars() string
}

func (k Key) Validate(v KeyValidator) error {
	if len(k) > v.KeyMaxLen() {
		return ErrKeyLength{
			Max: v.KeyMaxLen(),
			Got: len(k),
		}
	}

	if forbiddenCharIndex := strings.IndexAny(
		string(k), v.KeyForbiddenChars(),
	); forbiddenCharIndex >= 0 {
		return ErrKeyForbiddenChar{
			Char:     rune(k[forbiddenCharIndex]),
			Position: forbiddenCharIndex,
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
	Char     rune
	Position int
}

func (e ErrKeyForbiddenChar) Error() string {
	return fmt.Sprintf(
		"key contains forbidden character [%s] at position [%d]",
		string(e.Char), e.Position,
	)
}

type keyValidator struct {
	keyMaxLen         int
	keyForbiddenChars string
}

func (v keyValidator) KeyMaxLen() int {
	return v.keyMaxLen
}

func (v keyValidator) KeyForbiddenChars() string {
	return v.keyForbiddenChars
}

func NewKeyValidator(keyMaxLen int, forbiddenChars string) KeyValidator {
	return keyValidator{
		keyMaxLen:         keyMaxLen,
		keyForbiddenChars: forbiddenChars,
	}
}
