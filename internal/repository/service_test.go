package repository_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/repository"

	"github.com/stretchr/testify/suite"

	mocks "github.com/marktrs/gitsast/testutil/mocks"
	modelMock "github.com/marktrs/gitsast/testutil/mocks/model"
	queueMock "github.com/marktrs/gitsast/testutil/mocks/queue"
)

type ServiceTestSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	report  *modelMock.MockIReportRepo
	repo    *modelMock.MockIRepositoryRepo
	queue   *queueMock.MockHandler
	testApp *mocks.TestApp

	service repository.IService
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.report = modelMock.NewMockIReportRepo(suite.ctrl)
	suite.repo = modelMock.NewMockIRepositoryRepo(suite.ctrl)
	suite.queue = queueMock.NewMockHandler(suite.ctrl)
	suite.testApp = mocks.StartTestApp(context.Background())
	suite.testApp.App.SetQueue(suite.queue)

	suite.service = repository.NewService(
		suite.testApp.App, suite.repo, suite.report)
}

func (suite *ServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *ServiceTestSuite) TestGetById() {
	suite.repo.EXPECT().GetById(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err := suite.service.GetById(context.Background(), "fake-uuid")
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestList() {
	suite.repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err := suite.service.List(context.Background(), &model.RepositoryFilter{})
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestAdd() {
	suite.repo.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err := suite.service.
		Add(
			context.Background(),
			&repository.AddRepositoryRequest{
				Name:      "test",
				RemoteURL: "https://github.com/test/test.git",
			})
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestUpdate() {
	suite.repo.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err := suite.service.
		Update(
			context.Background(),
			"fake-uuid",
			&repository.UpdateRepositoryRequest{
				Name:      "test",
				RemoteURL: "https://github.com/test/test.git",
			})
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestRemove() {
	suite.repo.EXPECT().Remove(gomock.Any(), gomock.Any()).Return(nil)
	err := suite.service.Remove(context.Background(), "fake-uuid")
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestCreateReport() {
	suite.repo.EXPECT().
		GetById(gomock.Any(), gomock.Any()).Return(&model.Repository{}, nil)
	suite.report.EXPECT().
		GetByRepoId(gomock.Any(), gomock.Any()).Return(nil, nil)
	suite.report.EXPECT().
		Add(gomock.Any(), gomock.Any()).Return(&model.Report{}, nil)
	suite.queue.EXPECT().AddTask(gomock.Any()).Return(nil)
	suite.report.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, nil)

	_, err := suite.service.CreateReport(context.Background(), "fake-uuid")
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestGetReport() {
	suite.report.EXPECT().GetByRepoId(gomock.Any(), gomock.Any()).Return(nil, nil)
	_, err := suite.service.GetReportByRepoId(context.Background(), "fake-uuid")
	suite.NoError(err)
}

func (suite *ServiceTestSuite) TestAddRepositoryRequestValidation() {
	suite.repo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		AnyTimes()
	{

		cases := []struct {
			name    string
			body    *repository.AddRepositoryRequest
			errMsg  string
			wantErr bool
		}{
			{
				name: "valid input",
				body: &repository.AddRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "https://github.com/test/test.git",
				},
				wantErr: false,
			},
			{
				name: "valid input 2",
				body: &repository.AddRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "http://github.com/test/test.git",
				},
				wantErr: false,
			},
			{
				name: "valid input 3",
				body: &repository.AddRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "git@github.com/test/test.git",
				},
				wantErr: false,
			},
			{
				name: "blank name",
				body: &repository.AddRepositoryRequest{
					Name:      "",
					RemoteURL: "https://github.com/test/test.git",
				},
				errMsg:  "Field validation for 'Name' failed on the 'required' tag",
				wantErr: true,
			},
			{
				name: "blank remote url",
				body: &repository.AddRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "",
				},
				errMsg:  "Field validation for 'RemoteURL' failed on the 'required' tag",
				wantErr: true,
			},
			{
				name: "long name",
				body: &repository.AddRepositoryRequest{
					Name: `Lorem ipsum dolor, sit amet consectetur adipisicing 
					elit. Repudiandaesunt nihil tenetur, eius quaerat, fugit b
					eatae asperiores`,
					RemoteURL: "https://github.com/test/test.git",
				},
				errMsg:  "Field validation for 'Name' failed on the 'max' tag",
				wantErr: true,
			},
			{
				name: "long repo",
				body: &repository.AddRepositoryRequest{
					Name: "lorem",
					RemoteURL: `https://github.com/test/Loremipsumdolorsitametc
					onsecteturadipisicingelitRepudiandaesuntnihiltenetureiusqua
					eratfugitbeataeasperiores.git`,
				},
				errMsg:  "Field validation for 'RemoteURL' failed on the 'max' tag",
				wantErr: true,
			},
		}

		for _, c := range cases {

			_, err := suite.service.Add(context.Background(), c.body)

			if c.wantErr {
				suite.Assert().ErrorContains(err, c.errMsg)
			} else {
				suite.NoError(err)
			}
		}

	}
}

func (suite *ServiceTestSuite) TestUpdateRepositoryRequestValidation() {
	suite.repo.EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	{

		cases := []struct {
			name    string
			body    *repository.UpdateRepositoryRequest
			errMsg  string
			wantErr bool
		}{
			{
				name: "valid input",
				body: &repository.UpdateRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "https://github.com/test/test.git",
				},
				wantErr: false,
			},
			{
				name: "valid input 2",
				body: &repository.UpdateRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "http://github.com/test/test.git",
				},
				wantErr: false,
			},
			{
				name: "valid input 3",
				body: &repository.UpdateRepositoryRequest{
					Name:      "lorem",
					RemoteURL: "git@github.com/test/test.git",
				},
				wantErr: false,
			},
			{
				name: "long name",
				body: &repository.UpdateRepositoryRequest{
					Name: `Lorem ipsum dolor, sit amet consectetur adipisi
					cing elit. Repudiandae sunt nihil tenetur, eius quaerat, fu
					git beatae asperiores`,
					RemoteURL: "https://github.com/test/test.git",
				},
				errMsg:  "Field validation for 'Name' failed on the 'max' tag",
				wantErr: true,
			},
			{
				name: "long repo",
				body: &repository.UpdateRepositoryRequest{
					Name: "lorem",
					RemoteURL: `https://github.com/test/Loremipsumdolorsitametc
					onsecteturadipisicingelit.Repudiandaesuntnihiltenetureiusqu
					aeratfugitbeataeasperiores.git`,
				},
				errMsg:  "Field validation for 'RemoteURL' failed on the 'max' tag",
				wantErr: true,
			},
		}

		for _, c := range cases {

			err := suite.service.Update(
				context.Background(),
				"fake-uuid",
				c.body,
			)

			if c.wantErr {
				suite.Assert().ErrorContains(err, c.errMsg)
			} else {
				suite.NoError(err)
			}
		}

	}
}
