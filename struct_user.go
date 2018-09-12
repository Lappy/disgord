package disgord

import (
	"encoding/json"
	"errors"
	"sync"
)

const (
	// StatusIdle presence status for idle
	StatusIdle = "idle"
	// StatusDnd presence status for dnd
	StatusDnd = "dnd"
	// StatusOnline presence status for online
	StatusOnline = "online"
	// StatusOffline presence status for offline
	StatusOffline = "offline"
)

const (
	ActivityFlagInstance    = 1 << 0
	ActivityFlagJoin        = 1 << 1
	ActivityFlagSpectate    = 1 << 2
	ActivityFlagJoinRequest = 1 << 3
	ActivityFlagSync        = 1 << 4
	ActivityFlagPlay        = 1 << 5
)

type UserInterface interface {
	Mention() string
	MentionNickname() string
	String() string
}

type ActivityParty struct {
	ID   string `json:"id,omitempty"`   // the id of the party
	Size []int  `json:"size,omitempty"` // used to show the party's current and maximum size
}

// Limit shows the maximum number of guests/people allowed
func (ap *ActivityParty) Limit() int {
	if len(ap.Size) != 2 {
		return 0
	}

	return ap.Size[1]
}

// NumberOfPeople shows the current number of people attending the Party
func (ap *ActivityParty) NumberOfPeople() int {
	if len(ap.Size) != 1 {
		return 0
	}

	return ap.Size[0]
}

type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"` // the id for a large asset of the activity, usually a snowflake
	LargeText  string `json:"large_text,omitempty"`  //text displayed when hovering over the large image of the activity
	SmallImage string `json:"small_image,omitempty"` // the id for a small asset of the activity, usually a snowflake
	SmallText  string `json:"small_text,omitempty"`  //	text displayed when hovering over the small image of the activity
}
type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`     // the secret for joining a party
	Spectate string `json:"spectate,omitempty"` // the secret for spectating a game
	Match    string `json:"match,omitempty"`    // the secret for a specific instanced match
}
type ActivityTimestamp struct {
	Start int `json:"start,omitempty"` // unix time (in milliseconds) of when the activity started
	End   int `json:"end,omitempty"`   // unix time (in milliseconds) of when the activity ends
}

// UserActivity https://discordapp.com/developers/docs/topics/gateway#activity-object-activity-structure
type UserActivity struct {
	Name          string               `json:"name"`                     // the activity's name
	Type          int                  `json:"type"`                     // activity type
	URL           *string              `json:"url,omitempty"`            //stream url, is validated when type is 1
	Timestamps    []*ActivityTimestamp `json:"timestamps,omitempty"`     // timestamps object	unix timestamps for start and/or end of the game
	ApplicationID Snowflake            `json:"application_id,omitempty"` //?	snowflake	application id for the game
	Details       *string              `json:"details,omitempty"`        //?	?string	what the player is currently doing
	State         *string              `json:"state,omitempty"`          //state?	?string	the user's current party status
	Party         *ActivityParty       `json:"party"`                    //party?	party object	information for the current party of the player
	Assets        *ActivityAssets      `json:"assets,omitempty"`         // assets?	assets object	images for the presence and their hover texts
	Secrets       *ActivitySecrets     `json:"secrets,omitempty"`        // secrets?	secrets object	secrets for Rich Presence joining and spectating
	Instance      bool                 `json:"instance,omitempty"`       // instance?	boolean	whether or not the activity is an instanced game session
	Flags         int                  `json:"flags,omitempty"`          // flags?	int	activity flags ORd together, describes what the payload includes
}

// ---------

func NewUser() *User {
	return &User{}
}

type User struct {
	ID            Snowflake `json:"id,omitempty"`
	Username      string    `json:"username,omitempty"`
	Discriminator string    `json:"discriminator,omitempty"`
	Email         string    `json:"email,omitempty"`
	Avatar        *string   `json:"avatar"` // data:image/jpeg;base64,BASE64_ENCODED_JPEG_IMAGE_DATA //TODO: pointer?
	Token         string    `json:"token,omitempty"`
	Verified      bool      `json:"verified,omitempty"`
	MFAEnabled    bool      `json:"mfa_enabled,omitempty"`
	Bot           bool      `json:"bot,omitempty"`

	sync.RWMutex `json:"-"`
}

func (u *User) Mention() string {
	return "<@" + u.ID.String() + ">"
}

func (u *User) MentionNickname() string {
	return "<@!" + u.ID.String() + ">"
}

func (u *User) String() string {
	return u.Username + "#" + u.Discriminator + "{" + u.ID.String() + "}"
}

// Partial check if this is not a complete user object
// Assumption: has a snowflake.
func (u *User) Partial() bool {
	return (u.Username + u.Discriminator) == ""
}

func (u *User) MarshalJSON() ([]byte, error) {
	if u.ID.Empty() {
		return []byte("{}"), nil
	}

	return json.Marshal(User(*u))
}

// func (u *User) UnmarshalJSON(data []byte) error {
// 	return json.Unmarshal(data, &u.userJSON)
// }

func (u *User) Clear() {
	//u.d.Avatar = nil
}

func (u *User) SendMsg(session Session, message *Message) (channel *Channel, msg *Message, err error) {
	channel, err = session.CreateDM(u.ID)
	if err != nil {
		return
	}

	msg, err = session.SendMsg(channel.ID, message)
	return
}

func (u *User) SendMsgString(session Session, content string) (channel *Channel, msg *Message, err error) {
	channel, msg, err = u.SendMsg(session, &Message{
		Content: content,
	})
	return
}

func (u *User) DeepCopy() (copy interface{}) {
	copy = NewUser()
	u.CopyOverTo(copy)

	return
}

func (u *User) CopyOverTo(other interface{}) (err error) {
	var user *User
	var valid bool
	if user, valid = other.(*User); !valid {
		err = NewErrorUnsupportedType("argument given is not a *User type")
		return
	}

	u.RLock()
	user.Lock()

	user.ID = u.ID
	user.Username = u.Username
	user.Discriminator = u.Discriminator
	user.Email = u.Email
	user.Token = u.Token
	user.Verified = u.Verified
	user.MFAEnabled = u.MFAEnabled
	user.Bot = u.Bot

	if u.Avatar != nil {
		avatar := *u.Avatar
		user.Avatar = &avatar
	}

	u.RUnlock()
	user.Unlock()

	return
}

func (u *User) SaveToDiscord(session Session) (err error) {
	// TODO: check snowflake if ID is current user
	// call both modify methods
	return errors.New("not implemented")
}

func (u *User) Valid() bool {
	return u.ID > 0
}

// -------

func NewUserPresence() *UserPresence {
	return &UserPresence{}
}

type UserPresence struct {
	User    *User         `json:"user"`
	Roles   []Snowflake   `json:"roles"`
	Game    *UserActivity `json:"activity"`
	GuildID Snowflake     `json:"guild_id"`
	Nick    string        `json:"nick"`
	Status  string        `json:"status"`
}

func (p *UserPresence) Update(status string) {
	// Update the presence.
	// talk to the discord api
}

func (p *UserPresence) String() string {
	return p.Status
}

func (p *UserPresence) Clear() {
	p.Game = nil
}

type UserConnection struct {
	ID           string                `json:"id"`           // id of the connection account
	Name         string                `json:"name"`         // the username of the connection account
	Type         string                `json:"type"`         // the service of the connection (twitch, youtube)
	Revoked      bool                  `json:"revoked"`      // whether the connection is revoked
	Integrations []*IntegrationAccount `json:"integrations"` // an array of partial server integrations
}
