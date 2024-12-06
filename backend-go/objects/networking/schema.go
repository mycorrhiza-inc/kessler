package networking

type BasePaginationNetworkSchema struct {
	MaxHits     uint `json:"max_hits"`
	StartOffset uint `json:"start_offset"`
}
