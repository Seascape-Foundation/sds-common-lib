package topic

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/key_value"
)

type (
	// String is the string representation of
	// the topic or topic filter
	String string
	// Topic defines the smartcontract path in SDS
	Topic struct {
		Organization  string `json:"o,omitempty"`
		Project       string `json:"p,omitempty"`
		NetworkId     string `json:"n,omitempty"`
		Group         string `json:"g,omitempty"`
		Smartcontract string `json:"s,omitempty"`
		Event         string `json:"e,omitempty"`
	}
)

// AsTopicString Converts the string to the TopicString
func AsTopicString(topicString string) String {
	return String(topicString)
}

// Convert the TopicString to string
func (topicString String) String() string {
	return string(topicString)
}

// Converts the topic to the TopicString
// If one of the parameters is missing,
// then it will return an empty string
//
// Doesn't matter which topic level user want's get
// All topic should be equal
func (t *Topic) String(level uint8) String {
	if level < 1 || level > 6 {
		return ""
	}
	err := t.validateMissingLevel(t.Level())
	if err != nil {
		return ""
	}
	// we request inner level,
	// when it's not given in the topic
	if t.Level() < level {
		return ""
	}

	str := ""

	if level >= 1 {
		str += "o:" + t.Organization
	}
	if level >= 2 {
		str += ";p:" + t.Project
	}
	if level >= 3 {
		str += ";n:" + t.NetworkId
	}
	if level >= 4 {
		str += ";g:" + t.Group
	}
	if level >= 5 {
		str += ";s:" + t.Smartcontract
	}

	if level >= 6 {
		str += ";e:" + t.Event
	}

	return String(str)
}

// Level Calculates the level
// From the bottom level to up.
// If it's an empty, then it returns 0
func (t *Topic) Level() uint8 {
	if len(t.Event) > 0 {
		return FullLevel
	}
	if len(t.Smartcontract) > 0 {
		return SmartcontractLevel
	}
	if len(t.Group) > 0 {
		return GroupLevel
	}
	if len(t.NetworkId) > 0 {
		return NetworkIdLevel
	}
	if len(t.Project) > 0 {
		return ProjectLevel
	}
	if len(t.Organization) > 0 {
		return OrganizationLevel
	}

	return 0
}

// ParseJSON Parse JSON into the Topic
func ParseJSON(parameters key_value.KeyValue) (*Topic, error) {
	organization, err := parameters.GetString("o")
	if err != nil {
		return nil, fmt.Errorf("parameters.GetString(`o`): %w", err)
	}
	if len(organization) == 0 {
		return nil, errors.New("empty 'o' parameter")
	}
	project, err := parameters.GetString("p")
	if err != nil {
		return nil, fmt.Errorf("parameters.GetString(`p`): %w", err)
	}
	if len(project) == 0 {
		return nil, errors.New("empty 'p' parameter")
	}
	topic := Topic{
		Organization:  organization,
		Project:       project,
		NetworkId:     "",
		Group:         "",
		Smartcontract: "",
		Event:         "",
	}

	networkId, err := parameters.GetString("n")
	if err == nil {
		topic.NetworkId = networkId
	}

	group, err := parameters.GetString("g")
	if err == nil {
		topic.Group = group
	}

	smartcontract, err := parameters.GetString("s")
	if err == nil {
		topic.Smartcontract = smartcontract
	}

	event, err := parameters.GetString("e")
	if err == nil {
		topic.Event = event
	}

	err = topic.validateMissingLevel(topic.Level())
	if err != nil {
		return nil, fmt.Errorf("missing upper level: %w", err)
	}

	return &topic, nil
}

func isPathName(name string) bool {
	return name == "o" || name == "p" || name == "n" || name == "g" || name == "s" || name == "e"
}

// the name should be valid literal
// it's only alphanumeric characters, - and _
func isLiteral(val string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9 _-]*$`).MatchString(val)
}

func (t *Topic) setPath(path string, val string) error {
	switch path {
	case "o":
		if len(t.Organization) > 0 {
			return fmt.Errorf("the duplicate organization path name. already set as " + t.Organization)
		} else {
			t.Organization = val
		}
	case "p":
		if len(t.Project) > 0 {
			return fmt.Errorf("the duplicate project path name. already set as " + t.Project)
		} else {
			t.Project = val
		}
	case "n":
		if len(t.NetworkId) > 0 {
			return fmt.Errorf("the duplicate network id path name. already set as " + t.NetworkId)
		} else {
			t.NetworkId = val
		}
	case "g":
		if len(t.Group) > 0 {
			return fmt.Errorf("the duplicate group path name. already set as " + t.Group)
		} else {
			t.Group = val
		}
	case "s":
		if len(t.Smartcontract) > 0 {
			return fmt.Errorf("the duplicate smartcontract path name. already set as " + t.Smartcontract)
		} else {
			t.Smartcontract = val
		}
	case "e":
		if len(t.Event) > 0 {
			return fmt.Errorf("the duplicate event path name. already set as " + t.Event)
		} else {
			t.Event = val
		}
	}

	return nil
}

// The topic paths are in the order.
// The order is called level.
// If the bottom level's value is given, then the top
// level's parameters should be given too.
//
// Make sure that the upper level parameter is set.
func (t *Topic) validateMissingLevel(level uint8) error {
	switch level {
	case OrganizationLevel:
		if len(t.Organization) == 0 {
			return fmt.Errorf("missing organization")
		}
		return nil
	case ProjectLevel:
		if len(t.Project) == 0 {
			return fmt.Errorf("missing project")
		}
		return t.validateMissingLevel(OrganizationLevel)
	case NetworkIdLevel:
		if len(t.NetworkId) == 0 {
			return fmt.Errorf("missing network id")
		}
		return t.validateMissingLevel(ProjectLevel)
	case GroupLevel:
		if len(t.Group) == 0 {
			return fmt.Errorf("missing group")
		}
		return t.validateMissingLevel(NetworkIdLevel)
	case SmartcontractLevel:
		if len(t.Smartcontract) == 0 {
			return fmt.Errorf("missing smartcontract")
		}
		return t.validateMissingLevel(GroupLevel)
	case FullLevel:
		if len(t.Event) == 0 {
			return fmt.Errorf("missing event")
		}
		return t.validateMissingLevel(SmartcontractLevel)
	default:
		return fmt.Errorf("unsupported level")
	}
}

// Validate the topic for empty values, for valid names.
// The topic parameters should be defined as literals in popular programming languages.
// Finally, the path of topic if it's converted to the TopicString should be valid as well.
//
// That means, if user wants to create a topic to access t.Project, then it's upper parent
// the t.Organization should be defined as well.
// But other topic parameters could be left as empty.
func (t *Topic) Validate() error {
	level := t.Level()
	str := t.String(level)
	_, err := ParseString(str)
	if err != nil {
		return fmt.Errorf("failed to validate: %w", err)
	}

	return nil
}

// ParseString This method converts Topic String to the Topic Struct.
//
// The topic string is provided in the following string format:
//
//	`o:<organization>;p:<project>;n:<network id>;g:<group>;s:<smartcontract>;m:<method>`
//	`o:<organization>;p:<project>;n:<network id>;g:<group>;s:<smartcontract>;e:<event>`
//
// ----------------------
//
// Rules
//
//   - the topic string can have either `method` or `event` but not both at the same time.
//   - Topic string should contain at least 'organization' and 'project'
//   - Order of the path names does not matter: o:org;p:proj == p:proj;o:org
//   - The values between `<` and `>` are literals and should return true by `is_literal(literal)` function
func ParseString(topicString String) (Topic, error) {
	parts := strings.Split(topicString.String(), ";")
	length := len(parts)
	if length < 2 {
		return Topic{}, fmt.Errorf("%s should have atleast 2 parts divided by ';'", topicString)
	}

	if length > 6 {
		return Topic{}, fmt.Errorf("%s should have at most 6 parts divided by ';'", topicString)
	}

	t := Topic{}

	for i, part := range parts {
		keyValue := strings.Split(part, ":")
		if len(keyValue) != 2 {
			return Topic{}, fmt.Errorf("part[%d] is %s, it can't be divided to two elements by ':'", i, part)
		}

		if !isPathName(keyValue[0]) {
			return Topic{}, fmt.Errorf("part[%d] is_path_name(%s) false", i, keyValue[0])
		}

		if !isLiteral(keyValue[1]) {
			return Topic{}, fmt.Errorf("part[%d] ('%s') is_literal(%v) false", i, keyValue[0], keyValue[1])
		}

		err := t.setPath(keyValue[0], keyValue[1])
		if err != nil {
			return t, fmt.Errorf("part[%d] set_path: %w", i, err)
		}
	}

	err := t.validateMissingLevel(t.Level())
	if err != nil {
		return Topic{}, fmt.Errorf("missing upper level: %w", err)
	}

	return t, nil
}

const OrganizationLevel uint8 = 1  // only organization.
const ProjectLevel uint8 = 2       // only organization and project.
const NetworkIdLevel uint8 = 3     // only organization, project and, network id.
const GroupLevel uint8 = 4         // only organization and project, network id and group.
const SmartcontractLevel uint8 = 5 // smartcontract level path, till the smartcontract of the smartcontract
const FullLevel uint8 = 6          // full topic path
