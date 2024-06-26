package builtins

import (
	"github.com/bwmarrin/discordgo"
)

var Constants = map[string]any{
	"Permissions": map[string]int64{
		"ReadMessages":          discordgo.PermissionViewChannel,
		"SendMessages":          discordgo.PermissionSendMessages,
		"SendTTSMessages":       discordgo.PermissionSendTTSMessages,
		"ManageMessages":        discordgo.PermissionManageMessages,
		"EmbedLinks":            discordgo.PermissionEmbedLinks,
		"AttachFiles":           discordgo.PermissionAttachFiles,
		"ReadMessageHistory":    discordgo.PermissionReadMessageHistory,
		"MentionEveryone":       discordgo.PermissionMentionEveryone,
		"UseExternalEmojis":     discordgo.PermissionUseExternalEmojis,
		"UseSlashCommands":      discordgo.PermissionUseSlashCommands,
		"ManageThreads":         discordgo.PermissionManageThreads,
		"CreatePublicThreads":   discordgo.PermissionCreatePublicThreads,
		"CreatePrivateThreads":  discordgo.PermissionCreatePrivateThreads,
		"UseExternalStickers":   discordgo.PermissionUseExternalStickers,
		"SendMessagesInThreads": discordgo.PermissionSendMessagesInThreads,
		"VoicePrioritySpeaker":  discordgo.PermissionVoicePrioritySpeaker,
		"VoiceStreamVideo":      discordgo.PermissionVoiceStreamVideo,
		"VoiceConnect":          discordgo.PermissionVoiceConnect,
		"VoiceSpeak":            discordgo.PermissionVoiceSpeak,
		"VoiceMuteMembers":      discordgo.PermissionVoiceMuteMembers,
		"VoiceDeafenMembers":    discordgo.PermissionVoiceDeafenMembers,
		"VoiceMoveMembers":      discordgo.PermissionVoiceMoveMembers,
		"VoiceUseVAD":           discordgo.PermissionVoiceUseVAD,
		"VoiceRequestToSpeak":   discordgo.PermissionVoiceRequestToSpeak,
		"UseActivities":         discordgo.PermissionUseActivities,
		"ChangeNickname":        discordgo.PermissionChangeNickname,
		"ManageNicknames":       discordgo.PermissionManageNicknames,
		"ManageRoles":           discordgo.PermissionManageRoles,
		"ManageWebhooks":        discordgo.PermissionManageWebhooks,
		"ManageEmojis":          discordgo.PermissionManageEmojis,
		"ManageEvents":          discordgo.PermissionManageEvents,
		"CreateInstantInvite":   discordgo.PermissionCreateInstantInvite,
		"KickMembers":           discordgo.PermissionKickMembers,
		"BanMembers":            discordgo.PermissionBanMembers,
		"Administrator":         discordgo.PermissionAdministrator,
		"ManageChannels":        discordgo.PermissionManageChannels,
		"ManageServer":          discordgo.PermissionManageServer,
		"AddReactions":          discordgo.PermissionAddReactions,
		"ViewAuditLogs":         discordgo.PermissionViewAuditLogs,
		"ViewChannel":           discordgo.PermissionViewChannel,
		"ViewGuildInsights":     discordgo.PermissionViewGuildInsights,
		"ModerateMembers":       discordgo.PermissionModerateMembers,
		"AllText":               discordgo.PermissionAllText,
		"AllVoice":              discordgo.PermissionAllVoice,
		"AllChannel":            discordgo.PermissionAllChannel,
		"All":                   discordgo.PermissionAll,
	},
	"MessageFlag": map[string]discordgo.MessageFlags{
		"CrossPosted":                      discordgo.MessageFlagsCrossPosted,
		"IsCrossPosted":                    discordgo.MessageFlagsIsCrossPosted,
		"SuppressEmbeds":                   discordgo.MessageFlagsSuppressEmbeds,
		"SupressEmbeds":                    discordgo.MessageFlagsSupressEmbeds,
		"SourceMessageDeleted":             discordgo.MessageFlagsSourceMessageDeleted,
		"Urgent":                           discordgo.MessageFlagsUrgent,
		"HasThread":                        discordgo.MessageFlagsHasThread,
		"Ephemeral":                        discordgo.MessageFlagsEphemeral,
		"Loading":                          discordgo.MessageFlagsLoading,
		"FailedToMentionSomeRolesInThread": discordgo.MessageFlagsFailedToMentionSomeRolesInThread,
		"SuppressNotifications":            discordgo.MessageFlagsSuppressNotifications,
		"IsVoiceMessage":                   discordgo.MessageFlagsIsVoiceMessage,
	},
	"MessageType": map[string]discordgo.MessageType{
		"Default":                               discordgo.MessageTypeDefault,
		"RecipientAdd":                          discordgo.MessageTypeRecipientAdd,
		"RecipientRemove":                       discordgo.MessageTypeRecipientRemove,
		"Call":                                  discordgo.MessageTypeCall,
		"ChannelNameChange":                     discordgo.MessageTypeChannelNameChange,
		"ChannelIconChange":                     discordgo.MessageTypeChannelIconChange,
		"ChannelPinnedMessage":                  discordgo.MessageTypeChannelPinnedMessage,
		"GuildMemberJoin":                       discordgo.MessageTypeGuildMemberJoin,
		"UserPremiumGuildSubscription":          discordgo.MessageTypeUserPremiumGuildSubscription,
		"UserPremiumGuildSubscriptionTierOne":   discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne,
		"UserPremiumGuildSubscriptionTierTwo":   discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo,
		"UserPremiumGuildSubscriptionTierThree": discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree,
		"ChannelFollowAdd":                      discordgo.MessageTypeChannelFollowAdd,
		"GuildDiscoveryDisqualified":            discordgo.MessageTypeGuildDiscoveryDisqualified,
		"GuildDiscoveryRequalified":             discordgo.MessageTypeGuildDiscoveryRequalified,
		"ThreadCreated":                         discordgo.MessageTypeThreadCreated,
		"Reply":                                 discordgo.MessageTypeReply,
		"ChatInputCommand":                      discordgo.MessageTypeChatInputCommand,
		"ThreadStarterMessage":                  discordgo.MessageTypeThreadStarterMessage,
		"ContextMenuCommand":                    discordgo.MessageTypeContextMenuCommand,
	},
	"Status": map[string]discordgo.Status{
		"Online":       discordgo.StatusOnline,
		"Idle":         discordgo.StatusIdle,
		"DoNotDisturb": discordgo.StatusDoNotDisturb,
		"Invisible":    discordgo.StatusInvisible,
		"Offline":      discordgo.StatusOffline,
	},
	"UserFlags": map[string]discordgo.UserFlags{
		"DiscordEmployee":           discordgo.UserFlagDiscordEmployee,
		"DiscordPartner":            discordgo.UserFlagDiscordPartner,
		"HypeSquadEvents":           discordgo.UserFlagHypeSquadEvents,
		"BugHunterLevel1":           discordgo.UserFlagBugHunterLevel1,
		"HouseBravery":              discordgo.UserFlagHouseBravery,
		"HouseBrilliance":           discordgo.UserFlagHouseBrilliance,
		"HouseBalance":              discordgo.UserFlagHouseBalance,
		"EarlySupporter":            discordgo.UserFlagEarlySupporter,
		"TeamUser":                  discordgo.UserFlagTeamUser,
		"System":                    discordgo.UserFlagSystem,
		"BugHunterLevel2":           discordgo.UserFlagBugHunterLevel2,
		"VerifiedBot":               discordgo.UserFlagVerifiedBot,
		"VerifiedBotDeveloper":      discordgo.UserFlagVerifiedBotDeveloper,
		"DiscordCertifiedModerator": discordgo.UserFlagDiscordCertifiedModerator,
		"BotHTTPInteractions":       discordgo.UserFlagBotHTTPInteractions,
		"ActiveBotDeveloper":        discordgo.UserFlagActiveBotDeveloper,
	},
	"RoleFlags": map[string]discordgo.RoleFlags{
		"InPrompt": discordgo.RoleFlagInPrompt,
	},
	"SelectMenuType": map[string]discordgo.SelectMenuType{
		"String":      discordgo.StringSelectMenu,
		"User":        discordgo.UserSelectMenu,
		"Role":        discordgo.RoleSelectMenu,
		"Mentionable": discordgo.MentionableSelectMenu,
		"Channel":     discordgo.ChannelSelectMenu,
	},
	"ComponentType": map[string]discordgo.ComponentType{
		"ActionsRow":            discordgo.ActionsRowComponent,
		"Button":                discordgo.ButtonComponent,
		"SelectMenu":            discordgo.SelectMenuComponent,
		"TextInput":             discordgo.TextInputComponent,
		"UserSelectMenu":        discordgo.UserSelectMenuComponent,
		"RoleSelectMenu":        discordgo.RoleSelectMenuComponent,
		"MentionableSelectMenu": discordgo.MentionableSelectMenuComponent,
		"ChannelSelectMenu":     discordgo.ChannelSelectMenuComponent,
	},
	"EmbedType": map[string]discordgo.EmbedType{
		"Rich":    discordgo.EmbedTypeRich,
		"Image":   discordgo.EmbedTypeImage,
		"Video":   discordgo.EmbedTypeVideo,
		"Gifv":    discordgo.EmbedTypeGifv,
		"Article": discordgo.EmbedTypeArticle,
		"Link":    discordgo.EmbedTypeLink,
	},
	"MfaLevel": map[string]discordgo.MfaLevel{
		"None":     discordgo.MfaLevelNone,
		"Elevated": discordgo.MfaLevelElevated,
	},
	"PermissionOverwriteType": map[string]discordgo.PermissionOverwriteType{
		"Role":   discordgo.PermissionOverwriteTypeRole,
		"Member": discordgo.PermissionOverwriteTypeMember,
	},
	"PremiumTier": map[string]discordgo.PremiumTier{
		"None":  discordgo.PremiumTierNone,
		"Tier1": discordgo.PremiumTier1,
		"Tier2": discordgo.PremiumTier2,
		"Tier3": discordgo.PremiumTier3,
	},
	"SelectMenuDefaultValueType": map[string]discordgo.SelectMenuDefaultValueType{
		"User":    discordgo.SelectMenuDefaultValueUser,
		"Role":    discordgo.SelectMenuDefaultValueRole,
		"Channel": discordgo.SelectMenuDefaultValueChannel,
	},
	"StageInstancePrivacyLevel": map[string]discordgo.StageInstancePrivacyLevel{
		"Public":    discordgo.StageInstancePrivacyLevelPublic,
		"GuildOnly": discordgo.StageInstancePrivacyLevelGuildOnly,
	},
	"StickerFormat": map[string]discordgo.StickerFormat{
		"PNG":    discordgo.StickerFormatTypePNG,
		"APNG":   discordgo.StickerFormatTypeAPNG,
		"Lottie": discordgo.StickerFormatTypeLottie,
		"GIF":    discordgo.StickerFormatTypeGIF,
	},
	"StickerType": map[string]discordgo.StickerType{
		"Standard": discordgo.StickerTypeStandard,
		"Guild":    discordgo.StickerTypeGuild,
	},
	"ExpireBehavior": map[string]discordgo.ExpireBehavior{
		"RemoveRole": discordgo.ExpireBehaviorRemoveRole,
		"Kick":       discordgo.ExpireBehaviorKick,
	},
	"ExplicitContentFilterLevel": map[string]discordgo.ExplicitContentFilterLevel{
		"Disabled":            discordgo.ExplicitContentFilterDisabled,
		"MembersWithoutRoles": discordgo.ExplicitContentFilterMembersWithoutRoles,
		"AllMembers":          discordgo.ExplicitContentFilterAllMembers,
	},
	"ForumLayout": map[string]discordgo.ForumLayout{
		"NotSet":      discordgo.ForumLayoutNotSet,
		"ListView":    discordgo.ForumLayoutListView,
		"GalleryView": discordgo.ForumLayoutGalleryView,
	},
	"ForumSortOrderType": map[string]discordgo.ForumSortOrderType{
		"LatestActivity": discordgo.ForumSortOrderLatestActivity,
		"CreationDate":   discordgo.ForumSortOrderCreationDate,
	},
	"GuildFeature": map[string]discordgo.GuildFeature{
		"AnimatedBanner":                discordgo.GuildFeatureAnimatedBanner,
		"AnimatedIcon":                  discordgo.GuildFeatureAnimatedIcon,
		"AutoModeration":                discordgo.GuildFeatureAutoModeration,
		"Banner":                        discordgo.GuildFeatureBanner,
		"Community":                     discordgo.GuildFeatureCommunity,
		"Discoverable":                  discordgo.GuildFeatureDiscoverable,
		"Featurable":                    discordgo.GuildFeatureFeaturable,
		"InviteSplash":                  discordgo.GuildFeatureInviteSplash,
		"MemberVerificationGateEnabled": discordgo.GuildFeatureMemberVerificationGateEnabled,
		"MonetizationEnabled":           discordgo.GuildFeatureMonetizationEnabled,
		"MoreStickers":                  discordgo.GuildFeatureMoreStickers,
		"News":                          discordgo.GuildFeatureNews,
		"Partnered":                     discordgo.GuildFeaturePartnered,
		"PreviewEnabled":                discordgo.GuildFeaturePreviewEnabled,
		"PrivateThreads":                discordgo.GuildFeaturePrivateThreads,
		"RoleIcons":                     discordgo.GuildFeatureRoleIcons,
		"TicketedEventsEnabled":         discordgo.GuildFeatureTicketedEventsEnabled,
		"VanityURL":                     discordgo.GuildFeatureVanityURL,
		"Verified":                      discordgo.GuildFeatureVerified,
		"VipRegions":                    discordgo.GuildFeatureVipRegions,
		"WelcomeScreenEnabled":          discordgo.GuildFeatureWelcomeScreenEnabled,
	},
	"GuildNSFWLevel": map[string]discordgo.GuildNSFWLevel{
		"Default":       discordgo.GuildNSFWLevelDefault,
		"Explicit":      discordgo.GuildNSFWLevelExplicit,
		"Safe":          discordgo.GuildNSFWLevelSafe,
		"AgeRestricted": discordgo.GuildNSFWLevelAgeRestricted,
	},
	"GuildOnboardingMode": map[string]discordgo.GuildOnboardingMode{
		"Default":  discordgo.GuildOnboardingModeDefault,
		"Advanced": discordgo.GuildOnboardingModeAdvanced,
	},
	"GuildOnboardingPromptType": map[string]discordgo.GuildOnboardingPromptType{
		"MultipleChoice": discordgo.GuildOnboardingPromptTypeMultipleChoice,
		"Dropdown":       discordgo.GuildOnboardingPromptTypeDropdown,
	},
	"GuildScheduledEventEntityType": map[string]discordgo.GuildScheduledEventEntityType{
		"StageInstance": discordgo.GuildScheduledEventEntityTypeStageInstance,
		"Voice":         discordgo.GuildScheduledEventEntityTypeVoice,
		"External":      discordgo.GuildScheduledEventEntityTypeExternal,
	},
	"GuildScheduledEventPrivacyLevel": map[string]discordgo.GuildScheduledEventPrivacyLevel{
		"GuildOnly": discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
	},
	"GuildScheduledEventStatus": map[string]discordgo.GuildScheduledEventStatus{
		"Scheduled": discordgo.GuildScheduledEventStatusScheduled,
		"Active":    discordgo.GuildScheduledEventStatusActive,
		"Completed": discordgo.GuildScheduledEventStatusCompleted,
		"Canceled":  discordgo.GuildScheduledEventStatusCanceled,
	},
	"Intent": map[string]discordgo.Intent{
		"Guilds":                      discordgo.IntentGuilds,
		"GuildMembers":                discordgo.IntentGuildMembers,
		"GuildModeration":             discordgo.IntentGuildModeration,
		"GuildEmojis":                 discordgo.IntentGuildEmojis,
		"GuildIntegrations":           discordgo.IntentGuildIntegrations,
		"GuildWebhooks":               discordgo.IntentGuildWebhooks,
		"GuildInvites":                discordgo.IntentGuildInvites,
		"GuildVoiceStates":            discordgo.IntentGuildVoiceStates,
		"GuildPresences":              discordgo.IntentGuildPresences,
		"GuildMessages":               discordgo.IntentGuildMessages,
		"GuildMessageReactions":       discordgo.IntentGuildMessageReactions,
		"GuildMessageTyping":          discordgo.IntentGuildMessageTyping,
		"GuildBans":                   discordgo.IntentGuildBans,
		"DirectMessages":              discordgo.IntentDirectMessages,
		"DirectMessageReactions":      discordgo.IntentDirectMessageReactions,
		"DirectMessageTyping":         discordgo.IntentDirectMessageTyping,
		"MessageContent":              discordgo.IntentMessageContent,
		"GuildScheduledEvents":        discordgo.IntentGuildScheduledEvents,
		"AutoModerationConfiguration": discordgo.IntentAutoModerationConfiguration,
		"AutoModerationExecution":     discordgo.IntentAutoModerationExecution,
		"AllWithoutPrivileged":        discordgo.IntentsAllWithoutPrivileged,
		"IntentsAll":                  discordgo.IntentsAll,
		"IntentsNone":                 discordgo.IntentsNone,
	},
	"InteractionResponseType": map[string]discordgo.InteractionResponseType{
		"Pong":                                 discordgo.InteractionResponsePong,
		"ChannelMessageWithSource":             discordgo.InteractionResponseChannelMessageWithSource,
		"DeferredChannelMessageWithSource":     discordgo.InteractionResponseDeferredChannelMessageWithSource,
		"DeferredMessageUpdate":                discordgo.InteractionResponseDeferredMessageUpdate,
		"UpdateMessage":                        discordgo.InteractionResponseUpdateMessage,
		"ApplicationCommandAutocompleteResult": discordgo.InteractionApplicationCommandAutocompleteResult,
		"Modal":                                discordgo.InteractionResponseModal,
	},
	"InteractionType": map[string]discordgo.InteractionType{
		"Ping":                           discordgo.InteractionPing,
		"ApplicationCommand":             discordgo.InteractionApplicationCommand,
		"MessageComponent":               discordgo.InteractionMessageComponent,
		"ApplicationCommandAutocomplete": discordgo.InteractionApplicationCommandAutocomplete,
		"ModalSubmit":                    discordgo.InteractionModalSubmit,
	},
	"InviteTargetType": map[string]discordgo.InviteTargetType{
		"Stream":              discordgo.InviteTargetStream,
		"EmbeddedApplication": discordgo.InviteTargetEmbeddedApplication,
	},
	"Locale": map[string]discordgo.Locale{
		"EnglishUS":    discordgo.EnglishUS,
		"EnglishGB":    discordgo.EnglishGB,
		"Bulgarian":    discordgo.Bulgarian,
		"ChineseCN":    discordgo.ChineseCN,
		"ChineseTW":    discordgo.ChineseTW,
		"Croatian":     discordgo.Croatian,
		"Czech":        discordgo.Czech,
		"Danish":       discordgo.Danish,
		"Dutch":        discordgo.Dutch,
		"Finnish":      discordgo.Finnish,
		"French":       discordgo.French,
		"German":       discordgo.German,
		"Greek":        discordgo.Greek,
		"Hindi":        discordgo.Hindi,
		"Hungarian":    discordgo.Hungarian,
		"Italian":      discordgo.Italian,
		"Japanese":     discordgo.Japanese,
		"Korean":       discordgo.Korean,
		"Lithuanian":   discordgo.Lithuanian,
		"Norwegian":    discordgo.Norwegian,
		"Polish":       discordgo.Polish,
		"PortugueseBR": discordgo.PortugueseBR,
		"Romanian":     discordgo.Romanian,
		"Russian":      discordgo.Russian,
		"SpanishES":    discordgo.SpanishES,
		"SpanishLATAM": discordgo.SpanishLATAM,
		"Swedish":      discordgo.Swedish,
		"Thai":         discordgo.Thai,
		"Turkish":      discordgo.Turkish,
		"Ukrainian":    discordgo.Ukrainian,
		"Vietnamese":   discordgo.Vietnamese,
		"Unknown":      discordgo.Unknown,
	},
	"MemberFlags": map[string]discordgo.MemberFlags{
		"DidRejoin":            discordgo.MemberFlagDidRejoin,
		"CompletedOnboarding":  discordgo.MemberFlagCompletedOnboarding,
		"BypassesVerification": discordgo.MemberFlagBypassesVerification,
		"StartedOnboarding":    discordgo.MemberFlagStartedOnboarding,
	},
	"MembershipState": map[string]discordgo.MembershipState{
		"Invited":  discordgo.MembershipStateInvited,
		"Accepted": discordgo.MembershipStateAccepted,
	},
	"MessageActivityType": map[string]discordgo.MessageActivityType{
		"Join":        discordgo.MessageActivityTypeJoin,
		"Spectate":    discordgo.MessageActivityTypeSpectate,
		"Listen":      discordgo.MessageActivityTypeListen,
		"JoinRequest": discordgo.MessageActivityTypeJoinRequest,
	},
	"MessageNotifications": map[string]discordgo.MessageNotifications{
		"AllMessages":  discordgo.MessageNotificationsAllMessages,
		"OnlyMentions": discordgo.MessageNotificationsOnlyMentions,
	},
}
