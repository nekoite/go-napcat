package message

type Sex string

type GroupRole string

const (
	SexMale    = "male"
	SexFemale  = "female"
	SexUnknown = "unknown"

	GroupRoleOwner  = "owner"
	GroupRoleAdmin  = "admin"
	GroupRoleMember = "member"
)

type ISender interface {
	GetUserId() int64
	GetNickname() string
	GetSex() Sex
	GetAge() int32
}

type Sender struct {
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      Sex    `json:"sex"`
	Age      int32  `json:"age"`
}

type GroupSender struct {
	Sender
	Card  string    `json:"card"`
	Area  string    `json:"area"`
	Level string    `json:"level"`
	Role  GroupRole `json:"role"`
	Title string    `json:"title"`
}

func (s *Sender) GetUserId() int64 {
	return s.UserId
}

func (s *Sender) GetNickname() string {
	return s.Nickname
}

func (s *Sender) GetSex() Sex {
	return s.Sex
}

func (s *Sender) GetAge() int32 {
	return s.Age
}
