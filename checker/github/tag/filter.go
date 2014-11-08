package tag

type filterFunc func(*Tag) bool

func MakeChanFilter(filter filterFunc) func(Chan) Chan {
	return func(in Chan) Chan {
		return filterChan(in, filter)
	}
}

func filterChan(in Chan, isOK filterFunc) Chan {
	out := make(chanRW)
	go func() {
		for t := range in {
			if isOK(t) {
				out <- t
			}
		}
		close(out)
	}()
	return onlyReadable(out)
}
