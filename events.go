package hubspot

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// EventType represents an event type used for adding events to the timeline
type EventType struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
	// HeaderTemplate is used when displaying the event in the timeline and uses Handlebars.js templating
	HeaderTemplate string `json:"headerTemplate"`
	// DetailTemplate is used when displaying the event in the timeline and uses Handlebars.js templating
	DetailTemplate string `json:"detailTemplate"`
	ApplicationID  string `json:"applicationId"`
	// ObjectType should be one of `CONTACT` `COMPANY` `DEAL`
	ObjectType string `json:"objectType"`
}

// Event represents a timeline event
type Event struct {
	ID string `json:"id"`
	// The VID of the contact
	ObjectID       int64              `json:"objectId"`
	EventTypeID    int64              `json:"eventTypeId"`
	ExtraData      *map[string]string `json:"extraData,omitempty"`
	TimelineIFrame *EventIFrame       `json:"timelineIFrame,omitempty"`
}

// EventIFrame represent an event IFrame which will open an iframe in the display of the timeline event. All fields are needed.
type EventIFrame struct {
	LinkLabel   string `json:"linkLabel"`
	IFrameLabel string `json:"iframeLabel"`
	IFrameURI   string `json:"iframeUri"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

// CreateNewEventType creates a new event type for events on the timeline
//
// API Doc: https://developers.hubspot.com/docs/methods/timeline/create-event-type
func CreateNewEventType(input *EventType) error {
	input.ApplicationID = Config.HubspotApplicationID
	if input.Name == "" || input.ObjectType == "" || input.ApplicationID == "" {
		return APIError{
			HTTPCode:   http.StatusBadRequest,
			SystemCode: CodeEventTypeMissingData,
			Message:    "name, applicationID, and objectType are all required",
			Body:       nil,
		}
	}
	input.ObjectType = strings.ToUpper(input.ObjectType)
	if input.HeaderTemplate == "" {
		input.HeaderTemplate = "# Title for event {{id}}\nThis is an event for {{objectType}}"
	}
	if input.DetailTemplate == "" {
		input.DetailTemplate = "This event happened on {{#formatDate timestamp}}{{/formatDate}}"
	}

	ret, err := prepareCall(EndpointCreateEventType, map[string]string{}, input)

	if err != nil {
		if apiErr, apiErrOK := err.(APIError); apiErrOK {
			apiErr.SystemCode = CodeEventTypeCouldNotBeCreated
			return apiErr
		}
		return err
	}
	if body, bodyOK := ret.Body.(map[string]interface{}); bodyOK {
		if idF, idFOK := body["id"].(float64); idFOK {
			input.ID = int64(idF)
		}
	}

	return nil
}

// CreateOrUpdateEvent creates or update an event on the timeline. The ID should be provided externally. If it isn't, this func will
// generate a random ID based upon the current timestamp. This returns no data since the return header is a 204
//
// API Doc: https://developers.hubspot.com/docs/methods/timeline/create-or-update-event
func CreateOrUpdateEvent(input *Event) error {
	// check for some required fields
	if input.ObjectID == 0 || input.EventTypeID == 0 {
		return APIError{
			HTTPCode:   http.StatusBadRequest,
			SystemCode: CodeEventMissingData,
			Message:    "objectID and eventTypeID are required",
			Body:       nil,
		}
	}
	if input.ID == "" {
		rand.Seed(time.Now().UnixNano())
		input.ID = fmt.Sprintf("%d%d%d", rand.Int63(), time.Now().Unix(), input.ObjectID)
	}

	_, err := prepareCall(EndpointCreateEvent, map[string]string{}, input)

	if err != nil {
		if apiErr, apiErrOK := err.(APIError); apiErrOK {
			apiErr.SystemCode = CodeEventCouldNotBeCreated
			return apiErr
		}
	}

	return err
}

// DeleteEventTypeByID deletes an event type by the id
//
// API Doc: https://developers.hubspot.com/docs/methods/timeline/delete-event-type
func DeleteEventTypeByID(eventTypeID int64) error {
	_, err := prepareCall(EndpointDeleteEventType, map[string]string{
		":eventTypeID": fmt.Sprintf("%d", eventTypeID),
	}, nil)

	if err != nil {
		if apiErr, apiErrOK := err.(APIError); apiErrOK {
			apiErr.SystemCode = CodeEventTypeCouldNotBeDeleted
			return apiErr
		}
	}

	return err
}
