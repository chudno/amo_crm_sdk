package amo_crm_sdk

import (
	"testing"
)

func TestVersion(t *testing.T) {
	// Поскольку константы нельзя изменять, мы просто проверим
	// что формат версии соответствует ожидаемому формату
	
	// Получаем текущую версию
	got := Version()

	// Ожидаемая версия на основе текущих констант
	expectedWithoutSuffix := "1.0.0"
	expectedWithSuffix := "1.0.0-"

	// Проверяем соответствие версии ожиданиям
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
