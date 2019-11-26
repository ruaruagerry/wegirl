package rconst

// HomeTag 导航栏
type HomeTag struct {
	CID   string `json:"cid"`
	Title string `json:"title"`
}

// HomeImg 图片
type HomeImg struct {
	Title  string `json:"title"`
	Href   string `json:"href"`
	Large  string `json:"large"`
	Thumb  string `json:"thumb"`
	Small  string `json:"small"`
	Filter bool   `json:"filter"`
}

const (
	// HashHomeTagsConfig 导航栏配置
	HashHomeTagsConfig = "wegirl:home:tagsconfig"
)
