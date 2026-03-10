package massage

func EntityTypeToTitle(entityType string) string {
	titles := map[string]string{
		"GROUP_INVITE":    "คำเชิญกลุ่มใหม่ 👥",
		"GROUP":           "กิจกรรมในกลุ่ม 👥",
		"GROUP_REMOVED":   "ออกจากกลุ่ม",
		"ASSET":           "รายการแชร์ใหม่ 💼",
		"FRIEND_REQUEST":  "คำขอเป็นเพื่อน 👋",
		"FRIEND_ACCEPTED": "ตอบรับเพื่อน ✅",
		"ACCESS_GRANTED":  "ได้รับสิทธิ์เข้าถึง 🔓",
		"INSURANCE":       "แจ้งเตือนประกัน ⏰",
	}
	if t, ok := titles[entityType]; ok {
		return t
	}
	return "แจ้งเตือนใหม่"
}
