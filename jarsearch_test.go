package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func MakeRequestMock(url string) string {
	return "{\"response\":{\"numFound\": 12}}"
}

func TestSearch(t *testing.T) {
	MakeRequest = MakeRequestMock

	dependencies := search("http://localhost/mock")
	var expDep Dependencies
	expDep.Response.NumFound = 12
	assert.Equal(t, expDep, dependencies)
}
