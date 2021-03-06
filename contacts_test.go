package hubspot

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

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

func TestContactDelete(t *testing.T) {
	ConfigSetup()

	// 0 is an error
	err := DeleteContactByVID(0)
	assert.NotNil(t, err)
	apiErr, cOK := err.(APIError)
	require.True(t, cOK)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPCode)
	assert.Equal(t, CodeContactVIDZero, apiErr.SystemCode)

	// 1 can't be found
	err = DeleteContactByVID(1)
	assert.NotNil(t, err)
	apiErr, cOK = err.(APIError)
	require.True(t, cOK)
	assert.Equal(t, http.StatusNotFound, apiErr.HTTPCode)
	assert.Equal(t, CodeContactNotFound, apiErr.SystemCode)

	// let's create, use that one as our delete subject
	contact := Contact{
		Email: "test@test.com",
	}
	err = CreateOrUpdateContact(&contact)
	assert.Nil(t, err)
	assert.NotZero(t, contact.VID)

	err = DeleteContactByVID(contact.VID)
	assert.Nil(t, err)
}

func TestContactGet(t *testing.T) {
	ConfigSetup()

	_, err := GetContactByEmail("")
	assert.NotNil(t, err)
	apiErr, cOK := err.(APIError)
	require.True(t, cOK)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPCode)
	assert.Equal(t, CodeContactNoEmail, apiErr.SystemCode)

	// create a random email
	rand.Seed(time.Now().UnixNano())
	random1 := rand.Int63()
	random2 := rand.Int63()
	randomEmail := fmt.Sprintf("%d-rand@%d.notrealdomain", random1, random2)

	_, err = GetContactByEmail(randomEmail)
	assert.NotNil(t, err)
	apiErr, cOK = err.(APIError)
	require.True(t, cOK)
	assert.Equal(t, http.StatusNotFound, apiErr.HTTPCode)
	assert.Equal(t, CodeContactNotFound, apiErr.SystemCode)

	// now, create a contact and then get it and delete it
	goodEmail := fmt.Sprintf("test-%d@test.com", random1)
	contact := Contact{
		Email:     goodEmail,
		FirstName: fmt.Sprintf("First-%d", random1),
		LastName:  fmt.Sprintf("Last-%d", random2),
		Company:   "Wagz",
		Address:   "230 Commerce St",
		City:      "Portsmouth",
		State:     "NH",
		Zip:       "03801",
		Website:   "https://www.wagz.com",
		Phone:     "1-800-GET-WAGZ",
	}
	err = CreateOrUpdateContact(&contact)
	assert.Nil(t, err)
	assert.NotZero(t, contact.VID)

	// make sure we clean up after ourselves
	defer DeleteContactByVID(contact.VID)

	foundContact, err := GetContactByEmail(goodEmail)
	assert.Nil(t, err)
	assert.Equal(t, contact.VID, foundContact.VID)
	assert.Equal(t, contact.Email, foundContact.Email)
	assert.Equal(t, contact.FirstName, foundContact.FirstName)
	assert.Equal(t, contact.LastName, foundContact.LastName)
	assert.Equal(t, contact.Company, foundContact.Company)
	assert.Equal(t, contact.Address, foundContact.Address)
	assert.Equal(t, contact.City, foundContact.City)
	assert.Equal(t, contact.State, foundContact.State)
	assert.Equal(t, contact.Zip, foundContact.Zip)
	assert.Equal(t, contact.Website, foundContact.Website)
	assert.Equal(t, contact.Phone, foundContact.Phone)
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
	props := input.convertContactToProperties()
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
