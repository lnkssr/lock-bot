package models

type Headers struct {
	ContentType   string
	Accept        string
	Authorization string
}

func (h Headers) ToMap() map[string]string {
	m := map[string]string{}
	if h.ContentType != "" {
		m["Content-Type"] = h.ContentType
	}
	if h.Accept != "" {
		m["Accept"] = h.Accept
	}
	if h.Authorization != "" {
		m["Authorization"] = h.Authorization
	}
	return m
}
