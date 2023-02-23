package analyzer_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/queue/task/analyzer"
	analyzerMock "github.com/marktrs/gitsast/testutil/mocks/queue/analyzer"
	"github.com/stretchr/testify/assert"
)

func TestScanFilesForIssues(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		content       string
		rules         []*model.Rule
		expectedIssue []*model.Issue
		wantError     error
	}{
		{
			name:    "public key leak simple",
			path:    "pub.key",
			content: `xibcuvsdf: public_key=sbodufsdfin`,
			expectedIssue: []*model.Issue{
				{
					RuleID: "G001",
					Location: model.Location{
						Path: "pub.key",
						Line: 0,
					},
					Description: "Public key leak",
					Severity:    "LOW",
					Keyword:     `public_key`,
				},
			},
			rules: []*model.Rule{
				{
					ID:          1,
					Name:        "Public key leak",
					Keyword:     `public_key`,
					Description: "Public key leak",
					Severity:    model.Low,
				},
			},
		},
		{
			name:    "Private key leak simple",
			path:    "priv.key",
			content: `xibcuvsdf: private_key=sbodufsdfin`,
			expectedIssue: []*model.Issue{
				{
					RuleID: "G002",
					Location: model.Location{
						Path: "priv.key",
						Line: 0,
					},
					Description: "Private key leak",
					Severity:    "LOW",
					Keyword:     `private_key`,
				},
			},
			rules: []*model.Rule{
				{
					ID:          2,
					Name:        "Private key leak",
					Keyword:     `private_key`,
					Description: "Private key leak",
					Severity:    model.High,
				},
			},
		},
	}

	for _, tc := range testCases {
		prepareTestFiles(t, tc.path, tc.content)
		defer cleanupTestFiles(t, tc.path)

		ctrl := gomock.NewController(t)
		detector := analyzerMock.NewMockDetector(ctrl)
		scanner := analyzer.NewScanner(detector)

		detector.EXPECT().DetectIssueLocation(gomock.Any(), gomock.Any()).Return(tc.expectedIssue)

		issues, err := scanner.ScanFilesForIssues(tc.path, []string{tc.path}, tc.rules)
		if tc.wantError != nil {
			assert.EqualError(t, err, tc.wantError.Error())
			continue
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, tc.expectedIssue, issues)
	}
}

func prepareTestFiles(t *testing.T, path, content string) {
	_, err := os.Create(path)
	assert.NoError(t, err)
	assert.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

func cleanupTestFiles(t *testing.T, path string) {
	assert.NoError(t, os.Remove(path))
}
