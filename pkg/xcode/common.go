package xcode

// var
//
//	@param Canceled = New(498
//	@param "CANCELED")
//	@param "INTERNAL_ERROR")
//	@param "DEADLINE_EXCEEDED")
//	@author kunarc
//	@update 2024-10-19 06:14:02
var (
	Canceled  = New(498, "CANCELED")
	ServerErr = New(500, "INTERNAL_ERROR")
	Deadline  = New(504, "DEADLINE_EXCEEDED")
)
