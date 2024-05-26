// Advanced JSON Customization

package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

// advanced json customization
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}

// concept of custom json decoding: more on this later
