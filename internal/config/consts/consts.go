package consts

const AppName = "go-blog"

// App consts
const (
	LongTextLength    = 32767 //  int(int16(^uint16(0) >> 1)) // equivalent of short.MaxValue
	DefaultTextLength = 100

	TitleTextLengthTiny   = 12
	TitleTextLengthSmall  = 25
	TitleTextLengthInfo   = 35
	TitleTextLengthMedium = 50
	TitleTextLengthLarge  = 100
)

const (
	// PathAPI represents the group of PathAPI.
	PathAPI        = "/api"
	PathAuthSignin = "/auth/signin"
)
const (
	RoleAdmin = "admin"
)

//nolint:gosec
const (
	PathSysMetricsAPI = "/sys/api/metrics"

	PathBlogPingDebugAPI = "/blog/api/ping"

	PathBlog       = "/blog"
	PathBlogAssets = "/blog/assets"

	PathBlogPostsEntity          = "/blog/posts/:code"          //
	PathBlogStatusAPI            = "/blog/api/status"           // get _csrf, user related, no-cache
	PathBlogConfigAPI            = "/blog/api/config"           // public
	PathBlogPostsEntityByCodeAPI = "/blog/api/posts/:code/code" //

)
