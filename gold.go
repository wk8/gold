package gold

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

var CaseFoldString = CaseFoldStringFast
var CaseFoldStringFast = generateCaseFoldString(CaseFoldBytesFast)
var CaseFoldStringLowMem = generateCaseFoldString(CaseFoldBytesLowMem)
var CaseFoldBytes = CaseFoldBytesFast

func generateCaseFoldString(f func(bytes []byte) []byte) func(str string) string {
	return func(str string) string {
		return string(f([]byte(str)))
	}
}

func CaseFoldBytesFast(bytes []byte) []byte {
	newLength, runesCount := 0, 0
	newRunes := make([]rune, len(bytes))

	foldOverCanonicalRunes(bytes, func(rn rune, rnSize int) {
		newLength += rnSize
		newRunes[runesCount] = rn
		runesCount++
	})

	newBytes := make([]byte, newLength)
	current := newBytes
	for i := 0; i < runesCount; i++ {
		rnSize := utf8.EncodeRune(current, newRunes[i])
		current = current[rnSize:]
	}

	return newBytes
}

func CaseFoldBytesLowMem(bytes []byte) []byte {
	newLength := 0
	foldOverCanonicalRunes(bytes, func(_ rune, rnSize int) {
		newLength += rnSize
	})

	newBytes := make([]byte, newLength)
	current := newBytes
	foldOverCanonicalRunes(bytes, func(rn rune, rnSize int) {
		if size := utf8.EncodeRune(current, rn); size != rnSize {
			panic(fmt.Sprintf("Expected to write %v bytes, wrote %v", rnSize, size))
		}
		current = current[rnSize:]
	})

	return newBytes
}

const asciiUpperToLowerDiff = 'a' - 'A'

func foldOverCanonicalRunes(bytes []byte, f func(rune, int)) {
	var rn, current, next rune
	var rnSize int

	for len(bytes) != 0 {
		// extract the next rune, and reduce to its canonical equivalent rune,
		// defined as being the biggest ASCII equivalent one if one exists,
		// otherwise the smallest one

		// fast check for ASCII
		if bytes[0] < utf8.RuneSelf {
			rn, rnSize = rune(bytes[0]), 1
		} else {
			// not ASCII, let's extract the rune
			rn, rnSize = utf8.DecodeRune(bytes)

			// and let's cycle through equivalent unicode runes until we hit the
			// biggest
			current, next = rn, unicode.SimpleFold(rn)
			for next > rn {
				current = next
				next = unicode.SimpleFold(current)
			}

			if next < utf8.RuneSelf {
				// we've found an ASCII equivalent, let it fall through to the
				// ASCII case below
				rn = next
			} else {
				// no ASCII equivalent, we keep the biggest equivalent
				rn = current
			}
		}

		// if ASCII, let's convert upper to lower
		if rn < utf8.RuneSelf {
			if 'A' <= rn && rn <= 'Z' {
				rn += asciiUpperToLowerDiff
			}
			f(rn, 1)
		} else {
			f(rn, utf8.RuneLen(rn))
		}

		bytes = bytes[rnSize:]
	}
}
