package errorhandler

type DatabaseError struct {
	Message string
}

func (databaseError DatabaseError) Error() string {
	return databaseError.Message
}

type DataNotFoundError struct {
	Message string
}

func (dataNotFoundError DataNotFoundError) Error() string {
	return dataNotFoundError.Message
}
