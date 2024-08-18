package qq

type Sex string

type GroupRole string

const (
	SexMale    Sex = "male"
	SexFemale  Sex = "female"
	SexUnknown Sex = "unknown"

	GroupRoleOwner  GroupRole = "owner"
	GroupRoleAdmin  GroupRole = "admin"
	GroupRoleMember GroupRole = "member"
)

type IUser interface {
	GetUserId() int64
	GetNickname() string
}

type BasicUser struct {
	UserId   int64  `json:"user_id"`
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

func (s BasicUser) GetUserId() int64 {
	return s.UserId
}

func (s BasicUser) GetNickname() string {
	return s.Nickname
}
