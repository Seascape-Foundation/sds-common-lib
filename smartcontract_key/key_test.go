package smartcontract_key

import (
	"testing"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestKeySuite struct {
	suite.Suite
	key Key
}

func (suite *TestKeySuite) SetupTest() {
	networkId := "1"
	address := "0x123"
	uintMap := key_value.Empty().
		Set("network_id", networkId).
		Set("address", address)

	key, _ := New(networkId, address)
	suite.Require().Equal(networkId, key.NetworkId)
	suite.Require().Equal(address, key.Address)

	mapKey, err := NewFromKeyValue(uintMap)
	suite.Require().NoError(err)
	suite.Require().Equal(networkId, mapKey.NetworkId)
	suite.Require().Equal(address, mapKey.Address)

	uintMap = key_value.Empty().
		Set("address", address)
	_, err = NewFromKeyValue(uintMap)
	suite.Require().Error(err)

	uintMap = key_value.Empty().
		Set("network_id", networkId)
	_, err = NewFromKeyValue(uintMap)
	suite.Require().Error(err)

	uintMap = key_value.Empty()
	_, err = NewFromKeyValue(uintMap)
	suite.Require().Error(err)

	uintMap = key_value.Empty().
		Set("network_id", networkId).
		Set("address", address).
		Set("additional_param", uint64(1))
	_, err = NewFromKeyValue(uintMap)
	suite.Require().NoError(err)

	suite.key = key
}

func (suite *TestKeySuite) TestToString() {
	keyString := "1.0x123"
	suite.Require().Equal(keyString, suite.key.String())

	key, err := NewFromString(keyString)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.key, key)

	noNetworkString := ".0x123"
	_, err = NewFromString(noNetworkString)
	suite.Require().Error(err)

	noAddress := "1."
	_, err = NewFromString(noAddress)
	suite.Require().Error(err)

	noParameters := "."
	_, err = NewFromString(noParameters)
	suite.Require().Error(err)

	empty := ""
	_, err = NewFromString(empty)
	suite.Require().Error(err)

	tooManyParameters := "1.0x1232.123213"
	_, err = NewFromString(tooManyParameters)
	suite.Require().Error(err)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestBlockHeader(t *testing.T) {
	suite.Run(t, new(TestKeySuite))
}
