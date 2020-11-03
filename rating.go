package irapi

import (
	"fmt"
	"strconv"
)

// Rating represents a user's ttRating or iRating
type Rating int64

// UnmarshalJSON decodes a Rating from a JSON value
func (r *Rating) UnmarshalJSON(s []byte) error {
	str := string(s)

	if str == "\"---\"" {
		r = nil
		return nil
	}

	v, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		fmt.Printf("Parse Error for '%s', %s\n", str, err)
		return err
	}

	*r = Rating(v)

	return nil
}
