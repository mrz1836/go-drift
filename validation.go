package drift

// requireString validates that a string is not empty
func requireString(val string, err error) error {
	if len(val) == 0 {
		return err
	}
	return nil
}

// requireID validates that an ID is not zero
func requireID(id uint64, err error) error {
	if id == 0 {
		return err
	}
	return nil
}
