package sqlpath

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Returns a list of SQL files from given paths.
func TestReturnsListOfSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/file1.sql", "testdata/file2.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{"testdata/file1.sql", "testdata/file2.sql"}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v, %v", expected, result, cmp.Diff(expected, result))
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestReturnsNilListWhenNoSQLFilesFound(t *testing.T) {
	// Arrange
	paths := []string{"testdata/extra.txt"}

	// Act
	result, err := Glob(paths)
	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v, %v", expected, result, cmp.Diff(expected, result))
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestIgnoresHiddenFilesWhenSearchingForSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/.hidden.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestIgnoresNonSQLFilesWhenSearchingForSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/extra.txt"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestExcludesSQLFilesEndingWithDownSQLWhenSearchingForSQLFiles(t *testing.T) {
	// Arrange
	paths := []string{"testdata/file1.sql", "testdata/file3.down.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{"testdata/file1.sql"}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestReturnsErrorWhenPathDoesNotExist(t *testing.T) {
	// Arrange
	paths := []string{"non_existent_path"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	} else {
		expectedError := fmt.Errorf("path error: stat non_existent_path: no such file or directory")
		if !cmp.Equal(err.Error(), expectedError.Error()) {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}
	}
}

func TestReturnsErrorWhenDirectoryCannotBeRead(t *testing.T) {
	// Arrange
	paths := []string{"testdata/unreadable"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	} else {
		expectedError := fmt.Errorf("path error: stat testdata/unreadable: no such file or directory")
		if !cmp.Equal(err.Error(), expectedError.Error()) {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}
	}
}

func TestDoesNotIncludesSQLFilesWithUppercaseExtension(t *testing.T) {
	// Arrange
	paths := []string{"testdata/file4.SQL"}

	// Act
	result, err := Glob(paths)

	// Assert
	var expected []string
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestNotIncludesHiddenFilesAnyPath(t *testing.T) {
	// Arrange
	paths := []string{
		"./testdata/.hiddendir/file1.sql", // pass
		"./testdata/.hidden.sql",          // skip
	}

	// Act
	result, err := Glob(paths)

	// Assert
	expectedAny := [][]string{
		{"./testdata/.hiddendir/file1.sql"},
		{"testdata/.hiddendir/file1.sql"},
	}

	match := false
	for _, expected := range expectedAny {
		if cmp.Equal(result, expected) {
			match = true
			break
		}
	}
	if !match {
		t.Errorf("Expected any of %v, but got %v", expectedAny, result)
	}

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestFollowSymlinks(t *testing.T) {
	// Arrange
	paths := []string{"testdata/symlink", "testdata/file1.symlink.sql"}

	// Act
	result, err := Glob(paths)

	// Assert
	expected := []string{
		"testdata/symlink/file1.sql",
		"testdata/symlink/file1.symlink.sql",
		"testdata/symlink/file2.sql",
		"testdata/file1.symlink.sql",
	}
	if !cmp.Equal(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}

func TestGlobPattern(t *testing.T) {
	// Arrange
	tests := []struct {
		pattern  string
		expected []string
	}{
		{
			pattern: "testdata/glob/*/queries",
			expected: []string{
				"testdata/glob/sub1/queries/file1.sql",
				"testdata/glob/sub2/queries/file2.sql",
				"testdata/glob/sub3/queries/file3.sql",
				"testdata/glob/sub3/queries/file4.sql",
			},
		},
		{
			pattern: "testdata/glob/sub3/queries/file?.sql",
			expected: []string{
				"testdata/glob/sub3/queries/file3.sql",
				"testdata/glob/sub3/queries/file4.sql",
			},
		},
		{
			pattern: "testdata/glob/sub3/queries/file[1-5].sql",
			expected: []string{
				"testdata/glob/sub3/queries/file3.sql",
				"testdata/glob/sub3/queries/file4.sql",
			},
		},
	}

	for _, test := range tests {
		// Act
		result, err := Glob([]string{test.pattern})

		// Assert
		if !cmp.Equal(result, test.expected) {
			t.Errorf("Pattern %v: Expected %v, but got %v", test.pattern, test.expected, result)
		}
		if err != nil {
			t.Errorf("Pattern %v: Expected no error, but got %v", test.pattern, err)
		}
	}
}
