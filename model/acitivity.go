package model

type Activity struct {
	Title string `json:"title"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

func (Activity) TebleName() string {
	return "activity"
}

//分页查询
func GetActivityByPaging(number, page int) ([]Activity, int, int, error) {
	var list []Activity
	var total int

	return list, total, page, nil
}
