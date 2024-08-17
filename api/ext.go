package api

import (
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/utils"
)

// GetNewResultFunc 用于创建响应对象，应该返回 *Resp[T] 类型
type GetNewResultFunc func() any

type ExtAction struct {
	Action Action
	// GetNewResultFunc 用于创建响应对象，应该返回 *Resp[T] 类型
	GetNewResultFunc GetNewResultFunc
}

type ApiExtension struct {
	Name    string
	Actions map[Action]ExtAction
}

var (
	extensions = make(map[string]ApiExtension)
	extActions = make(map[Action]ExtAction)
)

func NewExtension(name string) *ApiExtension {
	return &ApiExtension{
		Name:    name,
		Actions: make(map[Action]ExtAction),
	}
}

func (ext *ApiExtension) WithAction(action Action, getNewResultFunc GetNewResultFunc) *ApiExtension {
	ext.Actions[action] = ExtAction{
		Action:           action,
		GetNewResultFunc: getNewResultFunc,
	}
	return ext
}

func (ext *ApiExtension) WithActions(actions map[Action]GetNewResultFunc) *ApiExtension {
	for action, getNewResultFunc := range actions {
		ext.Actions[action] = ExtAction{
			Action:           action,
			GetNewResultFunc: getNewResultFunc,
		}
	}
	return ext
}

func (ext *ApiExtension) Register() error {
	return RegisterExtension(*ext)
}

func RegisterExtension(ext ApiExtension) error {
	if _, ok := extensions[ext.Name]; ok {
		return errors.ErrExtensionAlreadyRegistered
	}
	if err := registerActions(ext); err != nil {
		return err
	}
	extensions[ext.Name] = ext
	return nil
}

func registerActions(ext ApiExtension) error {
	alreadyRegisteredActions := utils.NewSet[Action]()
	for action, extAction := range ext.Actions {
		if _, ok := extActions[action]; ok {
			goto ERROR
		}
		alreadyRegisteredActions.Add(action)
		extActions[action] = extAction
	}
	return nil

ERROR:
	for action := range alreadyRegisteredActions {
		delete(extActions, action)
	}
	return errors.ErrActionAlreadyRegistered
}
