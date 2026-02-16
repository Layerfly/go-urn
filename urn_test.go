package urn

import (
	"regexp"
	"strings"
	"testing"
)

func TestEntity(t *testing.T) {
	entity, err := Entity("urn:orders:1234")
	if err != nil {
		t.Fatal(err)
	}
	if entity != "orders" {
		t.Errorf("expected orders, got %s", entity)
	}
}

func TestID(t *testing.T) {
	id, err := ID("urn:orders:1234")
	if err != nil {
		t.Fatal(err)
	}
	if id != "1234" {
		t.Errorf("expected 1234, got %s", id)
	}
}

func TestValueExisting(t *testing.T) {
	val, found, err := Value("urn:orders:1234:vendorCode:abcd", "vendorCode")
	if err != nil {
		t.Fatal(err)
	}
	if !found || val != "abcd" {
		t.Errorf("expected abcd, got %s (found=%v)", val, found)
	}
}

func TestValueProductSKU(t *testing.T) {
	val, found, err := Value("urn:product:65b2713b1267994147953b27:vendor:foo:sku:999", "sku")
	if err != nil {
		t.Fatal(err)
	}
	if !found || val != "999" {
		t.Errorf("expected 999, got %s", val)
	}
}

func TestValueNonExisting(t *testing.T) {
	_, found, err := Value("urn:product:65b2713b1267994147953b27:vendor:foo:sku:123", "foo")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Error("expected not found")
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid("urn:orders:1234") {
		t.Error("expected valid")
	}
}

func TestIsValidMissingEntity(t *testing.T) {
	if IsValid("urn::1234") {
		t.Error("expected invalid")
	}
}

func TestIsValidInvalidFormat(t *testing.T) {
	if IsValid("invalid:orders:1234") {
		t.Error("expected invalid")
	}
}

func TestIsValidTooLong(t *testing.T) {
	longURN := "urn:long" + strings.Repeat(":a", 250)
	if IsValid(longURN) {
		t.Error("expected invalid for long URN")
	}
}

func TestAddAttribute(t *testing.T) {
	updated, err := AddAttribute("urn:orders:1234", "customer", "john-doe")
	if err != nil {
		t.Fatal(err)
	}
	val, found, _ := Value(updated, "customer")
	if !found || val != "john-doe" {
		t.Errorf("expected john-doe, got %s", val)
	}
}

func TestRemoveAttribute(t *testing.T) {
	updated, err := RemoveAttribute("urn:orders:1234:customer:john-doe", "customer")
	if err != nil {
		t.Fatal(err)
	}
	_, found, _ := Value(updated, "customer")
	if found {
		t.Error("expected attribute removed")
	}
}

func TestGetAllAttributes(t *testing.T) {
	attrs, err := GetAllAttributes("urn:orders:1234:customer:john-doe:status:pending")
	if err != nil {
		t.Fatal(err)
	}
	if attrs["customer"] != "john-doe" || attrs["status"] != "pending" {
		t.Errorf("unexpected attributes: %v", attrs)
	}
}

func TestCreateUUID(t *testing.T) {
	result := CreateUUID("session")
	matched, _ := regexp.MatchString(`^urn:session:[a-f0-9-]{36}$`, result)
	if !matched {
		t.Errorf("unexpected UUID URN format: %s", result)
	}
}

func TestNormalize(t *testing.T) {
	normalized, err := Normalize("URN:EXAMPLE:Animal:Ferret:Nose")
	if err != nil {
		t.Fatal(err)
	}
	if normalized != "urn:example:Animal:Ferret:Nose" {
		t.Errorf("expected urn:example:Animal:Ferret:Nose, got %s", normalized)
	}
}

func TestVendor(t *testing.T) {
	val, found, err := Vendor("urn:orders:1234:vendor:amazon")
	if err != nil {
		t.Fatal(err)
	}
	if !found || val != "amazon" {
		t.Errorf("expected amazon, got %s", val)
	}
}

func TestVendorMissing(t *testing.T) {
	_, found, err := Vendor("urn:orders:1234")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Error("expected not found")
	}
}

func TestParseMalformed(t *testing.T) {
	_, err := Parse("invalidURN")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "Must start with the 'urn:' scheme") {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestParseMissingEntityOrID(t *testing.T) {
	_, err := Parse("urn::")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "Entity or ID is empty") {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestParseUnmatchedKeyValue(t *testing.T) {
	_, err := Parse("urn:orders:1234:status")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "Attribute key without value") {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestCompose(t *testing.T) {
	result, err := Compose("order", "12345", map[string]string{
		"vendor": "amazon",
		"status": "shipped",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result, "urn:order:12345") {
		t.Errorf("unexpected compose result: %s", result)
	}
	// Verify round-trip
	u, err := Parse(result)
	if err != nil {
		t.Fatal(err)
	}
	if u.Entity != "order" || u.ID != "12345" {
		t.Errorf("unexpected parse: %+v", u)
	}
	attrs := u.Attributes()
	if attrs["vendor"] != "amazon" || attrs["status"] != "shipped" {
		t.Errorf("unexpected attributes: %v", attrs)
	}
}

func TestComposeEmptyEntity(t *testing.T) {
	_, err := Compose("", "123")
	if err == nil {
		t.Fatal("expected error for empty entity")
	}
}

func TestURNString(t *testing.T) {
	u, err := Parse("urn:orders:1234:status:pending")
	if err != nil {
		t.Fatal(err)
	}
	if u.String() != "urn:orders:1234:status:pending" {
		t.Errorf("unexpected String(): %s", u.String())
	}
}
