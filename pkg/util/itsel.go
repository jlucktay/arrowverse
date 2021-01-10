/*
Copyright © 2021 James Lucktaylor <jlucktay@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package util

import "fmt"

// IteratingSelector is a helper to get us through those pesky 'td' selectors.
type IteratingSelector struct {
	selectorFmt string
	tdOffset    int
}

// NewIteratingSelector currently has hard-coded values because we only use it in one loop.
func NewIteratingSelector() *IteratingSelector {
	return &IteratingSelector{
		selectorFmt: "td:nth-of-type(%d)",
		tdOffset:    0,
	}
}

func (is *IteratingSelector) String() string {
	return fmt.Sprintf(is.selectorFmt, is.tdOffset)
}

// Current will return the iterator with its current value.
func (is *IteratingSelector) Current() string {
	return is.String()
}

// Next will first increment the value, and then return the iterator.
func (is *IteratingSelector) Next() string {
	is.tdOffset++

	return is.String()
}
