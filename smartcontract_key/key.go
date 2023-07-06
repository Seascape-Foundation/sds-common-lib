// Package smartcontract_key defines the unique
// smartcontract id within SeascapeSDS
//
// The [smartcontract_key.Key] is composed of a string: network_id + "." + address
package smartcontract_key

import (
	"fmt"
	"strings"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
)

// Key of smartcontract composed of network id + "." + address.
// Any smartcontract should have one unique key.
type Key struct {
	NetworkId string `json:"network_id"`
	Address   string `json:"address"`
}

// KeyToTopicString is the equivalent of `map(smartcontract_key => topicString)`
type KeyToTopicString map[Key]string

// New key from networkId and address
func New(networkId string, address string) (Key, error) {
	key := Key{NetworkId: networkId, Address: address}
	err := key.Validate()
	if err != nil {
		return Key{}, fmt.Errorf("key.Validate: %w", err)
	}

	return key, nil
}

// NewFromKeyValue converts the parameters to Key
func NewFromKeyValue(parameters key_value.KeyValue) (Key, error) {
	var key Key
	err := parameters.Interface(&key)
	if err != nil {
		return Key{}, fmt.Errorf("failed to convert key-value to interface: %w", err)
	}

	err = key.Validate()
	if err != nil {
		return Key{}, fmt.Errorf("key.Validate: %w", err)
	}

	return key, nil
}

// NewFromString converts the string to Key
func NewFromString(s string) (Key, error) {
	str := strings.Split(s, ".")
	if len(str) != 2 {
		return Key{}, fmt.Errorf("string '%s' doesn't have two parts", s)
	}

	if len(str[0]) == 0 ||
		len(str[1]) == 0 {
		return Key{}, fmt.Errorf("missing parameter or empty parameter")
	}

	key := Key{NetworkId: str[0], Address: str[1]}
	err := key.Validate()
	if err != nil {
		return Key{}, fmt.Errorf("key.Validate: %w", err)
	}

	return key, nil
}

// Returns the key as a string
// `<network_id>.<address>`
func (k *Key) String() string {
	return k.NetworkId + "." + k.Address
}

func (k *Key) Validate() error {
	if len(k.NetworkId) == 0 {
		return fmt.Errorf("missing network id")
	}
	if len(k.Address) == 0 {
		return fmt.Errorf("missing address")
	}

	return nil
}
