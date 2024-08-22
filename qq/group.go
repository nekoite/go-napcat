package qq

type AnonymousData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

type Group struct {
	GroupId        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int32  `json:"member_count"`
	MaxMemberCount int32  `json:"max_member_count"`
}

type GroupUser struct {
	User
	Card  string    `json:"card"`
	Area  string    `json:"area"`
	Level string    `json:"level"`
	Role  GroupRole `json:"role"`
	Title string    `json:"title"`
}

type DetailedGroupUser struct {
	GroupUser
	GroupId         int64 `json:"group_id"`
	JoinTime        int64 `json:"join_time"`
	LastSentTime    int64 `json:"last_sent_time"`
	Unfriendly      bool  `json:"unfriendly"`
	TitleExpireTime int64 `json:"title_expire_time"`
	CardChangeable  bool  `json:"card_changeable"`
}
