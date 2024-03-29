package model

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zimmski/osutil"

	"github.com/symflower/eval-symflower-codegen-testing/language"
)

func TestModelSymflowerGenerateTestsForFile(t *testing.T) {
	type testCase struct {
		Name string

		ModelSymflower *ModelSymflower
		Language       language.Language

		RepositoryPath   string
		RepositoryChange func(t *testing.T, repositoryPath string)
		FilePath         string

		ExpectedError     error
		ExpectedErrorText string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			temporaryPath := t.TempDir()
			repositoryPath := filepath.Join(temporaryPath, filepath.Base(tc.RepositoryPath))
			require.NoError(t, osutil.CopyTree(tc.RepositoryPath, repositoryPath))

			if tc.RepositoryChange != nil {
				tc.RepositoryChange(t, repositoryPath)
			}

			if tc.ModelSymflower == nil {
				tc.ModelSymflower = &ModelSymflower{}
			}
			actualErr := tc.ModelSymflower.GenerateTestsForFile(repositoryPath, tc.FilePath)

			if tc.ExpectedError != nil {
				assert.ErrorIs(t, tc.ExpectedError, actualErr)
			} else if actualErr != nil || tc.ExpectedErrorText != "" {
				assert.ErrorContains(t, actualErr, tc.ExpectedErrorText)
			}

			require.NoError(t, tc.Language.Execute(repositoryPath))
		})
	}

	validate(t, &testCase{
		Name: "Go",

		Language: &language.LanguageGolang{},

		RepositoryPath: "../testdata/golang/plain/",
		FilePath:       "plain.go",
	})
}
