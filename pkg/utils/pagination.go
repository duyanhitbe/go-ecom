package utils

func GetPaginationMeta(pg, perPg *int32) (page int32, perPage int32, offset int32) {
	if pg == nil {
		page = 1
	}
	if perPg == nil {
		perPage = 10
	}

	offset = (page - 1) * perPage

	return
}

func CalculateTotalPage(perPage, total *int32) *int32 {
	if perPage == nil || total == nil {
		return nil
	}
	totalPages := *total / *perPage
	if *total%*perPage != 0 {
		totalPages++
	}
	return &totalPages
}

func CalculateNextPage(page, totalPages *int32) *int32 {
	if page == nil || totalPages == nil {
		return nil
	}
	var nextPage *int32
	if *page < *totalPages {
		next := *page + 1
		nextPage = &next
	}
	return nextPage
}

func CalculatePrevPage(page *int32) *int32 {
	if page == nil {
		return nil
	}
	var prevPage *int32
	if *page > 1 {
		prev := *page - 1
		prevPage = &prev
	}
	return prevPage
}
