package models

import "database/sql"

type Core struct {
	Service         Service
	StoreService    StoreService
	DB              *sql.DB
	RequestResponse RequestResponse
}

func (c *Core) ConnectResponseWriter() {
	rw := &BasicRequestResponse{}
	c.RequestResponse = rw
}
