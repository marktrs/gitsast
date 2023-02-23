package analyzer_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/queue/task/analyzer"
	mocks "github.com/marktrs/gitsast/testutil/mocks"
	modelMock "github.com/marktrs/gitsast/testutil/mocks/model"
	"github.com/stretchr/testify/suite"

	analyzerMock "github.com/marktrs/gitsast/testutil/mocks/queue/analyzer"
)

// TestAnalyzerTestSuite is the external unit tests set for AnalyzerTestSuite.
type AnalyzerTestSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	report  *modelMock.MockIReportRepo
	repo    *modelMock.MockIRepositoryRepo
	rule    *modelMock.MockIRuleRepo
	git     *analyzerMock.MockIClient
	testApp *mocks.TestApp

	scanner  *analyzerMock.MockScanner
	detector *analyzerMock.MockDetector
	analyzer analyzer.IAnalyzeTask
}

func TestAnalyzerTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzerTestSuite))
}

func (suite *AnalyzerTestSuite) SetupTest() {
	var err error
	suite.ctrl = gomock.NewController(suite.T())
	suite.report = modelMock.NewMockIReportRepo(suite.ctrl)
	suite.repo = modelMock.NewMockIRepositoryRepo(suite.ctrl)
	suite.rule = modelMock.NewMockIRuleRepo(suite.ctrl)
	suite.git = analyzerMock.NewMockIClient(suite.ctrl)
	suite.detector = analyzerMock.NewMockDetector(suite.ctrl)
	suite.scanner = analyzerMock.NewMockScanner(suite.ctrl)

	suite.testApp = mocks.StartTestApp(context.Background())
	suite.analyzer, err = analyzer.NewAnalyzer(
		suite.testApp.App,
		suite.repo,
		suite.report,
		suite.rule,
		suite.git,
		suite.detector,
		suite.scanner,
	)
	suite.NoError(err)
}

func (suite *AnalyzerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *AnalyzerTestSuite) TestAnalyze() {
	suite.repo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(&model.Repository{
		ID: "fake-repo-uuid",
	}, nil)
	suite.report.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(&model.Report{
		ID: "fake-report-uuid",
	}, nil)
	suite.rule.EXPECT().GetAll(gomock.Any()).Return(nil, nil)
	suite.git.EXPECT().GetPathsFromRemoteURL(gomock.Any(), gomock.Any()).Return(nil, nil)
	suite.report.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, nil).MaxTimes(2)
	suite.scanner.EXPECT().ScanFilesForIssues(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	err := suite.analyzer.Analyze("fake-uuid")
	suite.NoError(err)
}

func (suite *AnalyzerTestSuite) TestAnalyzeError() {
	suite.report.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
	err := suite.analyzer.Analyze("fake-uuid")
	suite.Error(err)
	suite.EqualError(err, "sql: no rows in result set")
}
