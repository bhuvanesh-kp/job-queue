package api

import "time"

type UserRequest struct {
	ReqType  string
	Duration time.Duration
}
