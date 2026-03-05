



package uuid

import (
	"database/sql/driver"
	"fmt"
)




func (uuid *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		return nil

	case string:
		
		if src == "" {
			return nil
		}

		
		u, err := Parse(src)
		if err != nil {
			return fmt.Errorf("Scan: %v", err)
		}

		*uuid = u

	case []byte:
		
		if len(src) == 0 {
			return nil
		}

		
		
		if len(src) != 16 {
			return uuid.Scan(string(src))
		}
		copy((*uuid)[:], src)

	default:
		return fmt.Errorf("Scan: unable to scan type %T into UUID", src)
	}

	return nil
}




func (uuid UUID) Value() (driver.Value, error) {
	return uuid.String(), nil
}
