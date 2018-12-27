package models

type Channel struct {
	Channelname string `xorm:"VARCHAR(255)"`
	Channel     string `xorm:"not null VARCHAR(255)"`
	Id          int    `xorm:"not null pk INT(11)"`
	Updatetime  int    `xorm:"default 1323855882 INT(11)"`
	Platform    string `xorm:"VARCHAR(255)"`
}
