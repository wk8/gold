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

func foldOverCanonicalRunes(bytes []byte, f func(rn rune, rnSize int)) {
	for len(bytes) != 0 {
		// extract the next rune, and reduce to its canonical equivalent rune,
		// defined as being the biggest ASCII equivalent one if one exists,
		// otherwise the smallest one
		var r rune
		n := 1

		// fast check for ASCII
		if bytes[0] < utf8.RuneSelf {
			r = rune(bytes[0])
		} else {
			// not ASCII, let's extract the rune
			r, n = utf8.DecodeRune(bytes)

			// and let's cycle through equivalent unicode runs until we hit the
			// biggest
			current, next := r, utf8.MaxRune+1
			for {
				next = unicode.SimpleFold(current)
				if next <= r {
					break
				}
				current = next
			}

			if next < utf8.RuneSelf {
				// we've found an ASCII equivalent, let it fall through to the
				// ASCII case below
				r = next
			} else {
				// no ASCII equivalent, we keep the biggest equivalent
				r = current
			}
		}

		// if ASCII, let's convert upper to lower
		if r < utf8.RuneSelf {
			if 'A' <= r && r <= 'Z' {
				r += asciiUpperToLowerDiff
			}
			f(r, 1)
		} else {
			f(r, utf8.RuneLen(r))
		}

		bytes = bytes[n:]
	}
}
