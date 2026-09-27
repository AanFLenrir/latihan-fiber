package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// POST — semua field wajib (Tambahan tag validate dari Langkah 6)
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email,max=120"`
	Password string `json:"password" validate:"required,min=8,max=72,nospace"`
}

// PUT — ganti seluruh isi, jadi field bertipe biasa dan semuanya wajib (Tambahan tag validate dari Langkah 6)
type ReplaceUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email,max=120"`
	IsActive bool   `json:"is_active"`
}

// PATCH — ubah sebagian, jadi field bertipe pointer supaya bisa dibedakan
// antara "tidak dikirim" (nil) dan "dikirim bernilai kosong"
type PatchUserRequest struct {
	// CATATAN LAPORAN: Di modul ditulis 'Username string'. Itu adalah KESALAHAN DISENGAJA.
	// Harus tetap '*string' agar tag omitnil berfungsi benar dan tidak menerima string kosong.
	Username *string `json:"username,omitempty" validate:"omitnil,min=3,max=30,alphanum"`
	Email    *string `json:"email,omitempty" validate:"omitnil,email,max=120"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// Amplop baku untuk respons sukses
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// BENTUK RESPONSE KEGAGALAN BARU (Tambahan dari Langkah 1)
type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}

// Offset menghitung berapa baris yang dilewati untuk halaman ini.
func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

// BARU DITAMBAHKAN: AssignRoleRequest dipakai endpoint PATCH /users/:id/role.
// Saya tambahkan tag validate required agar tidak bisa dikirim kosong.
type AssignRoleRequest struct {
	Role string `json:"role" validate:"required"`
}