package list

type Row map[string]string

const RowIDKey = "_id"

func (r Row) ID() string {
	return r[RowIDKey]
}

func (r *Row) SetID(id string) {
	(*r)[RowIDKey] = id
}

func (r *Row) Keys() []string {
	keys := make([]string, 0)
	for k := range *r {
		keys = append(keys, k)
	}
	return keys
}
