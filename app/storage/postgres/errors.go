package postgres

func checkSQLError(err error) error {
	var typedError error
	switch err {
	case nil:
		return nil
	default:
		return err
	}

	return typedError
}
