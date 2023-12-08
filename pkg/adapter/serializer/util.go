package serializer

import (
	"fmt"
	"strings"
)

// SplitDoguName splits a qualified dogu name into the namespace and the simple name or raises an error if this is not possible.
// "official/nginx" -> "official", "nginx"
func SplitDoguName(doguName string) (string, string, error) {
	splitName := strings.Split(doguName, "/")
	if len(splitName) != 2 {
		return "", "", fmt.Errorf("dogu name needs to be in the form 'namespace/dogu' but is '%s'", doguName)
	}
	return splitName[0], splitName[1], nil
}
