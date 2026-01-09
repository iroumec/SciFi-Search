package notifications

// ------------------------------------------------------------------

type Notifications struct {
	Name string
}

// ------------------------------------------------------------------

var (
	CookieNotification    = Notifications{Name: "Cookie"}
	HXTriggerNotification = Notifications{Name: "HXTrigger"}
)

// ------------------------------------------------------------------
