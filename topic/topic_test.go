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
type TestTopicSuite struct {
	suite.Suite
	topic       Topic
	topicString String
}

// SetupTests
// Setup checks the New() functions
// Setup checks ToMap() functions
func (suite *TestTopicSuite) SetupTest() {
	sample := Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		NetworkId:     "1",
		Group:         "test-suite",
		Smartcontract: "TestErc20",
		Event:         "Transfer",
	}
	topicString := AsTopicString(`o:seascape;p:sds-core;n:1;g:test-suite;s:TestErc20;e:Transfer`)

	suite.topic = sample
	suite.topicString = topicString

	suite.Require().Equal(topicString, sample.String(FullLevel))
}

func (suite *TestTopicSuite) TestStringParse() {
	newTopic, err := ParseString(suite.topicString)
	suite.Require().NoError(err)
	suite.Require().EqualValues(suite.topic, newTopic)

	// additional parameter in the topic string should fail
	topicString := AsTopicString(`o:seascape;p:sds-core;n:1;g:test-suite;s:TestErc20;e:Transfer;m:transfer`)
	_, err = ParseString(topicString)
	suite.Require().Error(err)

	// case sensitive
	topicString = AsTopicString(`O:seascape;p:sds-core;n:1;g:test-suite;s:TestErc20;e:Transfer`)
	_, err = ParseString(topicString)
	suite.Require().Error(err)

	// additional semicolon should fail
	topicString = AsTopicString(`o:seascape;p:sds-core;n:1;g:test-suite;s:TestErc20;e:Transfer;`)
	_, err = ParseString(topicString)
	suite.Require().Error(err)

	// missing the one of the paths
	// if the event is given, then all previous levels
	// should be given too.
	// missing "network_id"
	topicString = AsTopicString(`o:seascape;p:sds-core;g:test-suite;s:TestErc20;e:Transfer`)
	_, err = ParseString(topicString)
	suite.Require().Error(err)

	// value of the topic path is not a literal
	// it has not required tokens.
	topicString = AsTopicString(`o:seascape:network;p:sds-core;n:1;g:test-suite;s:TestErc20;e:Transfer`)
	_, err = ParseString(topicString)
	suite.Require().Error(err)
}

func (suite *TestTopicSuite) TestParsingJson() {
	kv := key_value.Empty().
		Set("o", "seascape").
		Set("p", "sds-core").
		Set("n", "1").
		Set("g", "test-suite").
		Set("s", "TestErc20").
		Set("e", "Transfer")

	newTopic, err := ParseJSON(kv)
	suite.Require().NoError(err)
	suite.Require().EqualValues(suite.topic, *newTopic)

	// changing the orders doesn't affect the topic
	kv = key_value.Empty().
		Set("o", "seascape").
		Set("n", "1").
		Set("p", "sds-core").
		Set("g", "test-suite").
		Set("s", "TestErc20").
		Set("e", "Transfer")

	newTopic, err = ParseJSON(kv)
	suite.Require().NoError(err)
	suite.Require().EqualValues(suite.topic, *newTopic)

	// additional parameter in the topic string
	// should succeed, but the value will be missed
	kv.Set("m", "transfer")
	_, err = ParseJSON(kv)
	suite.Require().NoError(err)

	// setting with the empty parameter should fail
	// empty group
	invalidKv := key_value.Empty().
		Set("o", "seascape").
		Set("p", "sds").
		Set("n", "1").
		Set("g", "").
		Set("s", "TestErc20").
		Set("e", "Transfer")
	_, err = ParseJSON(invalidKv)
	suite.Require().Error(err)

	// case-sensitive
	// Group name is given as 'G', should be 'g'
	invalidKv = key_value.Empty().
		Set("o", "seascape").
		Set("p", "sds").
		Set("n", "1").
		Set("G", "test-suite").
		Set("s", "TestErc20").
		Set("e", "Transfer")
	_, err = ParseJSON(invalidKv)
	suite.Require().Error(err)

	// missing the one of the paths
	// if the event is given, then all previous levels
	// should be given too.
	// missing "group"
	invalidKv = key_value.Empty().
		Set("o", "seascape").
		Set("p", "sds").
		Set("n", "1").
		Set("s", "TestErc20").
		Set("e", "Transfer")
	_, err = ParseJSON(invalidKv)
	suite.Require().Error(err)
}

func (suite *TestTopicSuite) TestToString() {
	topic := Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		NetworkId:     "1",
		Group:         "test-suite",
		Smartcontract: "TestErc20",
		Event:         "Transfer",
	}

	topicString := topic.String(0)
	suite.Require().Empty(topicString)

	topicString = topic.String(7)
	suite.Require().Empty(topicString)

	expectedTopicString := String(`o:seascape;p:sds-core;n:1;g:test-suite;s:TestErc20;e:Transfer`)
	topicString = topic.String(6)
	suite.Require().EqualValues(expectedTopicString, topicString)

	expectedTopicString = `o:seascape;p:sds-core;n:1;g:test-suite;s:TestErc20`
	topicString = topic.String(5)
	suite.Require().EqualValues(expectedTopicString, topicString)

	expectedTopicString = `o:seascape;p:sds-core;n:1;g:test-suite`
	topicString = topic.String(4)
	suite.Require().EqualValues(expectedTopicString, topicString)

	expectedTopicString = `o:seascape;p:sds-core;n:1`
	topicString = topic.String(3)
	suite.Require().EqualValues(expectedTopicString, topicString)

	expectedTopicString = `o:seascape;p:sds-core`
	topicString = topic.String(2)
	suite.Require().EqualValues(expectedTopicString, topicString)

	expectedTopicString = `o:seascape`
	topicString = topic.String(1)
	suite.Require().EqualValues(expectedTopicString, topicString)

	expectedTopicString = `o:seascape`
	suite.Require().EqualValues(expectedTopicString, topicString)

	topic = Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		NetworkId:     "",
		Group:         "test-suite",
		Smartcontract: "TestErc20",
		Event:         "Transfer",
	}
	topicString = topic.String(FullLevel)
	suite.Require().Empty(topicString)

	// NetworkId is empty, the upper root exists
	// But all topic should be valid
	topicString = topic.String(ProjectLevel)
	suite.Require().Empty(topicString)

	topic = Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		Group:         "test-suite",
		Smartcontract: "TestErc20",
		Event:         "Transfer",
	}
	topicString = topic.String(FullLevel)
	suite.Require().Empty(topicString)

	topic = Topic{
		Organization:  "seascape",
		Project:       "sds-core",
		NetworkId:     "1",
		Group:         "test-suite",
		Smartcontract: "TestErc20",
	}
	// the topic is FullLevel,
	// but we try to get full path
	// it should fail
	topicString = topic.String(FullLevel)
	suite.Require().Empty(topicString)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTopic(t *testing.T) {
	suite.Run(t, new(TestTopicSuite))
}
