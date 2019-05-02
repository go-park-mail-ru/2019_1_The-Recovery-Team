package postgresql

import "github.com/jackc/pgx"

type Conn interface {
	Close()
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
	Query(sql string, args ...interface{}) (Rows, error)
	QueryRow(sql string, args ...interface{}) Row
	Begin() (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
	Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error)
	QueryRow(sql string, args ...interface{}) Row
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) (err error)
}

type Row interface {
	Scan(dest ...interface{}) (err error)
}

func NewConnPool(conn *pgx.ConnPool) Conn {
	return &ConnPool{
		conn: conn,
	}
}

type ConnPool struct {
	conn *pgx.ConnPool
}

func (c *ConnPool) Close() {
	c.conn.Close()
}

func (c *ConnPool) Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error) {
	return c.conn.Exec(sql, arguments...)
}

func (c *ConnPool) Query(sql string, args ...interface{}) (Rows, error) {
	return c.conn.Query(sql, args...)
}

func (c *ConnPool) QueryRow(sql string, args ...interface{}) Row {
	return c.conn.QueryRow(sql, args...)
}

func (c *ConnPool) Begin() (Tx, error) {
	tx, err := c.conn.Begin()
	txPgx := TxPgx{
		Tx: tx,
	}
	return &txPgx, err
}

type TxPgx struct {
	Tx *pgx.Tx
}

func (t *TxPgx) Commit() error {
	return t.Tx.Commit()
}

func (t *TxPgx) Rollback() error {
	return t.Tx.Rollback()
}

func (t *TxPgx) Exec(sql string, arguments ...interface{}) (commandTag pgx.CommandTag, err error) {
	return t.Tx.Exec(sql, arguments...)
}

func (t *TxPgx) QueryRow(sql string, args ...interface{}) Row {
	return t.Tx.QueryRow(sql, args...)
}
