package aiclient

import "errors"

var (
	ErrAiResponseBadFormat   = errors.New("response from AI model is not in JSON format")
	ErrAiResponseEmpty       = errors.New("response from AI model is empty")
	ErrAiResponseNotApproved = errors.New("self-critique loop ends without approval")
)
