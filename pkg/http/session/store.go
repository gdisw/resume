package session

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/meehow/securebytes"
)

type Key string

func init() {
	gob.Register(Key(""))
}

type Store interface {
	Put(*http.Request, http.ResponseWriter, Key, any) error
	PutAll(*http.Request, http.ResponseWriter, map[Key]any) error
	Get(*http.Request, Key) (any, error)
	GetAll(*http.Request, []Key) (map[Key]any, error)
	Delete(*http.Request, http.ResponseWriter) error
}

type store struct {
	sb   *securebytes.SecureBytes
	name string
}

func NewStore(sb *securebytes.SecureBytes, name string) *store {
	return &store{
		sb:   sb,
		name: name,
	}
}

func (s *store) read(r *http.Request) map[Key]any {
	session := make(map[Key]any)
	if cookie, err := r.Cookie(s.name); err == nil {
		_ = s.sb.DecryptBase64(cookie.Value, &session)
	}

	return session
}

func (s *store) save(session map[Key]any, w http.ResponseWriter) error {
	b64, err := s.sb.EncryptToBase64(session)
	if err != nil {
		log.Println(err)
		return err
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:    s.name,
		Value:   b64,
		Path:    "/",
		Expires: expiration,
		MaxAge:  int(time.Until(expiration).Seconds()),
	})
	return nil
}

func (s *store) Put(r *http.Request, w http.ResponseWriter, k Key, v any) error {
	session := s.read(r)

	session[k] = v
	return s.save(session, w)
}

func (s *store) PutAll(r *http.Request, w http.ResponseWriter, items map[Key]any) error {
	session := s.read(r)

	for k, v := range items {
		session[k] = v
	}

	return s.save(session, w)
}

func (s *store) Get(r *http.Request, k Key) (any, error) {
	session := s.read(r)

	return session[k], nil
}

func (s *store) GetAll(r *http.Request, ks []Key) (map[Key]any, error) {
	session := s.read(r)

	out := make(map[Key]any)
	for _, k := range ks {
		out[k] = session[k]
	}

	return out, nil
}

func (s *store) Delete(r *http.Request, w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Name:    s.name,
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	})
	return nil
}
