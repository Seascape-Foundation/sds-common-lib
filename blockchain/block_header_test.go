package blockchain

import (
	"testing"

	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestBlockHeaderSuite struct {
	suite.Suite
	header BlockHeader
}

// Test setup (inproc, tcp and sub)
//	Along with the reconnect
// Test Requests (router, remote)
// Test the timeouts
// Test close (attempt to request)

// Todo test in-process and external types of controllers
// Todo test the business of the controller
// Make sure that Account is set to five
// before each test
func (suite *TestBlockHeaderSuite) SetupTest() {
	uintNumber := uint64(10)
	uintTimestamp := uint64(123)
	uintMap := key_value.Empty().
		Set("block_number", uintNumber).
		Set("block_timestamp", uintTimestamp)

	header, _ := NewHeader(uintNumber, uintTimestamp)
	mapHeader, err := NewHeaderFromKeyValueParameter(uintMap)
	suite.Require().NoError(err)
	number, _ := NewNumber(uintNumber)
	mapNumber, err := NewNumberFromKeyValueParameter(uintMap)
	suite.Require().NoError(err)
	timestamp, _ := NewTimestamp(uintTimestamp)
	mapTimestamp, err := NewTimestampFromKeyValueParameter(uintMap)
	suite.Require().NoError(err)

	suite.Require().Equal(header, mapHeader)
	suite.Require().Equal(header.Number.Value(), uintNumber)
	suite.Require().Equal(header.Number, number)
	suite.Require().Equal(header.Number, mapNumber)
	suite.Require().Equal(header.Timestamp.Value(), uintTimestamp)
	suite.Require().Equal(header.Timestamp, timestamp)
	suite.Require().Equal(header.Timestamp, mapTimestamp)

	suite.header = header

	// Missing both parameters
	emptyMap := key_value.Empty()
	_, err = NewHeaderFromKeyValueParameter(emptyMap)
	suite.Require().Error(err)

	// Missing timestamp, should fail
	noTimestampMap := key_value.Empty().
		Set("block_number", uintNumber)
	_, err = NewHeaderFromKeyValueParameter(noTimestampMap)
	suite.Require().Error(err)

	// Missing timestamp, should fail
	noNumberMap := key_value.Empty().
		Set("block_timestamp", uintTimestamp)
	_, err = NewHeaderFromKeyValueParameter(noNumberMap)
	suite.Require().Error(err)
}

func (suite *TestBlockHeaderSuite) TestValueChange() {
	number, _ := NewNumber(11)
	suite.header.Number.Increment()
	suite.Equal(number, suite.header.Number)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBlockHeader(t *testing.T) {
	suite.Run(t, new(TestBlockHeaderSuite))
}
