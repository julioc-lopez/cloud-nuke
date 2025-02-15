package resources

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gruntwork-io/cloud-nuke/telemetry"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gruntwork-io/cloud-nuke/config"
	"github.com/stretchr/testify/require"
)

type mockedDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	DescribeTableOutputMap map[string]dynamodb.DescribeTableOutput
	ListTablesOutput       dynamodb.ListTablesOutput
	DeleteTableOutput      dynamodb.DeleteTableOutput
}

func (m mockedDynamoDB) ListTablesPages(input *dynamodb.ListTablesInput, fn func(*dynamodb.ListTablesOutput, bool) bool) error {
	fn(&m.ListTablesOutput, true)
	return nil
}

func (m mockedDynamoDB) DescribeTable(input *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	output := m.DescribeTableOutputMap[*input.TableName]
	return &output, nil
}

func (m mockedDynamoDB) DeleteTable(input *dynamodb.DeleteTableInput) (*dynamodb.DeleteTableOutput, error) {
	return &m.DeleteTableOutput, nil
}

func TestDynamoDB_GetAll(t *testing.T) {
	telemetry.InitTelemetry("cloud-nuke", "")
	t.Parallel()

	testName1 := "table1"
	testName2 := "table2"
	now := time.Now()
	ddb := DynamoDB{
		Client: mockedDynamoDB{
			ListTablesOutput: dynamodb.ListTablesOutput{
				TableNames: []*string{
					aws.String(testName1),
					aws.String(testName2),
				},
			},
			DescribeTableOutputMap: map[string]dynamodb.DescribeTableOutput{
				testName1: {
					Table: &dynamodb.TableDescription{
						TableName:        aws.String(testName1),
						CreationDateTime: aws.Time(now),
					},
				},
				testName2: {
					Table: &dynamodb.TableDescription{
						TableName:        aws.String(testName2),
						CreationDateTime: aws.Time(now.Add(1)),
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
					TimeAfter: aws.Time(now.Add(-1)),
				}},
			expected: []string{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			names, err := ddb.getAll(config.Config{
				DynamoDB: tc.configObj,
			})
			require.NoError(t, err)
			require.Equal(t, tc.expected, aws.StringValueSlice(names))
		})
	}
}

func TestDynamoDb_NukeAll(t *testing.T) {
	telemetry.InitTelemetry("cloud-nuke", "")
	t.Parallel()

	ddb := DynamoDB{
		Client: mockedDynamoDB{
			DeleteTableOutput: dynamodb.DeleteTableOutput{},
		},
	}

	err := ddb.nukeAll([]*string{aws.String("table1"), aws.String("table2")})
	require.NoError(t, err)
}
