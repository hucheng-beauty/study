package model

import "encoding/json"

type Profile struct {
	Name       string // 昵称
	Gender     string // 性别
	Age        int    // 年龄
	Height     int    // 身高
	Weight     int    // 体重
	Income     string // 收入
	Marriage   string // 婚姻状况
	Education  string // 学历
	Occupation string // 职业
	Hokou      string // 户口
	Xinzuo     string // 星座
	House      string // 房
	Car        string // 车
}

// FromJsonObj transform interface{} to Type Profile
func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}

	err = json.Unmarshal(s, &profile)
	return profile, err
}
