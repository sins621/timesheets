package request

import (
	"testing"
	"time"
)

func TestPostTimeSheetEntry(t *testing.T) {
	currentTime := time.Now()
	tsr := NewTimeSheetRequest("https://office.warpdevelopment.com")
	// TODO: Read Token From DB or Something
	err := tsr.PostTimeSheetEntry("token", 4767, 1322, 4, false, 4, currentTime, "Test Commit")
	if err != nil {
		t.Errorf("PostTimeSheetEntry failed with error: %v\n", err)
	}
}
