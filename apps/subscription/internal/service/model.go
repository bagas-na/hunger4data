package service

type Subscription struct {
	Id         int64
	Id_user    int64
	Id_country string
}

type Country struct {
	Id   int64
	Name string `json:"location_name"`
}
