// Package topic defines the special kind of data type called topic and topic filter.
// The topics are replacing the smartcontract addresses in order to detect the smartcontract
// that user wants to interact with.
//
// For example, if the user wants to interact with Crowns cryptocurrency
// on Ethereum network, then user will need to know the Crowns ABI interface,
// as well as the smartcontract address.
//
// In SDS its replaced with the Topic.
// Define the topic something like:
//
//		topic := topic.Topic{
//			Organization: "seascape",
//	     	NetworkId: "1",
//	     	Name: "Crowns"
//		}
//
// For example, use the topic in SDK to read the data from indexer.
// Viola, we don't need to remember smartcontract address.
package topic

import (
	"fmt"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
)

// Filter unlike Topic can omit the parameters
// Allows to define list of smartcontract that match the topic filter.
//
// Which means users can interact with multiple smartcontract at once.
type Filter struct {
	Organizations  []string `json:"o,omitempty"`
	Projects       []string `json:"p,omitempty"`
	NetworkIds     []string `json:"n,omitempty"`
	Groups         []string `json:"g,omitempty"`
	Smartcontracts []string `json:"s,omitempty"`
}

// convert properties to string
func reduceProperties(properties []string) string {
	str := ""
	for i, v := range properties {
		if i != 0 {
			str += ","
		}
		str += v
	}

	return str
}

func (t *Filter) hasNestedLevel(level uint8) bool {
	switch level {
	case OrganizationLevel:
		if !t.hasNestedLevel(ProjectLevel) {
			return len(t.Organizations) != 0
		}
		return true
	case ProjectLevel:
		if !t.hasNestedLevel(NetworkIdLevel) {
			return len(t.Projects) != 0
		}
		return true
	case NetworkIdLevel:
		if !t.hasNestedLevel(GroupLevel) {
			return len(t.NetworkIds) != 0
		}
		return true
	case GroupLevel:
		if !t.hasNestedLevel(SmartcontractLevel) {
			return len(t.Groups) != 0
		}
		return true
	case SmartcontractLevel:
		return len(t.Smartcontracts) != 0
	}
	return false
}

// Convert the topic filter object to the topic filter string.
func (t *Filter) String() Id {
	str := ""
	if len(t.Organizations) > 0 {
		str += "o:" + reduceProperties(t.Organizations)
		if t.hasNestedLevel(OrganizationLevel) {
			str += ";"
		}
	}
	if len(t.Projects) > 0 {
		str += "p:" + reduceProperties(t.Projects)
		if t.hasNestedLevel(ProjectLevel) {
			str += ";"
		}
	}
	if len(t.NetworkIds) > 0 {
		str += "n:" + reduceProperties(t.NetworkIds)
		if t.hasNestedLevel(NetworkIdLevel) {
			str += ";"
		}
	}
	if len(t.Groups) > 0 {
		str += "g:" + reduceProperties(t.Groups)
		if t.hasNestedLevel(GroupLevel) {
			str += ";"
		}
	}
	if len(t.Smartcontracts) > 0 {
		str += "s:" + reduceProperties(t.Smartcontracts)
		if t.hasNestedLevel(SmartcontractLevel) {
			str += ";"
		}
	}

	return Id(str)
}

// NewFromKeyValueParameter extracts the "topic_filter" parameter from parameters.
func NewFromKeyValueParameter(parameters key_value.KeyValue) (*Filter, error) {
	topicFilterMap, err := parameters.GetKeyValue("topic_filter")
	if err != nil {
		return nil, fmt.Errorf("missing `topic_filter` parameter")
	}

	var filter Filter
	err = topicFilterMap.Interface(&filter)

	if err != nil {
		return nil, fmt.Errorf("failed to convert the value to TopicFilter: %w", err)
	}

	return &filter, nil
}
