package model

type Rank struct {
	UserID       int     `json:"user_id"`
	GoldCoin     float64 `json:"gold_Coin"`
	Nickname     string  `json:"nickname"`
	FaceID       int     `json:"face_id"`
	HeadImageUrl string  `json:"head_image_url"`
	LevelNum     int     `json:"level_num"`
	Rank         int     `json:"rank,omitempty"`
	RoleID       int     `json:"role_id"`
	SuitID       int     `json:"suit_id"`
	PhotoFrameID int     `json:"photo_frame_id"`
}

func (Rank) TableName() string {
	return "rank"
}
