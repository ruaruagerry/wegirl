package rconst

// HomeTags 导航栏
type HomeTags struct {
	CID   string `json:"cid"`
	Title string `json:"title"`
}

const (
	// HashHomeTagsConfig 导航栏配置
	HashHomeTagsConfig = "wegirl:home:tagsconfig"
)
