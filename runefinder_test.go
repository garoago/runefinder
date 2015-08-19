package main

import (
	"testing"

	"gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

type MySuite struct{}

var _ = check.Suite(&MySuite{})

func (s *MySuite) TestFindRunes(c *check.C) {
	index := map[string][]rune{
		"REGISTERED": []rune{0xAE},
	}

	tests := map[string][]rune{
		"registered": []rune{0xAE},
		"nonesuch":   []rune{},
	}
	for query, found := range tests {
		c.Assert(findRunes(query, index), check.DeepEquals, found)
	}
}
