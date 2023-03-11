package controller

type Response struct {
	Errno  string `json:"errno" binding:"required"`
	Errmsg string `json:"errmsg" binding:"required"`
}
type User struct {
	Mobile   string `json:"mobile" binding:"required"`
	PassWord string `json:"password" binding:"required"`
	Sms_Code string `json:"sms_code" binding:"required"`
}
type HouseStu struct {
	Acreage   string   `json:"acreage"`
	Address   string   `json:"address"`
	AreaId    string   `json:"area_id"`
	Beds      string   `json:"beds"`
	Capacity  string   `json:"capacity"`
	Deposit   string   `json:"deposit"`
	Facility  []string `json:"facility"`
	MaxDays   string   `json:"max_days"`
	MinDays   string   `json:"min_days"`
	Price     string   `json:"price"`
	RoomCount string   `json:"room_count"`
	Title     string   `json:"title"`
	Unit      string   `json:"unit"`
}
type Resp struct {
	Errno  string `json:"errno" binding:"required"`
	Errmsg string `json:"errmsg" binding:"required"`
	Data   interface {
		//Name string `json:"name,omitempty"`
	} `json:"data,omitempty"`
}
type UserInfoResp struct {
	Errno  string `json:"errno" binding:"required"`
	Errmsg string `json:"errmsg" binding:"required"`
	Data   struct {
		ID         int    `json:"user_id"`    //用户编号
		Name       string `json:"name"`       //用户名
		Mobile     string `json:"mobile"`     //手机号
		Real_name  string `json:"real_name"`  //真实姓名  实名认证
		Id_card    string `json:"id_card"  `  //身份证号  实名认证
		Avatar_url string `json:"avatar_url"` //用户头像路径       通过fastdfs进行图片存储
	} `json:"data"`
}
