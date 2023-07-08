package topic

import (
	"testing"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/stretchr/testify/suite"
)

// Test creation
//   - from parameters
//   - from json
//   - from string
//     topic filter string to topic string
//     should fail
//
// compare the level (for each nesting) against constants
//
// Test the string creation
// for each level
type TestTopicFilterSuite struct {
	suite.Suite
	topic       *Filter
	topicString Id
}

// SetupTest
// Setup checks the New() functions
// Setup checks ToMap() functions
func (suite *TestTopicFilterSuite) SetupTest() {
	sample := Filter{
		Organizations:  []string{"seascape"},
		Projects:       []string{"sds-core"},
		NetworkIds:     []string{"1", "56", "imx"},
		Groups:         []string{"test-suite"},
		Smartcontracts: []string{"TestErc20"},
		Events:         []string{"Transfer"},
	}
	topicString := AsTopicString(`o:seascape;p:sds-core;n:1,56,imx;g:test-suite;s:TestErc20;e:Transfer`)

	suite.topic = &sample
	suite.topicString = topicString

	suite.Require().Equal(topicString, sample.String())
}

func (suite *TestTopicFilterSuite) TestKvParameterParsing() {
	// empty kv, there is no test_filter
	kv := key_value.Empty()
	_, err := NewFromKeyValueParameter(kv)
	suite.Require().Error(err)

	// nil parameter in the key_value
	kv.Set("topic_filter", nil)
	_, err = NewFromKeyValueParameter(kv)
	suite.Require().Error(err)

	// topic filter is not a key value
	kv.Set("topic_filter", []string{"hello"})
	_, err = NewFromKeyValueParameter(kv)
	suite.Require().Error(err)

	// the map format is wrong
	kv.Set("topic_filter", map[interface{}]interface{}{"hello": "world"})
	_, err = NewFromKeyValueParameter(kv)
	suite.Require().Error(err)

	// converting empty topic filter should be fine
	kv.Set("topic_filter", map[string]interface{}{})
	_, err = NewFromKeyValueParameter(kv)
	suite.Require().NoError(err)

	// TopicFilter is not a struct
	expected := Filter{
		Organizations:  []string{"seascape"},
		Projects:       []string{"sds-core"},
		NetworkIds:     []string{"1", "56", "imx"},
		Groups:         []string{"test-suite"},
		Smartcontracts: []string{"TestErc20"},
		Events:         []string{"Transfer"},
	}
	kv.Set("topic_filter", expected)
	_, err = NewFromKeyValueParameter(kv)
	suite.Require().Error(err)

	// filter with the parameters should be fine
	sampleKv, err := key_value.NewFromInterface(expected)
	suite.Require().NoError(err)

	kv.Set("topic_filter", sampleKv)
	fromKv, err := NewFromKeyValueParameter(kv)
	suite.Require().NoError(err)

	kv.Set("topic_filter", sampleKv.Map())
	fromMap, err := NewFromKeyValueParameter(kv)
	suite.Require().NoError(err)

	suite.Require().EqualValues(expected, *fromKv)
	suite.Require().EqualValues(expected, *fromMap)
}

func (suite *TestTopicFilterSuite) TestToString() {
	empty := Filter{}
	topicString := empty.String()
	suite.Require().Empty(topicString)

	expected := Filter{
		Organizations:  []string{"seascape"},
		Projects:       []string{"sds-core"},
		NetworkIds:     []string{"1", "56", "imx"},
		Groups:         []string{"test-suite"},
		Smartcontracts: []string{"TestErc20"},
		Events:         []string{"Transfer"},
	}
	topicString = expected.String()
	expectedTopicString := Id(`o:seascape;p:sds-core;n:1,56,imx;g:test-suite;s:TestErc20;e:Transfer`)
	suite.Require().EqualValues(expectedTopicString, topicString)

	// some parameters are missing
	expected = Filter{
		Organizations:  []string{"seascape"},
		Projects:       []string{"sds-core"},
		Groups:         []string{"test-suite"},
		Smartcontracts: []string{"TestErc20"},
		Events:         []string{"Transfer"},
	}
	topicString = expected.String()
	expectedTopicString = `o:seascape;p:sds-core;g:test-suite;s:TestErc20;e:Transfer`
	suite.Require().EqualValues(expectedTopicString, topicString)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTopicFilter(t *testing.T) {
	suite.Run(t, new(TestTopicFilterSuite))
}
