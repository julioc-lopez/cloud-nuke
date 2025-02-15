package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/gruntwork-io/cloud-nuke/telemetry"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

type mockedECR struct {
	ecriface.ECRAPI
	DescribeRepositoriesPagesOutput ecr.DescribeRepositoriesOutput
	DeleteRepositoryOutput          ecr.DeleteRepositoryOutput
}

func (m mockedECR) DescribeRepositoriesPages(input *ecr.DescribeRepositoriesInput, callback func(*ecr.DescribeRepositoriesOutput, bool) bool) error {
	callback(&m.DescribeRepositoriesPagesOutput, true)
	return nil
}

func (m mockedECR) DeleteRepository(input *ecr.DeleteRepositoryInput) (*ecr.DeleteRepositoryOutput, error) {
	return &m.DeleteRepositoryOutput, nil
}

func TestECR_GetAll(t *testing.T) {
	telemetry.InitTelemetry("cloud-nuke", "")
	t.Parallel()

	testName1 := "test-repo1"
	testName2 := "test-repo2"
	now := time.Now()
	er := ECR{
		Client: &mockedECR{
			DescribeRepositoriesPagesOutput: ecr.DescribeRepositoriesOutput{
				Repositories: []*ecr.Repository{
					{
						RepositoryName: aws.String(testName1),
						CreatedAt:      aws.Time(now),
					},
					{
						RepositoryName: aws.String(testName2),
						CreatedAt:      aws.Time(now.Add(1)),
					},
				},
			},
		},
	}

	tests := map[string]struct {
		configObj config.ResourceType
		expected  []string
	}{
		"emptyFilter": {
			configObj: config.ResourceType{},
			expected:  []string{testName1, testName2},
		},
		"nameExclusionFilter": {
			configObj: config.ResourceType{
				ExcludeRule: config.FilterRule{
					NamesRegExp: []config.Expression{{
						RE: *regexp.MustCompile(testName1),
					}}},
			},
			expected: []string{testName2},
		},
		"timeAfterExclusionFilter": {
			configObj: config.ResourceType{
				ExcludeRule: config.FilterRule{
					TimeAfter: aws.Time(now),
				}},
			expected: []string{testName1},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			names, err := er.getAll(config.Config{
				ECRRepository: tc.configObj,
			})
			require.NoError(t, err)
			require.Equal(t, tc.expected, aws.StringValueSlice(names))
		})
	}

}

func TestECR_NukeAll(t *testing.T) {
	telemetry.InitTelemetry("cloud-nuke", "")
	t.Parallel()

	er := ECR{
		Client: &mockedECR{
			DeleteRepositoryOutput: ecr.DeleteRepositoryOutput{},
		},
	}

	err := er.nukeAll([]string{"test-repo1", "test-repo2"})
	require.NoError(t, err)
}
