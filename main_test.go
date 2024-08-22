package main

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

const TEXT = `
“The rain it raineth on the just
And also on the unjust fella;
But chiefly on the just, because
The unjust hath the just’s umbrella.”
― Charles Bowen

"It's not the voting that's democracy, it's the counting!"
- Tom Stoppard

"Comment is free but facts are on expenses."
- Tom Stoppard

"O, but they say the tongues of dying men
Enforce attention like deep harmony:
Where words are scarce, they are seldom spent in vain,
For they breathe truth that breathe their words in
pain."
- William Shakespeare

`

func TestWordCount(t *testing.T) {

	c := qt.New(t)

	input := strings.NewReader(TEXT)
	wc := NewWordCounter()
	c.Assert(88, qt.Equals, wc.Count(input))
}

func BenchmarkWordCount(b *testing.B) {

	input := strings.NewReader(TEXT)

	wc := NewWordCounter()

	for i := 0; i < b.N; i++ {
		_ = wc.Count(input)
	}
}
