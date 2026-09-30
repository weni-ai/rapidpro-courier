package whatsapp_test

import (
	"testing"

	"github.com/nyaruka/courier/handlers/meta/whatsapp"
	"github.com/nyaruka/gocommon/urns"
	"github.com/stretchr/testify/assert"
)

func TestRecipientFields(t *testing.T) {
	tcs := []struct {
		urn               urns.URN
		expectedTo        string
		expectedRecipient string
	}{
		{urn: "whatsapp:250788123123", expectedTo: "250788123123", expectedRecipient: ""},
		{urn: "whatsapp:US.13491208655302741918", expectedTo: "", expectedRecipient: "US.13491208655302741918"},
	}

	for _, tc := range tcs {
		to, recipient := whatsapp.RecipientFields(tc.urn)
		assert.Equal(t, tc.expectedTo, to, "to mismatch for %s", tc.urn)
		assert.Equal(t, tc.expectedRecipient, recipient, "recipient mismatch for %s", tc.urn)
	}
}
