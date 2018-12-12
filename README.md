# gold

A simple golang lib to [case-fold](https://www.w3.org/International/wiki/Case_folding) UTF8 strings. Go-fold didn't sound quite as nice as `gold`!

## Why?

When you need to compare two UTF8 strings, [`strings.EqualFold`](https://golang.org/pkg/strings/#EqualFold) lets you test equality under Unicode case-folding, and that's great. But what if you want to do other comparisons? For example, what if you want to test whether one string is included in another one, modulo Unicode case-folding?

That's why this lib gives you functions to convert UTF-8 strings to their case-folded equivalents, similar to [python's `str.casefold()` function](https://docs.python.org/3/library/stdtypes.html#str.casefold).

## Installation

```
go get -u github.com/wk8/gold
```

Or use your favorite golang vendoring tool!

## Usage

The two main functions are:
```
func CaseFoldString(string) string
func CaseFoldBytes([]byte) []byte
```
They respectively convert a `string` or a `[]byte` to their case-folded representation.

For example:

```
package main

import (
	"fmt"
	"strings"

	"github.com/wk8/gold"
)

func main() {
	fmt.Println(gold.CaseFoldString("heLlo")) // => "hello"

	fmt.Println(strings.Contains(
		gold.CaseFoldString("Hey Σalμt toi"),
		gold.CaseFoldString("σalµT"))) // => true
}
```

If you're more concerned about memory usage than speed, you can also use the functionally equivalent
```
func CaseFoldStringLowMem(string) string
func CaseFoldBytesLowMem([]byte) []byte
```
functions instead (half the memory usage, but twice as slow).

Finally, please note that all of `gold`'s functions expect valid UTF8 strings as inputs, and do not verify that. If you need to validate your inputs, please use `utf8`'s functions [`Valid`](https://golang.org/pkg/unicode/utf8/#Valid) or [`ValidString`](https://golang.org/pkg/unicode/utf8/#ValidString).
