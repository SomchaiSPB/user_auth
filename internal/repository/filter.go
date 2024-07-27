package repository

type Filter interface {
	Query() interface{}
	Args() interface{}
}
