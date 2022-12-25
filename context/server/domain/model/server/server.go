package server

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/block"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// ブロックの上限数です
	BlockMaxAmount = 10
)

// Discordのサーバーです
type Server struct {
	id          ID
	adminRoleID RoleID
	block       []block.Block
}

// Discordのサーバーを作成します
func NewServer(id ID) (*Server, error) {
	s := &Server{}
	s.id = id
	s.adminRoleID = RoleID{}
	s.block = []block.Block{}

	if err := s.validate(); err != nil {
		return s, errors.NewError("検証に失敗しました", err)
	}

	return s, nil
}

// -------------------
// setter
// -------------------

// 管理者のロールIDを更新します
func (u *Server) UpdateAdminRoleID(admin RoleID) error {
	u.adminRoleID = admin

	if err := u.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}

// ブロックを更新します
func (u *Server) UpdateBlock(b []block.Block) error {
	u.block = b

	if err := u.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}

// -------------------
// getter
// -------------------

// IDを取得します
func (u *Server) ID() ID {
	return u.id
}

// 管理者のロールIDを取得します
func (u *Server) AdminRoleID() RoleID {
	return u.adminRoleID
}

// ブロックを取得します
func (u *Server) Block() []block.Block {
	return u.block
}

// -------------------
// validation
// -------------------

// 検証します
func (u *Server) validate() error {
	if len(u.block) > BlockMaxAmount {
		return errors.NewError("ブロックの数が上限を超えています")
	}

	return nil
}

// -------------------
// JSON
// -------------------

// 構造体をJSONに変換します
func (u *Server) MarshalJSON() ([]byte, error) {
	j := struct {
		ID          ID            `json:"id"`
		AdminRoleID RoleID        `json:"admin_role_id"`
		Block       []block.Block `json:"block"`
	}{
		ID:          u.id,
		AdminRoleID: u.adminRoleID,
		Block:       u.block,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません")
	}

	return b, nil
}

// JSONを構造体を変換します
func (u *Server) UnmarshalJSON(b []byte) error {
	j := &struct {
		ID          ID            `json:"id"`
		AdminRoleID RoleID        `json:"admin_role_id"`
		Block       []block.Block `json:"block"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	u.id = j.ID
	u.adminRoleID = j.AdminRoleID
	u.block = j.Block

	return nil
}
