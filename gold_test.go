package gold

import (
	"testing"
)

func TestCaseFoldString(t *testing.T) {
	inputToExpected := map[string]string{
		// hello => hello
		"hello": "hello",

		// hELlo => hello
		"hELlo": "hello",

		// with a null byte
		"HeL\000oO": "hel\000oo",

		// µ => μ
		"\u00b5": "\u03bc",

		// okKK => okkk
		"ok\u212AK": "okkk",

		// Σ => σ
		"\u03a3": "\u03c3",

		// AͅΣ => aισ
		"A\u0345\u03a3": "a\u1fbe\u03c3",
	}

	funcs := map[string]func(str string) string{
		"CaseFoldStringFast":   CaseFoldStringFast,
		"CaseFoldStringLowMem": CaseFoldStringLowMem,
	}

	for funcName, f := range funcs {
		t.Run("with function "+funcName, func(t *testing.T) {
			for input, expected := range inputToExpected {
				if actual := f(input); actual != expected {
					t.Errorf("For input: %v; expected output: %v VS actual: %v", input, expected, actual)

					if again := f(actual); again != actual {
						t.Errorf("should be idempotent, %v != %v", actual, again)
					}
				}
			}
		})
	}
}
