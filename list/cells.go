package list

type Row map[string]string

const RowIDKey = "_id"

func (r Row) ID() string {
	return r[RowIDKey]
}

func (r *Row) SetID(id string) {
	(*r)[RowIDKey] = id
}
