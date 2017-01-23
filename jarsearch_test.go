package main

import (
	"bytes"
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

func TestOutputMavenResults(t *testing.T) {
	var doc Doc
	doc.G = "groupIdX"
	doc.A = "artifactIdX"
	doc.LatestVersion = "1.0.0"

	var dependencies Dependencies
	dependencies.Response.Docs = []Doc{doc}

	out = new(bytes.Buffer)

	outputMavenResults(dependencies)

	actualResult := out.(*bytes.Buffer).String()
	expectedResult := "\n<dependency>\n<groupId>groupIdX</groupId>\n<artifactId>artifactIdX</artifactId>\n<version>1.0.0</version>\n</dependency>\n"
	assert.Equal(t, expectedResult, actualResult, "Maven output not as expected")
}

func TestOutputGradleResults(t *testing.T) {
	var doc Doc
	doc.Id = "groupIdX:artifactIdX"
	doc.LatestVersion = "1.0.0"

	var dependencies Dependencies
	dependencies.Response.Docs = []Doc{doc}

	out = new(bytes.Buffer)

	outputGradleResults(dependencies)

	actualResult := out.(*bytes.Buffer).String()
	expectedResult := "groupIdX:artifactIdX:1.0.0\n"
	assert.Equal(t, expectedResult, actualResult, "Gradle output not as expected")
}
