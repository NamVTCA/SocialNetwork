package utils

import "regexp"

// Chuyển số điện thoại từ định dạng Việt Nam thành E.164 (ví dụ: 0912... => +84912...)
func FormatPhoneToE164(phone string) string {
	if len(phone) == 0 {
		return phone
	}
	if phone[0] == '0' {
		return "+84" + phone[1:]
	}
	if phone[0] != '+' {
		return "+" + phone
	}
	return phone
}

// Kiểm tra định dạng E.164 hợp lệ
func IsValidPhoneE164(phone string) bool {
	re := regexp.MustCompile(`^\+[1-9]\d{6,14}$`)
	return re.MatchString(phone)
}

func FormatPhoneToVietnamese(phone string) string {
	if len(phone) == 0 {
		return phone
	}
	if phone[0] == '0' {
		return "84" + phone[1:]
	}
	if phone[:3] == "+84" {
		return phone[1:] // +84 -> 84
	}
	return phone
}