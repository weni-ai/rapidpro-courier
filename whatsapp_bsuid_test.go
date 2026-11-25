package courier_test

import (
	"testing"

	"github.com/nyaruka/courier/v26"
	"github.com/nyaruka/gocommon/urns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsWhatsAppBSUID(t *testing.T) {
	bsuid := urns.URN("whatsapp:US.IRLBSUID001")
	phone := urns.URN("whatsapp:5511987654321")
	parent := urns.URN("whatsapp:US.ENT.11815799212886844830")
	legacy := urns.URN("bsuid:US.A1B2C3")

	assert.True(t, courier.IsWhatsAppBSUID(bsuid))
	assert.True(t, courier.IsWhatsAppBSUID(legacy))
	assert.False(t, courier.IsWhatsAppBSUID(phone))
	assert.False(t, courier.IsWhatsAppBSUID(parent))

	require.NoError(t, bsuid.Validate())
	require.NoError(t, phone.Validate())
	assert.Error(t, parent.Validate())

	created, err := urns.New(urns.WhatsApp, "US.IRLBSUID001")
	require.NoError(t, err)
	assert.Equal(t, "whatsapp:US.IRLBSUID001", string(created))
}
