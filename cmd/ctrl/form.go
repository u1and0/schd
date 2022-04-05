package ctrl

// Form : Post, Update, Deleteで使うForm情報
type Form struct {
	ID0   string `form:"id0"`
	ID1   string `form:"id1"`
	Konpo `json:"梱包"`
	Syuka `json:"出荷"`
	Noki  `json:"納期"`
}
