package courier

import (
	"regexp"
	"strings"

	"github.com/nyaruka/gocommon/urns"
)

// Ireland stores WhatsApp BSUIDs on the whatsapp: scheme (phone OR CC.alphanumeric), matching
// RapidPro main / gocommon v1.75.7. Upstream gocommon v1.83.0 (courier 26.2) only allows digits
// on whatsapp: and moved BSUID to the bsuid: scheme. These helpers restore the Ireland contract
// without pinning gocommon back.
var whatsAppBSUIDPathRegex = regexp.MustCompile(`^[A-Z]{2}\.[a-zA-Z0-9]{1,128}$`)

func init() {
	origValidate := urns.WhatsApp.Validate
	urns.WhatsApp.Validate = func(path string) bool {
		if origValidate != nil && origValidate(path) {
			return true
		}
		return isWhatsAppBSUIDPath(path)
	}
	urns.WhatsApp.Normalize = func(path string) string {
		return normalizeWhatsAppBSUIDPath(path)
	}
}

func normalizeWhatsAppBSUIDPath(path string) string {
	path = strings.TrimSpace(path)
	if len(path) > 2 && path[2] == '.' {
		return strings.ToUpper(path[:2]) + path[2:]
	}
	return path
}

func isWhatsAppBSUIDPath(path string) bool {
	return whatsAppBSUIDPathRegex.MatchString(normalizeWhatsAppBSUIDPath(path))
}

// IsWhatsAppBSUID reports whether the URN is a WhatsApp business-scoped user ID
// (Ireland whatsapp:CC.xxx, or the upstream bsuid: scheme).
func IsWhatsAppBSUID(u urns.URN) bool {
	switch u.Scheme() {
	case urns.WhatsApp.Prefix, urns.BSUID.Prefix:
		return isWhatsAppBSUIDPath(u.Path())
	default:
		return false
	}
}
