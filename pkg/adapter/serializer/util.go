package serializer

import (
	"fmt"
	"strings"
)

func SplitDoguName(doguName string) (string, string, error) {
	splitName := strings.Split(doguName, "/")
	if len(splitName) != 2 {
		return "", "", fmt.Errorf("dogu name needs to be in the form 'namespace/dogu' but is '%s'", doguName)
	}
	return splitName[0], splitName[1], nil
}
