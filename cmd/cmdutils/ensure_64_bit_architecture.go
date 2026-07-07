package cmdutils

import (
	"strconv"

	"github.com/thenativeweb/esdm/logging"
)

func Ensure64BitArchitecture() {
	// Check the size of an int to determine whether the
	// application is running on a 64-bit architecture
	// or not.
	if strconv.IntSize != 64 {
		logging.Fatal("64-bit architecture required")
	}
}
