package utils

var isTesting bool
var testingHasBeenSetted bool

// SetTesting set IsTesting value, could be called only one time
func SetTesting(t bool) {
	if !testingHasBeenSetted {
		isTesting = t
		testingHasBeenSetted = true
	}
}

// IsTesting returns true if is testing
func IsTesting() bool {
	return isTesting
}
