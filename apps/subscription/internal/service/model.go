package service

type Subscription struct {
	id         int64
	id_user    int64
	id_country string
}

type Country struct {
	id       int64
	name     string
	category string
}
