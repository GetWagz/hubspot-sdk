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

// Contact is an individual person that you would like to create or update in Hubspot. Email is the only required field.
// By default, the VID field will be 0. After a create call, the VID will be filled in for you.
// Additional porperties are tricky
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
	props := contact.convertContactToProperties()
	send := map[string]interface{}{
		"properties": props,
	}

	ret, err := prepareCall(EndpointCreateContact, map[string]string{
		":email": contact.Email,
	}, send)

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

	_, err := prepareCall(EndpointDeleteContact, map[string]string{
		":vid": fmt.Sprintf("%d", vid),
	}, nil)
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

// GetContactByEmail gets a single contact by their email address
//
// API Doc: https://developers.hubspot.com/docs/methods/contacts/get_contact_by_email
func GetContactByEmail(email string) (Contact, error) {
	contact := Contact{}
	if email == "" || !strings.Contains(email, "@") {
		return contact, APIError{
			HTTPCode:   http.StatusBadRequest,
			SystemCode: CodeContactNoEmail,
			Message:    "you must specify an email address for that contact",
			Body:       nil,
		}
	}
	res, err := prepareCall(EndpointGetContact, map[string]string{
		":email": email,
	}, nil)
	if err != nil {
		if apiErr, apiErrOK := err.(APIError); apiErrOK {
			if apiErr.HTTPCode == 404 {
				apiErr.SystemCode = CodeContactNotFound
				return contact, apiErr
			}
			apiErr.SystemCode = CodeGeneralError
			return contact, apiErr
		}
	}
	// this part is fun; we don't want to worry about all of the properties,
	// so let's loop and figure it out
	fields := res.Body.(map[string]interface{})
	contact.populateContactFields(fields)
	return contact, err
}

func (contact *Contact) populateContactFields(fields map[string]interface{}) {
	for k, v := range fields {
		if k == "vid" {
			vidF := v.(float64)
			contact.VID = int64(vidF)
		}
		if k == "properties" {
			// this is the really big one
			// we need to loop over all of the possible proprties; this is a map of string/interfaces
			properties := v.(map[string]interface{})
			for pK, pV := range properties {
				vals := pV.(map[string]interface{})
				switch pK {
				case "firstname":
					contact.FirstName = vals["value"].(string)
				case "lastname":
					contact.LastName = vals["value"].(string)
				case "email":
					contact.Email = vals["value"].(string)
				case "address":
					contact.Address = vals["value"].(string)
				case "city":
					contact.City = vals["value"].(string)
				case "state":
					contact.State = vals["value"].(string)
				case "zip":
					contact.Zip = vals["value"].(string)
				case "company":
					contact.Company = vals["value"].(string)
				case "phone":
					contact.Phone = vals["value"].(string)
				case "website":
					contact.Website = vals["value"].(string)
				}
			}
		}
	}
}

func (contact *Contact) convertContactToProperties() []map[string]string {
	props := []map[string]string{}
	contactParsed := structs.Map(contact)
	for k, v := range contactParsed {
		// get the string value, if possible
		if val, valOK := v.(string); valOK && v != "" {
			props = append(props, map[string]string{
				"property": strings.ToLower(k),
				"value":    val,
			})
		}
	}
	// now we need to merge the additional properties
	if contact.AdditionalProperties != nil {
		for _, p := range *contact.AdditionalProperties {
			props = append(props, map[string]string{
				"property": p.Property,
				"value":    p.Value,
			})
		}
	}
	return props
}
