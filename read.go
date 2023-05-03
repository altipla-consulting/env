package env

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// MustRead asserts an env variable exists and it is not empty. It reads its value.
// If the value starts with "base64://" it will be decoded as well before returning.
func MustRead(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic(fmt.Sprintf("missing %s environment variable", name))
	}

	if strings.HasPrefix(v, "base64://") {
		decoded, err := base64.StdEncoding.DecodeString(v[9:])
		if err != nil {
			panic(fmt.Sprintf("invalid base64 %s environment variable: %s", name, err.Error()))
		}
		v = string(decoded)
	}

	return v
}
