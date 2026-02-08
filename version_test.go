package amo_crm_sdk

import (
	"testing"
)

func TestVersion(t *testing.T) {
	got := Version()

	expectedWithoutSuffix := "0.2.0"
	expectedWithSuffix := "0.2.0-"

	if VersionSuffix == "" {
		if got != expectedWithoutSuffix {
			t.Errorf("Version() = %q, хотим %q", got, expectedWithoutSuffix)
		}
	} else {
		expectedWithSuffix += VersionSuffix
		if got != expectedWithSuffix {
			t.Errorf("Version() = %q, хотим %q", got, expectedWithSuffix)
		}
	}
}
