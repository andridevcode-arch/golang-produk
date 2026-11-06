package cookie

import "net/http"

// Settings menampung atribut cookie autentikasi
type Settings struct {
	Name     string
	Path     string
	Domain   string
	Secure   bool
	SameSite http.SameSite
	MaxAge   int
}

// Set menulis cookie autentikasi sebagai HttpOnly cookie
func (s Settings) Set(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.Name,
		Value:    value,
		Path:     s.Path,
		Domain:   s.Domain,
		HttpOnly: true,
		Secure:   s.Secure,
		SameSite: s.SameSite,
		MaxAge:   maxAge,
	})
}

// Clear menghapus cookie autentikasi.
// Atribut path, domain, secure, dan sameSite harus sama dengan saat Set,
// jika tidak browser tidak akan menghapus cookie tersebut.
func (s Settings) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.Name,
		Value:    "",
		Path:     s.Path,
		Domain:   s.Domain,
		HttpOnly: true,
		Secure:   s.Secure,
		SameSite: s.SameSite,
		MaxAge:   -1,
	})
}

// FromRequest mengambil nilai token dari cookie.
// Mengembalikan string kosong jika cookie tidak ada atau nilainya kosong.
func (s Settings) FromRequest(r *http.Request) string {
	c, err := r.Cookie(s.Name)
	if err != nil || c == nil {
		return ""
	}
	return c.Value
}
