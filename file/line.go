package file

func def(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

type line struct {
	record []string
	l      int
}

func (l *line) reload(record []string) {
	l.record = record
	l.l = len(record)
}

func (l line) Pos(idx int) string {
	if idx >= l.l {
		return ""
	}
	return l.record[idx]
}
