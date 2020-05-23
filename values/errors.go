package values

import "errors"

var (
	ErrIncorrectUUID        = errors.New("Incorrect UUID")
	ErrInvalidUser          = errors.New("Invalid user")
	ErrInvalidDetails       = errors.New("Invalid signin details")
	ErrRetrieveUUID         = errors.New("Could not retrieve UUID")
	ErrMarshal              = errors.New("Could not marshal content")
	ErrWrite                = errors.New("Error while sending content")
	ErrAuthentication       = errors.New("Authentication error")
	ErrIllicitJoinRequest   = errors.New("User was not originally requested to join")
	ErrUserExistInRoom      = errors.New("User already exist in room")
	ErrUserAlreadyRequested = errors.New("User already requested to join room")
)
