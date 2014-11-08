package repository

type Chan <-chan *Repository
type chanW chan<- *Repository
type chanRW chan *Repository

// TODO(leon): This is ugly
func onlyReadable(in chanRW) <-chan *Repository {
	return in
}

// TODO(leon): This is ugly
func onlyWritable(in chanRW) chan<- *Repository {
	return in
}
