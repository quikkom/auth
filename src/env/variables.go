package env

import "os"

var (
	// Secret key which used to generate JWT tokens
	AUTH_SECRET = ""
)

// Set globally used env variables
func Fill() {
	AUTH_SECRET = FindEnv("AUTH_SECRET", "I_AM_SECRET")
}

// Searches an env variable and returns it.
// If it couldn't find it, returns `defaultValue`
func FindEnv(name string, defaultValue string) string {
	env, found := os.LookupEnv(name)

	if !found {
		return defaultValue
	}

	return env
}
