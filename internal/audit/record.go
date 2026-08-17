package audit

type AuditRecord struct {
	RowNumber int
	Values    map[string]string
}
