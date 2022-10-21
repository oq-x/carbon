package framework

import (
	"reflect"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
)

var perms = make(map[string]int64)
var Permissions = struct {
	CREATE_INSTANT_INVITE      int64
	KICK_MEMBERS               int64
	BAN_MEMBERS                int64
	ADMINISTRATOR              int64
	MANAGE_CHANNELS            int64
	MANAGE_GUILD               int64
	ADD_REACTIONS              int64
	VIEW_AUDIT_LOG             int64
	PRIORITY_SPEAKER           int64
	STREAM                     int64
	VIEW_CHANNEL               int64
	SEND_MESSAGES              int64
	SEND_TTS_MESSAGES          int64
	MANAGE_MESSAGES            int64
	EMBED_LINKS                int64
	ATTACH_FILES               int64
	READ_MESSAGE_HISTORY       int64
	MENTION_EVERYONE           int64
	USE_EXTERNAL_EMOJIS        int64
	VIEW_GUILD_INSIGHTS        int64
	CONNECT                    int64
	SPEAK                      int64
	MUTE_MEMBERS               int64
	DEAFEN_MEMBERS             int64
	MOVE_MEMBERS               int64
	USE_VAD                    int64
	CHANGE_NICKNAME            int64
	MANAGE_NICKNAMES           int64
	MANAGE_ROLES               int64
	MANAGE_WEBHOOKS            int64
	MANAGE_EMOJIS_AND_STICKERS int64
	USE_SLASH_COMMANDS         int64
	REQUEST_TO_SPEAK           int64
	MANAGE_EVENTS              int64
	MANAGE_THREADS             int64
	CREATE_PUBLIC_THREADS      int64
	CREATE_PRIVATE_THREADS     int64
	USE_EXTERNAL_STICKERS      int64
	SEND_MESSAGES_IN_THREADS   int64
	USE_EMBEDDED_ACTIVITIES    int64
	MODERATE_MEMBERS           int64
}{
	CREATE_INSTANT_INVITE:      0x0000000000000001,
	KICK_MEMBERS:               0x0000000000000002,
	BAN_MEMBERS:                0x0000000000000004,
	ADMINISTRATOR:              0x0000000000000008,
	MANAGE_CHANNELS:            0x0000000000000010,
	MANAGE_GUILD:               0x0000000000000020,
	ADD_REACTIONS:              0x0000000000000040,
	VIEW_AUDIT_LOG:             0x0000000000000080,
	PRIORITY_SPEAKER:           0x0000000000000100,
	STREAM:                     0x0000000000000200,
	VIEW_CHANNEL:               0x0000000000000400,
	SEND_MESSAGES:              0x0000000000000800,
	SEND_TTS_MESSAGES:          0x0000000000001000,
	MANAGE_MESSAGES:            0x0000000000002000,
	EMBED_LINKS:                0x0000000000004000,
	ATTACH_FILES:               0x0000000000008000,
	READ_MESSAGE_HISTORY:       0x0000000000010000,
	MENTION_EVERYONE:           0x0000000000020000,
	USE_EXTERNAL_EMOJIS:        0x0000000000040000,
	VIEW_GUILD_INSIGHTS:        0x0000000000080000,
	CONNECT:                    0x0000000000100000,
	SPEAK:                      0x0000000000200000,
	MUTE_MEMBERS:               0x0000000000400000,
	DEAFEN_MEMBERS:             0x0000000000800000,
	MOVE_MEMBERS:               0x0000000001000000,
	USE_VAD:                    0x0000000002000000,
	CHANGE_NICKNAME:            0x0000000004000000,
	MANAGE_NICKNAMES:           0x0000000008000000,
	MANAGE_ROLES:               0x0000000010000000,
	MANAGE_WEBHOOKS:            0x0000000020000000,
	MANAGE_EMOJIS_AND_STICKERS: 0x0000000040000000,
	USE_SLASH_COMMANDS:         0x0000000080000000,
	REQUEST_TO_SPEAK:           0x0000000100000000,
	MANAGE_EVENTS:              0x0000000200000000,
	MANAGE_THREADS:             0x0000000400000000,
	CREATE_PUBLIC_THREADS:      0x0000000800000000,
	CREATE_PRIVATE_THREADS:     0x0000001000000000,
	USE_EXTERNAL_STICKERS:      0x0000002000000000,
	SEND_MESSAGES_IN_THREADS:   0x0000004000000000,
	USE_EMBEDDED_ACTIVITIES:    0x0000008000000000,
	MODERATE_MEMBERS:           0x0000010000000000,
}

func CalculateBasePermissions(guild discordgo.Guild, roles []string) int64 {
	var permissions int64 = 0

	for _, i := range roles {
		for _, o := range guild.Roles {
			if o.ID != i {
				continue
			}
			permissions |= o.Permissions
		}
	}
	return permissions
}

func CalculatePermissions(permissionBits int64) []string {

	var permissions []string
	var allPermissions []string
	v := reflect.ValueOf(Permissions)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		perms[typeOfS.Field(i).Name] = v.Field(i).Interface().(int64)
		allPermissions = append(allPermissions, typeOfS.Field(i).Name)
	}

	for _, i := range allPermissions {
		if permissionBits&perms[i] != 0 {
			permissions = append(permissions, i)
		}
	}

	return permissions
}
func GetHighestRole(client bot.Client, guild discord.Guild, member *discord.Member) discord.Role {
	memberRoles := client.Caches().MemberRoles(*member)
	var everyoneRole discord.Role
	if len(memberRoles) == 0 {
		return everyoneRole
	}
	var memberHighestRole discord.Role
	for _, i := range memberRoles {
		for _, role := range client.Caches().Roles().GroupAll(guild.ID) {
			if i.ID != role.ID {
				continue
			}
			if memberHighestRole.Position < role.Position || memberHighestRole.Position == role.Position {
				memberHighestRole = role
			}
		}
	}

	return memberHighestRole
}

func HigherRolePosition(
	guild discord.Guild,
	role discord.Role,
	otherRole discord.Role,
) bool {
	if role.Position == otherRole.Position {
		return role.ID < otherRole.ID
	}

	return role.Position > otherRole.Position
}

func HigherMember(client bot.Client, guild discord.Guild, firstMember *discord.Member, secondMember *discord.Member) string {
	if IsGuildOwner(guild, firstMember) {
		return firstMember.User.ID.String()
	} else if IsGuildOwner(guild, secondMember) {
		return secondMember.User.ID.String()
	} else {
		firstMemberHighestRole := GetHighestRole(client, guild, firstMember)
		secondMemberHighestRole := GetHighestRole(client, guild, secondMember)

		if HigherRolePosition(guild, firstMemberHighestRole, secondMemberHighestRole) {
			return firstMember.User.ID.String()
		} else if HigherRolePosition(guild, secondMemberHighestRole, firstMemberHighestRole) {
			return secondMember.User.ID.String()
		}
	}
	return ""
}

func IsGuildOwner(guild discord.Guild, member *discord.Member) bool {
	return guild.OwnerID == member.User.ID
}
