package app

import "time"

type TickMsg struct{ T time.Time }
type ErrorMsg struct{ Err error }
