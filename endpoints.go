package hubspot

import "net/http"

// Endpoints represent specific calls to the Hubspot API and some meta data about them. These should largely be ignored
// by the consumer of the API
const (
	EndpointCreateContact = "endpointCreateContact"
	EndpointGetContact    = "endpointGetContact"
	EndpointDeleteContact = "endpointDeleteContact"

	EndpointCreateEventType = "endpointCreateEventType"
	EndpointDeleteEventType = "endpointDeleteEventType"
	EndpointCreateEvent     = "endpointCreateEvent"
)

type endpoint struct {
	Method       string
	Path         string
	RequireOAuth bool
	MockGoodHTTP int
	MockGood     map[string]interface{}
}

var endpoints = map[string]endpoint{
	// Contacts
	EndpointCreateContact: endpoint{
		Method:   http.MethodPost,
		Path:     "/contacts/v1/contact/createOrUpdate/email/:email",
		MockGood: nil,
	},
	EndpointGetContact: endpoint{
		Method:   http.MethodGet,
		Path:     "/contacts/v1/contact/email/:email/profile",
		MockGood: nil,
	},
	EndpointDeleteContact: endpoint{
		Method:   http.MethodDelete,
		Path:     "/contacts/v1/contact/vid/:vid",
		MockGood: nil,
	},
	// Events
	EndpointCreateEventType: endpoint{
		Method:       http.MethodPost,
		Path:         "/integrations/v1/:applicationID/timeline/event-types",
		MockGoodHTTP: http.StatusCreated,
		MockGood: map[string]interface{}{
			"id":             float64(123),
			"name":           "Test Event Type",
			"headerTemplate": "# Title for event {{id}}\nThis is an event for {{objectType}}",
			"detailTemplate": "This event happened on {{#formatDate timestamp}}{{/formatDate}}",
			"applicationId":  123,
			"objectType":     "CONTACT",
		},
	},
	EndpointDeleteEventType: endpoint{
		Method:       http.MethodDelete,
		Path:         "/integrations/v1/:applicationID/timeline/event-types/:eventTypeID",
		MockGoodHTTP: http.StatusNoContent,
		MockGood:     map[string]interface{}{},
	},
	EndpointCreateEvent: endpoint{
		Method:       http.MethodPut,
		Path:         "/integrations/v1/:applicationID/timeline/event",
		RequireOAuth: true,
		MockGoodHTTP: http.StatusNoContent,
		MockGood:     nil,
	},
}
