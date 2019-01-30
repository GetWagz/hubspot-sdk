package hubspot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEventType(t *testing.T) {
	ConfigSetup()

	badInput := EventType{}
	err := CreateNewEventType(&badInput)
	assert.NotNil(t, err)

	input := EventType{
		Name:           "Test Event Type",
		HeaderTemplate: "# Title for event {{id}}\nThis is an event for {{objectType}}",
		DetailTemplate: "This event happened on {{#formatDate timestamp}}{{/formatDate}}",
		ApplicationID:  "123",
		ObjectType:     "CONTACT",
	}

	err = CreateNewEventType(&input)
	require.Nil(t, err)
	// this is mocked in most cirumstances, so just make sure the data is sane
	assert.NotZero(t, input.ID)
	assert.NotEqual(t, "", input.Name)
}

func TestCreateEvent(t *testing.T) {
	ConfigSetup()

	badInput := Event{}
	err := CreateOrUpdateEvent(&badInput)
	assert.NotNil(t, err)

	// we need a valid event type, so we are going to create one
	eventTypeInput := EventType{
		Name:       "Demo Event Type",
		ObjectType: "CONTACT",
	}
	err = CreateNewEventType(&eventTypeInput)
	require.Nil(t, err)

	// create a contact and delete it
	contact := Contact{
		Email: "test@test.com",
	}
	err = CreateOrUpdateContact(&contact)
	require.Nil(t, err)
	require.NotZero(t, contact.VID)

	defer DeleteContactByVID(contact.VID)

	input := Event{
		ObjectID:    contact.VID,
		EventTypeID: eventTypeInput.ID,
		ExtraData:   nil,
		TimelineIFrame: &EventIFrame{
			LinkLabel:   "Click Here",
			IFrameLabel: "My IFrame",
			IFrameURI:   "https://www.hubspot.com",
			Width:       200,
			Height:      500,
		},
	}
	err = CreateOrUpdateEvent(&input)
	assert.Nil(t, err)
	// this is mocked in most cirumstances, so just make sure the data is sane
	assert.NotEqual(t, "", input.ID)

	err = DeleteEventTypeByID(eventTypeInput.ID)
	require.Nil(t, err)
}
