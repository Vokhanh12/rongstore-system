package errors

var catalog map[string]AppError

func InitCatalog(m map[string]AppError) {
	catalog = m
}

func Lookup(code string) *AppError {
	if e, ok := catalog[code]; ok {
		return copy(e)
	}
	return copy(UNKNOWN_DOMAIN_KE)
}
