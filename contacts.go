package hubspot

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/structs"
)

type ContactProperty struct {
	Property string `json:"property"`
	Value    string `json:"value"`
}

type Contact struct {
	VID                  int64
	Email                string
	FirstName            string
	LastName             string
	Website              string
	Company              string
	Phone                string
	Address              string
	City                 string
	State                string
	Zip                  string
	AdditionalProperties *[]ContactProperty
}

// CreateOrUpdateContact creates or updates a contact using an email address and additional properties. If the VID is 0, it
// is assumed that the user has not been created or updated yet.
//
// API Doc: https://developers.hubspot.com/docs/methods/contacts/create_or_update
func CreateOrUpdateContact(contact *Contact) error {
	// if email is blank, return an error
	if contact.Email == "" {
		return APIError{
			HTTPCode:   http.StatusBadRequest,
			SystemCode: CodeContactNoEmail,
			Message:    "you must provide an Email for a contact",
			Body:       nil,
		}
	}

	// generate the props
	props := contact.convertContactProperties()
	send := map[string]interface{}{
		"properties": props,
	}

	ret, err := makeCall("POST", fmt.Sprintf("/contacts/v1/contact/createOrUpdate/email/%s", contact.Email), send)
	if err != nil {
		if apiErr, apiErrOK := err.(APIError); apiErrOK {
			apiErr.SystemCode = CodeContactCouldNotBeCreated
			return apiErr
		}
		return err
	}
	if body, bodyOK := ret.Body.(map[string]interface{}); bodyOK {
		if vidF, vidOK := body["vid"].(float64); vidOK {
			contact.VID = int64(vidF)
		}
	}

	return nil
}

// DeleteContactByVID deletes a single contact by it's VID
//
// API Doc: https://developers.hubspot.com/docs/methods/contacts/delete_contact
func DeleteContactByVID(vid int64) error {
	if vid == 0 {
		return APIError{
			HTTPCode:   http.StatusBadRequest,
			SystemCode: CodeContactVIDZero,
			Message:    "the VID for the contact cannot be 0 when deleting",
			Body:       nil,
		}
	}
	_, err := makeCall("DELETE", fmt.Sprintf("/contacts/v1/contact/vid/%d", vid), nil)
	if err != nil {
		if apiErr, apiErrOK := err.(APIError); apiErrOK {
			if apiErr.HTTPCode == 404 {
				apiErr.SystemCode = CodeContactNotFound
				return apiErr
			}
			apiErr.SystemCode = CodeContactCouldNotBeCreated
			return apiErr
		}
	}
	return err
}

func (input *Contact) convertContactProperties() []map[string]string {
	props := []map[string]string{}
	contact := structs.Map(input)
	for k, v := range contact {
		// get the string value, if possible
		if val, valOK := v.(string); valOK && v != "" {
			props = append(props, map[string]string{
				"property": strings.ToLower(k),
				"value":    val,
			})
		}
	}
	// now we need to merge the additional properties
	if input.AdditionalProperties != nil {
		for _, p := range *input.AdditionalProperties {
			props = append(props, map[string]string{
				"property": p.Property,
				"value":    p.Value,
			})
		}
	}
	return props
}
