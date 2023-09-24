package topic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ahmetson/common-lib/data_type/key_value"
)

type (
	// Id is the string representation of
	// the topic or topic filter
	Id string
	// Topic defines the smartcontract path in SDS
	Topic struct {
		Organization string `json:"org,omitempty"`
		Project      string `json:"proj,omitempty"`
		NetworkId    string `json:"net,omitempty"`
		Group        string `json:"group,omitempty"`
		Name         string `json:"name,omitempty"`
	}
)

var allPaths = []string{"org", "proj", "net", "group", "name"}

// Id from the topic to the Id,
// If one of the parameters is missing,
// then it will return an empty string
//
// Doesn't matter which topic level user want's get
// All topic should be equal
func (t *Topic) Id() Id {
	str := ""

	if len(t.Name) > 0 {
		str = "name-" + t.Name
	}
	if len(t.Group) > 0 {
		str = "group-" + t.Group + "." + str
	}

	if len(t.NetworkId) > 0 {
		str = "net-" + t.NetworkId + "." + str
	}

	if len(t.Project) > 0 {
		str = "proj-" + t.Project + "." + str
	}

	if len(t.Organization) > 0 {
		str = "org-" + t.Organization + "." + str
	}

	return Id(str)
}

// Level Calculates the level
// From the bottom level to up.
// If it's empty, then it returns 0
func (t *Topic) Level() uint8 {
	if len(t.Name) > 0 {
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
	topic := Topic{
		Organization: "",
		Project:      "",
		NetworkId:    "",
		Group:        "",
		Name:         "",
	}
	err := parameters.Interface(&topic)
	if err != nil {
		return nil, fmt.Errorf("failed to convert json into topic")
	}

	return &topic, nil
}

func isPathName(name string) bool {
	return name == "org" || name == "proj" || name == "net" || name == "group" || name == "name"
}

// the name should be valid literal
// it's only alphanumeric characters, and _
func isLiteral(val string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9 _]*$`).MatchString(val)
}

func (t *Topic) setPath(path string, val string) {
	switch path {
	case "org":
		t.Organization = val
	case "proj":
		t.Project = val
	case "net":
		t.NetworkId = val
	case "group":
		t.Group = val
	case "name":
		t.Name = val
	}
}

func (t *Topic) getValue(path string) string {
	switch path {
	case "org":
		return t.Organization
	case "proj":
		return t.Project
	case "net":
		return t.NetworkId
	case "group":
		return t.Group
	case "name":
		return t.Name
	}
	return ""
}

// Has the given paths or not. If not, then
// return an error.
//
// If the path argument has an unsupported path name, then that will be skipped
func (t *Topic) Has(paths ...string) bool {
	for _, path := range paths {
		if !isPathName(path) {
			return false
		}
		if len(t.getValue(path)) == 0 {
			return false
		}
	}

	return true
}

// ValidateMissingLevel The topic paths are in the order.
// The order is called level.
// If the bottom level's value is given, then the top
// level's parameters should be given too.
//
// Make sure that the upper level parameter is set.
func (t *Topic) ValidateMissingLevel(level uint8) error {
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
		return t.ValidateMissingLevel(OrganizationLevel)
	case NetworkIdLevel:
		if len(t.NetworkId) == 0 {
			return fmt.Errorf("missing network id")
		}
		return t.ValidateMissingLevel(ProjectLevel)
	case GroupLevel:
		if len(t.Group) == 0 {
			return fmt.Errorf("missing group")
		}
		return t.ValidateMissingLevel(NetworkIdLevel)
	case SmartcontractLevel:
		if len(t.Name) == 0 {
			return fmt.Errorf("missing smartcontract")
		}
		return t.ValidateMissingLevel(GroupLevel)
	default:
		return fmt.Errorf("unsupported level")
	}
}

// Validate the topic for empty values, for valid names.
// The topic parameters should be defined as literals in popular programming languages.
// Finally, the path of topic if it's converted to the TopicString should be valid as well.
//
// That means, if a user wants to create a topic to access t.Project, then its upper parent
// the t.Organization should be defined as well.
// But other topic parameters could be left as empty.
func (t *Topic) Validate() error {
	_, err := t.Id().Unmarshal()
	if err != nil {
		return fmt.Errorf("failed to validate: %w", err)
	}

	return nil
}

// Unmarshal This method converts Topic Id to the Topic Struct.
//
// The topic string is provided in the following string format:
//
//	`org-<organization>.proj-<project>.net-<network_id>.group-<group>.name-<name>`
//
// ----------------------
//
// Rules
//
//   - Order of the path names does not matter: org-<org>.proj-<proj> == proj-<proj>.org-<org>
//   - The values between `<` and `>` are literals and should return true by `isLiteral(literal)` function
func (id Id) Unmarshal() (Topic, error) {
	parts := strings.Split(string(id), ".")
	length := len(parts)
	if length < 2 {
		return Topic{}, fmt.Errorf("%s should have atleast 2 parts divided by ';'", id)
	}

	if length > 5 {
		return Topic{}, fmt.Errorf("%s should have at most 6 parts divided by ';'", id)
	}

	t := Topic{}

	for i, part := range parts {
		keyValue := strings.Split(part, "-")
		if len(keyValue) != 2 {
			return Topic{}, fmt.Errorf("part[%d] is %s, it can't be divided to two elements by '-'", i, part)
		}

		if !isPathName(keyValue[0]) {
			return Topic{}, fmt.Errorf("part[%d] isPathName(%s) false", i, keyValue[0])
		}

		if !isLiteral(keyValue[1]) {
			return Topic{}, fmt.Errorf("part[%d] ('%s') isLiteral(%v) false", i, keyValue[0], keyValue[1])
		}

		t.setPath(keyValue[0], keyValue[1])
	}

	return t, nil
}

// Only keep the given paths, and the rest are removed.
// In any case of error, the id returns itself.
// The error silently ignored.
func (id Id) Only(paths ...string) Id {
	topic, err := id.Unmarshal()
	if err != nil {
		return id
	}

	for _, path := range allPaths {
		keep := false
		for _, pathToKeep := range paths {
			if path == pathToKeep {
				keep = true
				break
			}
		}

		if !keep {
			topic.setPath(path, "")
		}
	}

	return topic.Id()
}

// Has the given paths or not. If not, then
// return an error.
//
// If the path argument has an unsupported path name, then that will be skipped
func (id Id) Has(paths ...string) bool {
	topic, err := id.Unmarshal()
	if err != nil {
		return false
	}
	return topic.Has(paths...)
}

const OrganizationLevel uint8 = 1  // only organization.
const ProjectLevel uint8 = 2       // only organization and project.
const NetworkIdLevel uint8 = 3     // only organization, project and, network id.
const GroupLevel uint8 = 4         // only organization and project, network id and group.
const SmartcontractLevel uint8 = 5 // smartcontract level path, till the smartcontract of the smartcontract
