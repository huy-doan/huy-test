package object

type UserAgent string

func (ua UserAgent) Browser() string {
	userAgent := string(ua)

	if contains(userAgent, "Chrome") && !contains(userAgent, "Chromium") {
		return "Chrome"
	} else if contains(userAgent, "Firefox") {
		return "Firefox"
	} else if contains(userAgent, "Safari") && !contains(userAgent, "Chrome") {
		return "Safari"
	} else if contains(userAgent, "Edge") {
		return "Edge"
	} else if contains(userAgent, "MSIE") || contains(userAgent, "Trident") {
		return "Internet Explorer"
	}

	return "Unknown"
}

func (ua UserAgent) OS() string {
	userAgent := string(ua)

	if contains(userAgent, "Windows") {
		return "Windows"
	} else if contains(userAgent, "Mac OS") {
		return "Mac OS"
	} else if contains(userAgent, "Linux") {
		return "Linux"
	} else if contains(userAgent, "Android") {
		return "Android"
	} else if contains(userAgent, "iOS") {
		return "iOS"
	}

	return "Unknown"
}

func (ua UserAgent) IsMobile() bool {
	userAgent := string(ua)

	return contains(userAgent, "Mobile") ||
		contains(userAgent, "Android") ||
		contains(userAgent, "iPhone") ||
		contains(userAgent, "iPad")
}

func (ua UserAgent) String() string {
	return string(ua)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
