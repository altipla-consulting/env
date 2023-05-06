package env

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// MustRead asserts an env variable exists and it is not empty. It reads its value.
// If the value starts with "base64://" it will be decoded as well before returning.
func MustRead(name string) string {
	v := OptionalRead(name)
	if v == "" {
		panic(fmt.Sprintf("missing %s environment variable", name))
	}
	return v
}

// MustReadJSON asserts an env variable exists and it is not empty. It reads its value
// and unmarshals it into the given destination. If the value starts with "base64://"
// it will be decoded as well before returning.
func MustReadJSON(name string, dest any) {
	v := MustRead(name)
	if err := json.Unmarshal([]byte(v), dest); err != nil {
		panic(fmt.Sprintf("invalid json %s environment variable: %s", name, err.Error()))
	}
}

// OptionalRead reads an env variable if it exists and it is not empty. It returns
// the value. If the value starts with "base64://" it will be decoded as well before
// returning.
func OptionalRead(name string) string {
	v := os.Getenv(name)
	if strings.HasPrefix(v, "base64://") {
		decoded, err := base64.StdEncoding.DecodeString(v[9:])
		if err != nil {
			panic(fmt.Sprintf("invalid base64 %s environment variable: %s", name, err.Error()))
		}
		v = string(decoded)
	}
	return v
}
