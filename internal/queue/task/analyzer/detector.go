package analyzer

import (
	"regexp"

	"github.com/marktrs/gitsast/internal/model"
)

// detectIssueLocation - detect issue location
type Detector interface {
	DetectIssueLocation(fragment Fragment, rule *model.Rule) []*model.Issue
}

type detector struct{}

func NewDetector() Detector {
	return &detector{}
}

func (d *detector) DetectIssueLocation(fragment Fragment, rule *model.Rule) []*model.Issue {
	issues := make([]*model.Issue, 0)
	regex := regexp.MustCompile(rule.Keyword)
	matchIndices := regex.FindAllStringIndex(fragment.Raw, -1)

	for _, matchIndex := range matchIndices {
		loc := location(fragment, matchIndex)

		if matchIndex[1] > loc.endLineIndex {
			loc.endLineIndex = matchIndex[1]
		}

		issues = append(issues, &model.Issue{
			RuleID: model.GetFormattedRuleId(rule.ID),
			Location: model.Location{
				Path: fragment.FilePath,
				Line: uint64(loc.startLine),
			},
			Description: rule.Description,
			Severity:    rule.Severity.String(),
			Keyword:     rule.Keyword,
		})
	}

	return issues
}
