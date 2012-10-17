package lwes

import "testing"

// just so it doesn't get changed without fully understanding what you are
// doing.
func TestMAX_MSG_SIZE(t *testing.T) {
    if MAX_MSG_SIZE != 65507 {
        t.Fatalf("MAX_MSG_SIZE expected to be 65507")
    }
}

func TestMAX_SHORT_STRING_SIZE(t *testing.T) {
    if MAX_SHORT_STRING_SIZE != 255 {
        t.Fatalf("MAX_SHORT_STRING_SIZE expected to be 255")
    }
}
