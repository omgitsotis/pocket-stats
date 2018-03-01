package dao

import "github.com/omgitsotis/pocket-stats/pocket/model"

type DAO interface {
	AddUser(string) (int64, error)
	AddArticle(model.Row) error
	IsUser(int64) (bool, error)
	GetStatsForDate(int64) error
}
