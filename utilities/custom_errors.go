package utilities

import "errors"

/*
 All custom inventory management errors should be defined here.
*/

var LimitError = errors.New("Limit value cannot be <= 0")
var OffsetError = errors.New("Offset value cannot be < 0")
