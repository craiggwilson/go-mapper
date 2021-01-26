package naming

// adapted version of strings.FieldsFunc.
func splitFunc(s string, f func(rune) bool) []string {
	type span struct {
		start int
		end   int
	}
	spans := make([]span, 0, 4)

	start := 0
	for end, r := range s {
		if f(r) && start < end {
			spans = append(spans, span{start, end})
			start = end
		}
	}

	if start < len(s) {
		spans = append(spans, span{start, len(s)})
	}

	a := make([]string, len(spans))
	for i, span := range spans {
		a[i] = s[span.start:span.end]
	}
	return a
}

