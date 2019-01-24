package hubspot

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactCreate(t *testing.T) {
	ConfigSetup()

	// first should fail becaue there is no email address
	input := Contact{
		FirstName: "Test",
		LastName:  "User",
		City:      "Portsmouth",
		AdditionalProperties: &[]ContactProperty{
			ContactProperty{
				Property: "userId",
				Value:    "42",
			},
		},
	}
	err := CreateOrUpdateContact(&input)
	assert.NotNil(t, err)
	apiErr, cOK := err.(APIError)
	require.True(t, cOK)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPCode)
	assert.Equal(t, CodeContactNoEmail, apiErr.SystemCode)

	// this should fail since the userId isn't a valid property
	input = Contact{
		FirstName: "Test",
		LastName:  "User",
		City:      "Portsmouth",
		Email:     "test@test.com",
		AdditionalProperties: &[]ContactProperty{
			ContactProperty{
				Property: "userId",
				Value:    "42",
			},
		},
	}
	err = CreateOrUpdateContact(&input)
	assert.NotNil(t, err)
	apiErr, cOK = err.(APIError)
	require.True(t, cOK)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPCode)
	assert.Equal(t, CodeContactCouldNotBeCreated, apiErr.SystemCode)
	assert.NotNil(t, apiErr.Body)

	// this should successfully complete
	input = Contact{
		FirstName: "Test",
		LastName:  "User",
		City:      "Portsmouth",
		Email:     "test@test.com",
	}
	err = CreateOrUpdateContact(&input)
	assert.Nil(t, err)
	assert.NotZero(t, input.VID)
}

func TestConversion(t *testing.T) {
	input := Contact{
		FirstName: "Kevin",
		LastName:  "Eaton",
		City:      "",
		AdditionalProperties: &[]ContactProperty{
			ContactProperty{
				Property: "userId",
				Value:    "42",
			},
		},
	}
	props := input.convertContactProperties()
	// iterate and make sure that the fields we expect are the fields we get
	firstFound := false
	lastFound := false
	cityFound := false
	userIDFound := false

	for i := range props {
		if props[i]["property"] == "firstname" {
			if props[i]["value"] == input.FirstName {
				firstFound = true
			}
		}
		if props[i]["property"] == "lastname" {
			if props[i]["value"] == input.LastName {
				lastFound = true
			}
		}
		if props[i]["property"] == "city" {
			// this was blank and should be in the slice
			cityFound = true
		}
		if props[i]["property"] == "userId" {
			if props[i]["value"] == "42" {
				userIDFound = true
			}
		}
	}

	assert.True(t, firstFound)
	assert.True(t, lastFound)
	assert.False(t, cityFound)
	assert.True(t, userIDFound)
}
