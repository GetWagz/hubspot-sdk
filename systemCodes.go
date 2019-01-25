package hubspot

// System codes are used to convey specific statuses in return objects
const (
	CodeContactCouldNotBeCreated = "the contact could not be created at hubspot; you should check the Message field"
	CodeContactNoEmail           = "you failed to specify an email for that contact"
	CodeContactNotFound          = "contact could not be found"
	CodeContactVIDZero           = "the contact VID cannot be 0 for this action"

	CodeGeneralError = "a general error occurred"
)
