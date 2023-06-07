package models

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Models struct {
	db orm.DB
}

type QueryDefaulter interface {
	QueryDefault(*orm.Query) *orm.Query
}

func New(db orm.DB) (*Models, error) {
	m := &Models{
		db: db,
	}

	return m, nil
}

func (m *Models) Get(v any, id string) error {
	md := m.db.Model(v)

	pks := md.TableModel().Table().PKs

	if len(pks) != 1 {
		return fmt.Errorf("one primary key expected")
	}

	for _, pk := range pks {
		md = md.Where(fmt.Sprintf("%s = ?", pk.Column), id)
	}

	if err := md.Select(); err != nil {
		return err
	}

	return nil
}

func (m *Models) List(vs any) error {
	q := m.db.Model(vs)

	if reflect.TypeOf(vs).Kind() != reflect.Ptr || reflect.TypeOf(vs).Elem().Kind() != reflect.Slice {
		return fmt.Errorf("pointer to slice expected")
	}

	v := reflect.New(reflect.TypeOf(vs).Elem()).Interface()

	if qd, ok := v.(QueryDefaulter); ok {
		q = qd.QueryDefault(q)
	}

	return q.Select()
}

func (m *Models) Query(v any) *orm.Query {
	return m.db.Model(v)
}

func (m *Models) Save(v any, columns ...string) error {
	var md *orm.Query

	switch t := v.(type) {
	case *orm.Query:
		md = t
	default:
		md = m.db.Model(t)
	}

	pks := []string{}

	for _, pk := range md.TableModel().Table().PKs {
		pks = append(pks, string(pk.Column))
	}

	md = md.OnConflict(fmt.Sprintf("(%s) DO UPDATE", strings.Join(pks, ",")))

	if ups := m.updateColumns(v); ups != "" {
		md = md.Set(ups)
	}

	for _, column := range columns {
		md = md.Set(fmt.Sprintf("%q = EXCLUDED.%q", column, column))
	}

	if _, err := md.Insert(); err != nil {
		return err
	}

	return nil
}

func (m *Models) TableChanged(ctx context.Context, name string) <-chan string {
	ch := make(chan string)

	go m.listen(ctx, name, func(id string) {
		ch <- id
	})

	return ch
}

func (m *Models) transaction(fn func(*Models) error) error {
	db, ok := m.db.(*pg.DB)
	if !ok {
		return fmt.Errorf("transactions unsupported on model db")
	}

	return db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		mt := *m
		mt.db = tx
		return fn(&mt)
	})
}

func (m *Models) updateColumns(v interface{}, additional ...string) string {
	updates := map[string]bool{}

	for _, a := range additional {
		updates[a] = true
	}

	for field, attrs := range modelTags(v) {
		if attrs["update"] {
			for _, f := range m.db.Model(v).TableModel().Table().Fields {
				if f.GoName == field {
					updates[f.SQLName] = true
				}
			}
		}
	}

	statements := []string{}

	for k := range updates {
		statements = append(statements, fmt.Sprintf(`%q = EXCLUDED.%q`, k, k))
	}

	return strings.Join(statements, ",")
}

func (m *Models) listen(ctx context.Context, channel string, fn func(string)) {
	db, ok := m.db.(*pg.DB)
	if !ok {
		fmt.Println("error: listen unsupported on model db")
	}

	ln := db.Listen(ctx, channel)
	defer ln.Close()

	ch := ln.Channel()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			fn(msg.Payload)
		}
	}
}

func modelTags(v interface{}) map[string]map[string]bool {
	tags := map[string]map[string]bool{}

	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag, ok := f.Tag.Lookup("model"); ok {
			tags[f.Name] = map[string]bool{}
			for _, attr := range strings.Split(tag, ",") {
				tags[f.Name][strings.TrimSpace(attr)] = true
			}
		}
	}

	return tags
}

func sliceType(vs any) (any, error) {
	vst := reflect.TypeOf(vs)

	if vst.Kind() == reflect.Ptr {
		vst = reflect.TypeOf(vs).Elem()
	}

	if vst.Kind() != reflect.Slice {
		return nil, fmt.Errorf("slice expected")
	}

	return reflect.New(vst.Elem()).Interface(), nil
}
