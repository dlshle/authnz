package contract

/*
 * A contract connects a subject with multiple groups
 */

type Contract struct {
	ID        string `db:"id"`
	SubjectID string `db:"subject_id"`
	GroupID   string `db:"group_id"`
}
