package mysql

import (
	"database/sql"
	"fmt"
	"log"
)

// MySQL用のインフラです
// データベースハンドラやステートメントを管理する
type Infra struct {
	Tx *sql.Tx
}

// ステートメントを閉じます
func (i *Infra) CloseStmt(s *sql.Stmt) {
	if err := s.Close(); err != nil {
		log.Println(fmt.Errorf("SQLステートメントを閉じれません: %w", err).Error())
	}
}
