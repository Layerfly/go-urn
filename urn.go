package urn

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const MaxURNLength = 255

var entityRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{1,31}$`)

// InvalidURNError is returned when a URN string is malformed.
type InvalidURNError struct {
	Message string
}

func (e *InvalidURNError) Error() string {
	return e.Message
}

// attrPair preserves insertion order of attributes.
type attrPair struct {
	Key   string
	Value string
}

// URN represents a parsed Uniform Resource Name.
type URN struct {
	Entity     string
	ID         string
	attributes []attrPair
}

// Attributes returns a copy of the attributes as a map.
func (u *URN) Attributes() map[string]string {
	m := make(map[string]string, len(u.attributes))
	for _, p := range u.attributes {
		m[p.Key] = p.Value
	}
	return m
}

// String returns the composed URN string.
func (u *URN) String() string {
	s, _ := compose(u.Entity, u.ID, u.attributes)
	return s
}

// CreateUUID generates a URN with a new UUID as the identifier.
func CreateUUID(entity string) string {
	id := uuid.New().String()
	s, _ := Compose(entity, id)
	return s
}

// Compose constructs a URN string from the given components.
func Compose(entity, id string, attrs ...map[string]string) (string, error) {
	var pairs []attrPair
	if len(attrs) > 0 && attrs[0] != nil {
		for k, v := range attrs[0] {
			pairs = append(pairs, attrPair{Key: k, Value: v})
		}
	}
	return compose(entity, id, pairs)
}

func compose(entity, id string, pairs []attrPair) (string, error) {
	if entity == "" || id == "" {
		return "", &InvalidURNError{Message: "Cannot compose URN: 'entity' and 'id' are required"}
	}

	safeEntity := url.PathEscape(entity)
	safeID := url.PathEscape(id)
	var b strings.Builder
	b.WriteString("urn:")
	b.WriteString(safeEntity)
	b.WriteString(":")
	b.WriteString(safeID)

	for _, p := range pairs {
		b.WriteString(":")
		b.WriteString(url.PathEscape(p.Key))
		b.WriteString(":")
		b.WriteString(url.PathEscape(p.Value))
	}

	result := b.String()
	if len(result) > MaxURNLength {
		return "", &InvalidURNError{
			Message: fmt.Sprintf("Composed URN is too long (%d chars, max %d)", len(result), MaxURNLength),
		}
	}
	return result, nil
}

// Parse deconstructs a URN string into its components.
func Parse(urnStr string) (*URN, error) {
	if !strings.HasPrefix(strings.ToLower(urnStr), "urn:") {
		return nil, &InvalidURNError{Message: "Invalid URN: Must start with the 'urn:' scheme"}
	}
	content := urnStr[4:]
	parts := strings.Split(content, ":")

	if len(parts) < 2 {
		return nil, &InvalidURNError{Message: "Invalid URN: Missing entity or ID component"}
	}

	entity := parts[0]
	id := parts[1]
	if entity == "" || id == "" {
		return nil, &InvalidURNError{Message: "Invalid URN: Entity or ID is empty"}
	}

	rest := parts[2:]
	if len(rest)%2 != 0 {
		return nil, &InvalidURNError{Message: "Invalid URN: Attribute key without value"}
	}

	var attrs []attrPair
	for i := 0; i < len(rest); i += 2 {
		key := rest[i]
		value := rest[i+1]
		if key == "" || value == "" {
			return nil, &InvalidURNError{
				Message: fmt.Sprintf("Invalid URN: Attribute %s missing value", key),
			}
		}
		attrs = append(attrs, attrPair{Key: key, Value: value})
	}

	return &URN{Entity: entity, ID: id, attributes: attrs}, nil
}

// Entity extracts the entity from a URN string.
func Entity(urnStr string) (string, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return "", err
	}
	return u.Entity, nil
}

// ID extracts the identifier from a URN string.
func ID(urnStr string) (string, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return "", err
	}
	return u.ID, nil
}

// Value retrieves the value for a specific attribute key.
// Returns the value, whether it was found, and any parse error.
func Value(urnStr, key string) (string, bool, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return "", false, err
	}
	for _, p := range u.attributes {
		if p.Key == key {
			return p.Value, true, nil
		}
	}
	return "", false, nil
}

// IsValid checks whether a string is a valid URN.
func IsValid(urnStr string) bool {
	if urnStr == "" || len(urnStr) > MaxURNLength {
		return false
	}
	u, err := Parse(urnStr)
	if err != nil {
		return false
	}
	if !entityRegex.MatchString(u.Entity) {
		return false
	}
	return true
}

// AddAttribute appends or updates an attribute in the URN.
func AddAttribute(urnStr, key, value string) (string, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return "", err
	}
	safeKey := url.PathEscape(key)
	safeValue := url.PathEscape(value)

	found := false
	for i, p := range u.attributes {
		if p.Key == safeKey {
			u.attributes[i].Value = safeValue
			found = true
			break
		}
	}
	if !found {
		u.attributes = append(u.attributes, attrPair{Key: safeKey, Value: safeValue})
	}
	return compose(u.Entity, u.ID, u.attributes)
}

// RemoveAttribute removes an attribute by key from the URN.
func RemoveAttribute(urnStr, key string) (string, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return "", err
	}
	filtered := make([]attrPair, 0, len(u.attributes))
	for _, p := range u.attributes {
		if p.Key != key {
			filtered = append(filtered, p)
		}
	}
	u.attributes = filtered
	return compose(u.Entity, u.ID, u.attributes)
}

// GetAllAttributes returns all key-value attribute pairs from a URN.
func GetAllAttributes(urnStr string) (map[string]string, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return nil, err
	}
	return u.Attributes(), nil
}

// Vendor is a convenience method that extracts the "vendor" attribute.
func Vendor(urnStr string) (string, bool, error) {
	return Value(urnStr, "vendor")
}

// Normalize lowercases the entity and re-composes the URN.
func Normalize(urnStr string) (string, error) {
	u, err := Parse(urnStr)
	if err != nil {
		return "", err
	}
	return compose(strings.ToLower(u.Entity), u.ID, u.attributes)
}
