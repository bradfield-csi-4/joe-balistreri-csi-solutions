package db

import "errors"

var ErrNotFound = errors.New("NotFound")
var ErrKeyDeleted = errors.New("KeyDeleted")
