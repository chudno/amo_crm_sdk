// Информация о версии SDK
package amo_crm_sdk

import "fmt"

// Версия SDK
const (
	// VersionMajor - мажорный номер версии
	VersionMajor = 0
	// VersionMinor - минорный номер версии
	VersionMinor = 2
	// VersionPatch - патч версии
	VersionPatch = 0
	// VersionSuffix - суффикс версии (например, beta, rc)
	VersionSuffix = ""
)

// Version возвращает строку с полной версией SDK
func Version() string {
	version := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionSuffix != "" {
		version += "-" + VersionSuffix
	}
	return version
}
