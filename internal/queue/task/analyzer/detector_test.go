package analyzer

import (
	"testing"

	"github.com/marktrs/gitsast/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDetectIssueLocation(t *testing.T) {
	// NOTE: when certain issues are expected to occur,
	// the line numbers in the output will always be 0.
	// This is because line numbers are only added after the finding is created.
	testCases := []struct {
		name          string
		fragment      Fragment
		rule          *model.Rule
		expectedIssue []*model.Issue
	}{
		{
			name: "public key leak simple",
			fragment: Fragment{
				Raw:      `xibcuvsdf: public_key=sbodufsdfin`,
				FilePath: "tmp.txt",
			},
			rule: &model.Rule{
				ID:          1,
				Name:        "Public key leak",
				Keyword:     `public_key`,
				Description: "Public key leak",
				Severity:    model.Low,
			},
			expectedIssue: []*model.Issue{
				{
					RuleID: "G001",
					Location: model.Location{
						Path: "tmp.txt",
						Line: 0,
					},
					Description: "Public key leak",
					Severity:    "LOW",
					Keyword:     `public_key`,
				},
			},
		},
		{
			name: "Private key leak simple",
			fragment: Fragment{
				Raw:      `xibcuvsdf: private_key=sbodufsdfin`,
				FilePath: "tmp.txt",
			},
			rule: &model.Rule{
				ID:          2,
				Name:        "Private key leak",
				Keyword:     `private_key`,
				Description: "Private key leak",
				Severity:    model.High,
			},
			expectedIssue: []*model.Issue{
				{
					RuleID: "G002",
					Location: model.Location{
						Path: "tmp.txt",
						Line: 0,
					},
					Description: "Private key leak",
					Severity:    "HIGH",
					Keyword:     `private_key`,
				},
			},
		},
	}
	for _, tc := range testCases {
		d := NewDetector()
		issues := d.DetectIssueLocation(tc.fragment, tc.rule)
		assert.Equal(t, 1, len(issues), "Expected 1 issues, got %d", len(issues))
		assert.EqualValues(t, tc.expectedIssue, issues, "Expected issues does not match")
	}
}
