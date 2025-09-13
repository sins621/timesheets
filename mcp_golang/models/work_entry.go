package models

import "time"

type WorkEntry struct {
    Description string
    Date        time.Time
    Hours       int
    TaskID      int
    CostCodeID  int
}
