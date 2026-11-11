package greenapi

import "github.com/green-api/whatsapp-api-client-golang-v2/voip"

// VoIP returns the Python-compatible call API. Keep the returned object and use
// Connect to create an independently owned callsRtc connection.
func (c *GreenAPI) VoIP() *voip.Voip {
	return voip.New(voip.Config{Host: c.APIURL, IDInstance: c.IDInstance, APITokenInstance: c.APITokenInstance})
}
