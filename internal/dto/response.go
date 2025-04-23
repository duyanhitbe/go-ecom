package dto

import (
	"encoding/json"
	"github.com/duyanhitbe/go-ecom/pkg/utils"
	"net/http"
)

type Meta struct {
	Page      *int32 `json:"page"`
	PerPage   *int32 `json:"per_page"`
	Total     *int32 `json:"total_items"`
	TotalPage *int32 `json:"total_pages"`
	NextPage  *int32 `json:"next_page"`
	PrevPage  *int32 `json:"prev_page"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

type Response[T any] struct {
	StatusCode *int     `json:"status_code,omitempty"`
	Message    *string  `json:"message,omitempty"`
	Success    *bool    `json:"success,omitempty"`
	Data       *T       `json:"data,omitempty"`
	Meta       *Meta    `json:"meta,omitempty"`
	Errors     *[]Error `json:"errors,omitempty"`
	Error      *Error   `json:"error,omitempty"`
}

func NewMeta(perPage, page, total int32) *Meta {
	totalPages := utils.CalculateTotalPage(&perPage, &total)
	nextPage := utils.CalculateNextPage(&page, totalPages)
	prevPage := utils.CalculatePrevPage(&page)

	return &Meta{
		Page:      &page,
		PerPage:   &perPage,
		Total:     &total,
		TotalPage: totalPages,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}
}

func NewResponse[T any](statusCode int, success bool, data *T) *Response[T] {
	message := http.StatusText(statusCode)
	return &Response[T]{
		StatusCode: &statusCode,
		Message:    &message,
		Success:    &success,
		Data:       data,
	}
}

func NewOKResponse[T any](data *T) *Response[T] {
	statusCode := http.StatusOK
	return NewResponse(statusCode, true, data)
}

func NewCreatedResponse[T any](data *T) *Response[T] {
	statusCode := http.StatusCreated
	return NewResponse(statusCode, true, data)
}

func NewPaginatedResponse[T any](data T, meta *Meta) *Response[T] {
	statusCode := http.StatusOK
	rsp := NewResponse(statusCode, true, &data)
	rsp.Meta = meta
	return rsp
}

func NewErrResponse(statusCode int, err error) *Response[any] {
	var data *any = nil
	rsp := NewResponse(statusCode, true, data)
	msg := err.Error()
	rsp.Error = &Error{
		Message: msg,
	}
	return rsp
}

func NewErrorResponse(statusCode int, err Error) *Response[any] {
	var data *any = nil
	rsp := NewResponse(statusCode, true, data)
	rsp.Error = &err
	return rsp
}

func NewErrorsResponse(statusCode int, errs []Error) *Response[any] {
	var data *any = nil
	rsp := NewResponse(statusCode, true, data)
	rsp.Errors = &errs
	return rsp
}

func Write[T any](w http.ResponseWriter, rsp *Response[T]) {
	statusCode := rsp.StatusCode
	b, _ := json.Marshal(rsp)
	w.WriteHeader(*statusCode)
	w.Write(b)
}
