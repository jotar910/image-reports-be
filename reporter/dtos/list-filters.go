package dtos

type ListFilters struct {
	Page  int64 `form:"page"`
	Count int64 `form:"count"`
}

func (f *ListFilters) CheckDefaults() {
	if f.Page == 0 {
		f.Page = 1
	}
	if f.Count == 0 {
		f.Count = 50
	}
}
