package model

type Subscription struct {
	Id         int64
	Id_user    int64  `json:"id_user"`
	Id_country string `json:"id_country"`
}

type Country struct {
	Id   int64
	Name string `json:"location_name"`
}
