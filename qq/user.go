package qq

import "strconv"

type Sex string

type UserId int64

const (
	SexMale    Sex = "male"
	SexFemale  Sex = "female"
	SexUnknown Sex = "unknown"
)

type IUser interface {
	GetUserId() UserId
	GetNickname() string
}

type BasicUser struct {
	UserId   UserId `json:"user_id"`
	Nickname string `json:"nickname"`
}

type BasicUserWithAvatar struct {
	BasicUser
	Avatar string `json:"avatar"`
}

type User struct {
	BasicUser
	Sex Sex   `json:"sex"`
	Age int32 `json:"age"`
}

type BasicFriend struct {
	BasicUser
	Remark string `json:"remark"`
}

type Friend struct {
	BasicFriend
	Sex Sex   `json:"sex"`
	Age int32 `json:"age"`
}

func (s BasicUser) GetUserId() UserId {
	return s.UserId
}

func (s BasicUser) GetNickname() string {
	return s.Nickname
}

func (i UserId) String() string {
	return strconv.FormatInt(int64(i), 10)
}
