package main

import (
	"testing"
)

func TestIsDirectoryExcluded1(t *testing.T) {
	excludedDirectories := []string{"hola", "chao"}
	if isDirectoryExcluded("/home/david/buenas", excludedDirectories) {
		t.FailNow()
	}
}

func TestIsDirectoryExcluded2(t *testing.T) {
	excludedDirectories := []string{"/home/david/buena", "/home/david/buenachao"}
	if !isDirectoryExcluded("/home/david/buenas", excludedDirectories) {
		t.FailNow()
	}
}
func TestIsDirectoryExcluded3(t *testing.T) {
	excludedDirectories := []string{"/home/david/buena", "/home/david/buenas"}
	if !isDirectoryExcluded("/home/david/buenas", excludedDirectories) {
		t.FailNow()
	}
}
