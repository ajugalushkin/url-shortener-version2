package osexitcheck

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitCheck(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "./...")
}
