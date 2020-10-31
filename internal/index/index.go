package index

type Index interface {
	SortedEntries() []Entry
	Update(sum string, name string)
}
