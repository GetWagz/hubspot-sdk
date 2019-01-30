package hubspot

// System codes are used to convey specific statuses in return objects
const (
	CodeContactCouldNotBeCreated = "the contact could not be created at hubspot; you should check the Message field"
	CodeContactNoEmail           = "you failed to specify an email for that contact"
	CodeContactNotFound          = "contact could not be found"
	CodeContactVIDZero           = "the contact VID cannot be 0 for this action"

	CodeEventTypeCouldNotBeCreated = "the event type could not be created"
	CodeEventTypeMissingData       = "the input is missing required information"
	CodeEventTypeCouldNotBeDeleted = "that event type could not be deleted"

	CodeEventCouldNotBeCreated = "the event could not be created"
	CodeEventMissingData       = "the input is missing required information"

	CodeGeneralError = "a general error occurred"
)
