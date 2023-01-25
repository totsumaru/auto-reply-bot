package server

import (
	"encoding/json"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/comment/block"
	"github.com/techstart35/auto-reply-bot/context/server/domain/model/server/rule"
	"github.com/techstart35/auto-reply-bot/context/shared/errors"
)

const (
	// ブロックの上限数です
	BlockMaxAmount = 30
)

// Discordのサーバーです
type Server struct {
	id          model.ID
	adminRoleID model.RoleID
	block       []block.Block
	comment     comment.Comment
	rule        rule.Rule
}

// Discordのサーバーを作成します
func NewServer(id model.ID) (*Server, error) {
	s := &Server{}
	s.id = id
	s.adminRoleID = model.RoleID{}
	s.block = []block.Block{}
	s.rule = rule.Rule{}

	if err := s.validate(); err != nil {
		return s, errors.NewError("検証に失敗しました", err)
	}

	return s, nil
}

// -------------------
// setter
// -------------------

// 管理者のロールIDを更新します
func (u *Server) UpdateAdminRoleID(admin model.RoleID) error {
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

// ルールを更新します
func (u *Server) UpdateRule(r rule.Rule) error {
	u.rule = r

	if err := u.validate(); err != nil {
		return errors.NewError("検証に失敗しました", err)
	}

	return nil
}

// -------------------
// getter
// -------------------

// IDを取得します
func (u *Server) ID() model.ID {
	return u.id
}

// 管理者のロールIDを取得します
func (u *Server) AdminRoleID() model.RoleID {
	return u.adminRoleID
}

// ブロックを取得します
func (u *Server) Block() []block.Block {
	return u.block
}

// ルールを取得します
func (u *Server) Rule() rule.Rule {
	return u.rule
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
	// TODO: 修正が完了したら削除
	c, err := comment.NewComment(u.block)
	if err != nil {
		return nil, errors.NewError("コメントを作成できません", err)
	}

	j := struct {
		ID          model.ID        `json:"id"`
		AdminRoleID model.RoleID    `json:"admin_role_id"`
		Block       []block.Block   `json:"block"`
		Comment     comment.Comment `json:"comment"`
		Rule        rule.Rule       `json:"rule"`
	}{
		ID:          u.id,
		AdminRoleID: u.adminRoleID,
		Block:       u.block,
		Comment:     c,
		Rule:        u.rule,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, errors.NewError("構造体をJSONに変換できません", err)
	}

	return b, nil
}

// JSONを構造体を変換します
func (u *Server) UnmarshalJSON(b []byte) error {
	j := &struct {
		ID          model.ID        `json:"id"`
		AdminRoleID model.RoleID    `json:"admin_role_id"`
		Block       []block.Block   `json:"block"`
		Comment     comment.Comment `json:"comment"`
		Rule        rule.Rule       `json:"rule"`
	}{}

	if err := json.Unmarshal(b, &j); err != nil {
		return errors.NewError("JSONを構造体に変換できません")
	}

	u.id = j.ID
	u.adminRoleID = j.AdminRoleID
	u.block = j.Block
	u.rule = j.Rule

	// TODO: 修正が完了したら削除
	c, err := comment.NewComment(u.block)
	if err != nil {
		return errors.NewError("コメントを作成できません", err)
	}
	u.comment = c

	return nil
}
