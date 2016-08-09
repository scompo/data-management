package domain

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var testTime time.Time = time.Now()

func TestEncode(t *testing.T) {
	require := require.New(t)
	currentTime = func() time.Time {
		return testTime
	}
	testData := []struct {
		value    Project
		expected Project
	}{
		{
			Project{
				Name:        "testName",
			},
			Project{
				Name:         "testName",
				CreationDate: testTime,
			},
		},
	}
	for _, test := range testData {
		var buf bytes.Buffer
		test.value.Encode(&buf)
		expData, err := json.Marshal(test.expected)
		if err != nil {
			t.Errorf("%v", err)
		}
		require.Equal(string(expData[:])+"\n", buf.String())
	}
	currentTime = time.Now
}
