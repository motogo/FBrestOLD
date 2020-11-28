package config

import ("database/sql"
_"github.com/nakagami/firebirdsql"
)

func Conn() (db *sql.DB, err error) {

	db, err = sql.Open("firebirdsql", "SYSDBA:masterkey@localhost:3050/D:/Data/DokuMents/DOKUMENTS30.FDB")
	return
}