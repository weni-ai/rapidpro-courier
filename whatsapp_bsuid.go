package courier

import (
	"regexp"

	"github.com/nyaruka/gocommon/urns"
)

// gocommon v1.75.7 validates BSUID on whatsapp: but does not export IsWhatsAppBSUID (v1.87.0+).
var whatsAppBSUIDRegex = regexp.MustCompile(`^[A-Z]{2}\.[a-zA-Z0-9]{1,128}$`)

// IsWhatsAppBSUID reports whether the URN is a WhatsApp business-scoped user ID
// (whatsapp: with a CC.ALPHANUMERIC path rather than a phone number).
func IsWhatsAppBSUID(u urns.URN) bool {
	return u.Scheme() == urns.WhatsApp.Prefix && whatsAppBSUIDRegex.MatchString(u.Path())
}
