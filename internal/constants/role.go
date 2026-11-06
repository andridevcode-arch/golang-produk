package constants

// Role constants untuk pengguna
// Menambahkan role baru cukup tambahkan di sini, tidak perlu migrasi database
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// IsValidRole memeriksa apakah role valid
func IsValidRole(role string) bool {
	switch role {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}
