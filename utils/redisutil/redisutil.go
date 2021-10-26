package redisutil

func IsNoRecord(err error)  bool {
	if err.Error() == "redis: nil"{
		return true
	}
	return false
}