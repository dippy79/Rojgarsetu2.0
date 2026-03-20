package services

func clampPage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func clampLimit(limit int) int {
	if limit < 1 || limit > 100 {
		return 20
	}
	return limit
}
