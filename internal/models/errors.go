package models

import "errors"

var (
	ErrNotFound  = errors.New("metric not found")
	ErrWrongHash = errors.New("wrong hash")
	ErrNotDB     = errors.New("not connect to DB")
)
