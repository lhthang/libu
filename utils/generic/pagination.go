package generic

func Pagination(x int64, skip int64, size int64) (int64, int64) {
	limit := func() int64 {
		if skip+size > x {
			return x
		} else {
			return skip + size
		}

	}

	start := func() int64 {
		if skip > x {
			return x
		} else {
			return skip
		}

	}
	return start(), limit()
}
