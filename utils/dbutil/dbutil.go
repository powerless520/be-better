package dbutil

func IsNoRecord(err error)  bool {
	if err.Error() == "record not found"{
		return true
	}
	return false
}