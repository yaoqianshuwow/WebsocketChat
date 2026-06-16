package model

import "time"

// Contact 好友关系
type Contact struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserId      int64     `gorm:"column:user_id;index:idx_user_contact,unique;not null"`
	ContactId   int64     `gorm:"column:contact_id;index:idx_user_contact,unique;not null"`
	ContactType int32     `gorm:"column:contact_type;default:1"` // 1=好友
	Nickname    string    `gorm:"column:nickname;size:64"`
	Avatar      string    `gorm:"column:avatar;size:512"`
	Status      int32     `gorm:"column:status;default:0"` // 0=正常 1=黑名单
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Contact) TableName() string {
	return "contact"
}

// ContactApply 好友申请
type ContactApply struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	FromId    int64     `gorm:"column:from_id;index;not null"`
	ToId      int64     `gorm:"column:to_id;index;not null"`
	ApplyType int32     `gorm:"column:apply_type;default:1"` // 1=好友申请
	Remark    string    `gorm:"column:remark;size:255"`
	Status    int32     `gorm:"column:status;default:0"` // 0=待处理 1=同意 2=拒绝
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (ContactApply) TableName() string {
	return "contact_apply"
}

// Session 会话
type Session struct {
	Id             int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserId         int64     `gorm:"column:user_id;uniqueIndex:idx_session_user_peer_type;not null"`
	PeerId         int64     `gorm:"column:peer_id;uniqueIndex:idx_session_user_peer_type;not null"`
	SessionType    int32     `gorm:"column:session_type;uniqueIndex:idx_session_user_peer_type;default:1"` // 1=单聊 2=群聊
	SessionName    string    `gorm:"column:session_name;size:128"`
	Avatar         string    `gorm:"column:avatar;size:512"`
	LastMsgId      int64     `gorm:"column:last_msg_id;default:0"`
	LastMsgContent string    `gorm:"column:last_msg_content;size:512"`
	LastMsgTime    int64     `gorm:"column:last_msg_time;default:0"`
	UnreadCount    int32     `gorm:"column:unread_count;default:0"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Session) TableName() string {
	return "session"
}

// GroupInfo 群组信息
type GroupInfo struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;size:128;not null"`
	Avatar      string    `gorm:"column:avatar;size:255"`
	OwnerId     int64     `gorm:"column:owner_id;index;not null"`
	MemberCount int32     `gorm:"column:member_count;default:0"`
	AddMode     int32     `gorm:"column:add_mode;default:1"` // 1=需要验证 2=直接加入
	Status      int32     `gorm:"column:status;default:0"`   // 0=正常 1=已解散
	Notice      string    `gorm:"column:notice;size:512"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (GroupInfo) TableName() string {
	return "group_info"
}

// GroupMember 群组成员
type GroupMember struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	GroupId     int64     `gorm:"column:group_id;index:idx_group_user,unique;not null"`
	UserId      int64     `gorm:"column:user_id;index:idx_group_user,unique;not null"`
	Role        int32     `gorm:"column:role;default:0"` // 0=普通成员 1=管理员 2=群主
	Nickname    string    `gorm:"column:nickname;size:64"`
	Avatar      string    `gorm:"column:avatar;size:512"`
	UnreadCount int32     `gorm:"column:unread_count;default:0"`
	JoinedAt    time.Time `gorm:"column:joined_at"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (GroupMember) TableName() string {
	return "group_member"
}
