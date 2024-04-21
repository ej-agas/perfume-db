package postgresql

import "fmt"

var (
	ErrAcquiringConn = fmt.Errorf("error acquiring connection from database connection pool")
	ErrStartingDBTx  = fmt.Errorf("error starting database transaction")
)
