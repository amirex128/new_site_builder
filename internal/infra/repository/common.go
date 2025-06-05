package repository

import (
	"github.com/amirex128/new_site_builder/internal/contract/common"
)

func buildPaginationResponse[T any](items []T, paginationRequestDto common.PaginationRequestDto, count int64) (*common.PaginationResponseDto[T], error) {
	paginationResponse := &common.PaginationResponseDto[T]{
		Items:      items,
		PageNumber: paginationRequestDto.Page,
		TotalPages: int(count / int64(paginationRequestDto.PageSize)),
		TotalCount: count,
	}

	// Calculate total pages
	if count%int64(paginationRequestDto.PageSize) > 0 {
		paginationResponse.TotalPages++
	}

	// Calculate page navigation
	paginationResponse.HasPreviousPage = paginationRequestDto.Page > 1
	paginationResponse.HasNextPage = int64(paginationRequestDto.Page) < int64(paginationResponse.TotalPages)

	return paginationResponse, nil
}
