package hubspot

// System codes are used to convey specific statuses in return objects
const (
	CodeContactNoEmail           = "you failed to specify an email for that contact"
	CodeContactCouldNotBeCreated = "the contact could not be created at hubspot; you should check the Message field"

	CodeGeneralError = "a general error occurred"
)
