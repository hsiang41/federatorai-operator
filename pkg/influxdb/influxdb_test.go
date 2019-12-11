package influxdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getLocalClient() (Client, error) {
	return NewClient(Config{
		Address:  "https://localhost:8086",
		Username: "admin",
		Password: "adminpass",
	})
}

func TestListMeasurements(t *testing.T) {
	cli, err := getLocalClient()
	if err != nil {
		t.Error(err)
		return
	}

	type testCaseHave struct {
		db           string
		measurements []string
	}

	type testCase struct {
		have testCaseHave
		want []Measurement
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{
				db:           "alameda_cluster_status",
				measurements: []string{"pod", "namespace"},
			},
			want: []Measurement{},
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		actual, err := cli.ListMeasurements(testCase.have.db, testCase.have.measurements)
		assert.NoError(err)
		t.Logf("actual %+v", actual)
		assert.Equal(testCase.want, actual)
	}
}

func TestListMeasurementsInDatabase(t *testing.T) {
	cli, err := getLocalClient()
	if err != nil {
		t.Error(err)
		return
	}

	type testCase struct {
		have string
		want []string
	}

	testCases := []testCase{
		testCase{
			have: "alameda_cluster_status",
			want: []string{},
		},
	}

	assert := assert.New(t)
	for _, testCase := range testCases {
		actual, err := cli.ListMeasurementsInDatabase(testCase.have)
		assert.NoError(err)
		t.Logf("actual %+v", actual)
		assert.Equal(testCase.want, actual)
	}
}
