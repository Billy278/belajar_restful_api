package helper

import "database/sql"

func CommitOrRollBack(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		PanicIfError(errorRollback)
		panic(err)
	} else {
		errorCommit := tx.Commit()
		PanicIfError(errorCommit)

	}
}
