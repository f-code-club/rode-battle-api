package pagination

type Params struct {
	Page     int `query:"page" default:"1" minimum:"1" maximum:"100000"`
	PageSize int `query:"page_size" default:"20" minimum:"1" maximum:"100"`
}

func (p Params) Limit() int32 {
	return int32(p.PageSize)
}

func (p Params) Offset() int32 {
	return int32((p.Page - 1) * p.PageSize)
}
