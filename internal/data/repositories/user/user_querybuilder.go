package user

import (
	"fmt"
	"github.com/newpurr/stock/pkg/utils/gormutils"
	"gorm.io/gorm"
	"time"
)

type queryBuilder struct {
	table string
	db    *gorm.DB
	order []string
	group []string
	where []struct {
		prefix string
		value  []interface{}
	}
	limit  int
	offset int
}

// Builder 生成规则
// 1. 公用字段builder，可被其他builder进行组合
// 2. 索引字段才生成OrderBy、GroupBy
// 3. 其他字段生成Eq、WhereIn、WhereNotIn
type userQueryBuilder queryBuilder

func query(db *gorm.DB) *userQueryBuilder {
	return &userQueryBuilder{
		"kratos.user",
		db,
		[]string{},
		[]string{},
		[]struct {
			prefix string
			value  []interface{}
		}{},
		1,
		0,
	}
}

func (b *userQueryBuilder) from(table string) *userQueryBuilder {
	b.table = table
	return b
}

func (b *userQueryBuilder) build() *gorm.DB {
	ret := b.db.Table(b.table)
	for _, where := range b.where {
		ret = ret.Where(where.prefix, where.value...)
	}
	for _, order := range b.order {
		ret = ret.Order(order)
	}
	for _, group := range b.group {
		ret = ret.Group(group)
	}
	ret = ret.Limit(b.limit).Offset(b.offset)
	return ret
}

func (b *userQueryBuilder) resetWhere() *userQueryBuilder {
	b.where = []struct {
		prefix string
		value  []interface{}
	}{}
	return b
}

func (b *userQueryBuilder) addWhere(k string, p gormutils.Predicate, value interface{}) *userQueryBuilder {
	b.where = append(b.where, struct {
		prefix string
		value  []interface{}
	}{
		fmt.Sprintf("%v %v ?", k, p),
		[]interface{}{value},
	})
	return b
}

func (b *userQueryBuilder) addWhereBetween(k string, v1 interface{}, v2 interface{}) *userQueryBuilder {
	b.where = append(b.where, struct {
		prefix string
		value  []interface{}
	}{
		fmt.Sprintf("%v BETWEEN ? AND ?", k),
		[]interface{}{v1, v2},
	})
	return b
}
func (b *userQueryBuilder) whereIdCond(p gormutils.Predicate, value uint32) *userQueryBuilder {
	return b.addWhere("id", p, value)
}

func (b *userQueryBuilder) whereId(value uint32) *userQueryBuilder {
	return b.whereIdCond(gormutils.EqualPredicate, value)
}

func (b *userQueryBuilder) whereIdIn(value []uint32) *userQueryBuilder {
	return b.addWhere("id", gormutils.InPredicate, value)
}

func (b *userQueryBuilder) whereIdNotIn(value []uint32) *userQueryBuilder {
	return b.addWhere("id", gormutils.NotInPredicate, value)
}

func (b *userQueryBuilder) whereIdBetween(v1 uint32, v2 uint32) *userQueryBuilder {
	return b.addWhereBetween("id", v1, v2)
}
func (b *userQueryBuilder) orderById(asc bool) *userQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	b.order = append(b.order, "id "+order)
	return b
}
func (b *userQueryBuilder) groupById(asc bool) *userQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	b.order = append(b.order, "id "+order)
	return b
}
func (b *userQueryBuilder) whereNameCond(p gormutils.Predicate, value string) *userQueryBuilder {
	return b.addWhere("name", p, value)
}

func (b *userQueryBuilder) whereName(value string) *userQueryBuilder {
	return b.whereNameCond(gormutils.EqualPredicate, value)
}

func (b *userQueryBuilder) whereNameIn(value []string) *userQueryBuilder {
	return b.addWhere("name", gormutils.InPredicate, value)
}

func (b *userQueryBuilder) whereNameNotIn(value []string) *userQueryBuilder {
	return b.addWhere("name", gormutils.NotInPredicate, value)
}

func (b *userQueryBuilder) whereNameBetween(v1 string, v2 string) *userQueryBuilder {
	return b.addWhereBetween("name", v1, v2)
}
func (b *userQueryBuilder) whereNameLike(value string) *userQueryBuilder {
	return b.addWhere("name", gormutils.LikePredicate, "%"+value+"%")
}
func (b *userQueryBuilder) wherePwdCond(p gormutils.Predicate, value string) *userQueryBuilder {
	return b.addWhere("pwd", p, value)
}

func (b *userQueryBuilder) wherePwd(value string) *userQueryBuilder {
	return b.wherePwdCond(gormutils.EqualPredicate, value)
}

func (b *userQueryBuilder) wherePwdIn(value []string) *userQueryBuilder {
	return b.addWhere("pwd", gormutils.InPredicate, value)
}

func (b *userQueryBuilder) wherePwdNotIn(value []string) *userQueryBuilder {
	return b.addWhere("pwd", gormutils.NotInPredicate, value)
}

func (b *userQueryBuilder) wherePwdBetween(v1 string, v2 string) *userQueryBuilder {
	return b.addWhereBetween("pwd", v1, v2)
}
func (b *userQueryBuilder) wherePwdLike(value string) *userQueryBuilder {
	return b.addWhere("pwd", gormutils.LikePredicate, "%"+value+"%")
}
func (b *userQueryBuilder) whereIsDeletedCond(p gormutils.Predicate, value int32) *userQueryBuilder {
	return b.addWhere("is_deleted", p, value)
}

func (b *userQueryBuilder) whereIsDeleted(value int32) *userQueryBuilder {
	return b.whereIsDeletedCond(gormutils.EqualPredicate, value)
}

func (b *userQueryBuilder) whereIsDeletedIn(value []int32) *userQueryBuilder {
	return b.addWhere("is_deleted", gormutils.InPredicate, value)
}

func (b *userQueryBuilder) whereIsDeletedNotIn(value []int32) *userQueryBuilder {
	return b.addWhere("is_deleted", gormutils.NotInPredicate, value)
}

func (b *userQueryBuilder) whereIsDeletedBetween(v1 int32, v2 int32) *userQueryBuilder {
	return b.addWhereBetween("is_deleted", v1, v2)
}
func (b *userQueryBuilder) whereCreatedAtCond(p gormutils.Predicate, value time.Time) *userQueryBuilder {
	return b.addWhere("created_at", p, value)
}

func (b *userQueryBuilder) whereCreatedAt(value time.Time) *userQueryBuilder {
	return b.whereCreatedAtCond(gormutils.EqualPredicate, value)
}

func (b *userQueryBuilder) whereCreatedAtIn(value []time.Time) *userQueryBuilder {
	return b.addWhere("created_at", gormutils.InPredicate, value)
}

func (b *userQueryBuilder) whereCreatedAtNotIn(value []time.Time) *userQueryBuilder {
	return b.addWhere("created_at", gormutils.NotInPredicate, value)
}

func (b *userQueryBuilder) whereCreatedAtBetween(v1 time.Time, v2 time.Time) *userQueryBuilder {
	return b.addWhereBetween("created_at", v1, v2)
}
func (b *userQueryBuilder) whereUpdatedAtCond(p gormutils.Predicate, value time.Time) *userQueryBuilder {
	return b.addWhere("updated_at", p, value)
}

func (b *userQueryBuilder) whereUpdatedAt(value time.Time) *userQueryBuilder {
	return b.whereUpdatedAtCond(gormutils.EqualPredicate, value)
}

func (b *userQueryBuilder) whereUpdatedAtIn(value []time.Time) *userQueryBuilder {
	return b.addWhere("updated_at", gormutils.InPredicate, value)
}

func (b *userQueryBuilder) whereUpdatedAtNotIn(value []time.Time) *userQueryBuilder {
	return b.addWhere("updated_at", gormutils.NotInPredicate, value)
}

func (b *userQueryBuilder) whereUpdatedAtBetween(v1 time.Time, v2 time.Time) *userQueryBuilder {
	return b.addWhereBetween("updated_at", v1, v2)
}

/*--------------------------------------------------------------------------
| 添加gorm常用方法，减少使用者心智负担
|-------------------------------------------------------------------------- */

func (b *userQueryBuilder) Update(column string, value interface{}) (tx *gorm.DB) {
	return b.build().Update(column, value)
}

func (b *userQueryBuilder) Save(value interface{}) (tx *gorm.DB) {
	return b.build().Save(value)
}

func (b *userQueryBuilder) Scan(dest interface{}) (tx *gorm.DB) {
	return b.build().Scan(dest)
}

func (b *userQueryBuilder) Pluck(column string, dest interface{}) (tx *gorm.DB) {
	return b.build().Pluck(column, dest)
}

// FindInBatches https://gorm.io/zh_CN/docs/advanced_query.html#FindInBatches
func (b *userQueryBuilder) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB {
	return b.build().FindInBatches(dest, batchSize, fc)
}
