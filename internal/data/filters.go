package data

import "github.com/v3ronez/IDKN/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func ValidateFields(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greated than zero")
	v.Check(f.Page <= 1_000, "page", "must be maximum of one thousand")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be maximun of 100")

	v.Check(validator.PermittdValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}
