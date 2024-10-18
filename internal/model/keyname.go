package model

const ConcatSign string = ":"

// 用户序列

const UserClickRedisKeyName string = "user_click_seq"
const UserQueryClickRedisKeyName string = "user_query_click_seq"
const UserIntentListRedisKeyName string = "user_intent_list"
const UserLikeSeqRedisKeyName string = "user_like_seq"
const UserViewRedisKeyName string = "user_view_seq"
const UserDisLikeNoteSeqRedisKeyName string = "user_dislike_note_seq"
const UserDisLikeAuthorSeqRedisKeyName string = "user_dislike_author_seq"
const UserDisLikeTagsSeqRedisKeyName string = "user_dislike_tags_seq"
const UserFollowSeqRedisKeyName string = "user_follow_seq"
const UserForbiddenKeyName string = "forbidden:user"
const UserCFListKeyName string = "user_collaborative_filtering_list"
const UserFromChannelName string = "user_from_channel_name"

// 帖子序列

const NoteTopNewRedisKeyName string = "new_note_rank"
const NoteTopHotRedisKeyName string = "hot_note_rank"
const NoteForbiddenKeyName string = "forbidden:note"

// 实验配置

const AbTestKeyName string = "exp:abtest:recommend"
const AbTestWhiteListKeyName string = "exp:abtest:recommend:whitelist"
