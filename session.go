package mcgoweb

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// SessionDuration is used to set the length of a session
// before it expires
var SessionDuration time.Duration = 3 * 24 * time.Hour

// SessionUpdateWindow is used to determine when a session's
// expiration should be updated
var SessionUpdateWindow time.Duration = 12 * time.Hour

// SessionId represents a UUID session id.
type SessionId [16]byte

// Session represents a key-value session with
// expiration.
type Session struct {
	id SessionId
	values map[string]string
	expiration time.Time
	cache SessionCache
}

// SessionCache provides storage of sessions based
// off the session's SessionId.
type SessionCache interface {
	Retrieve(sessionId SessionId) (*Session, error)
	Store(sessionId SessionId, session *Session) error
	Delete(sessionId SessionId) error
}

// String returns the UUID string version of a SessionId.
func (id SessionId) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
}

// SessionIdFromString returns a SessionId parsed
// from the given string.
func SessionIdFromString(sessionString string) (SessionId, error) {
	var id SessionId
	n, err := fmt.Sscanf(sessionString, "%2x%2x%2x%2x-%2x%2x-%2x%2x-%2x%2x-%2x%2x%2x%2x%2x%2x", 
		&id[0], &id[1], &id[2], &id[3],
		&id[4], &id[5],
		&id[6], &id[7],
		&id[8], &id[9],
		&id[10], &id[11], &id[12], &id[13], &id[14], &id[15])
	if err != nil {
		return id, err
	}
	if n != 16 {
		panic("Invalid scan of session id")
	}
	return id, nil
}

// NewSessionId returns a new randomly generated SessionId
func NewSessionId() SessionId {
	var id SessionId
	for i := 0; i < 16; i++ {
		id[i] = byte(rand.Intn(255))
	}
	return id
}

// Store saves the session to the cache
func (session *Session) Store() error {
	return session.cache.Store(session.id, session)
}

// UpdateValue updates the value for given key and saves
// the result back to its cache.
func (session *Session) UpdateValue(key, value string) error {
	session.values[key] = value
	return session.cache.Store(session.id, session)
}

// GetValue returns the value for the key
func (session *Session) GetValue(key string) (string, bool) {
	value, ok := session.values[key]
	return value, ok
}

// Expire updates a session expiration to now and removes
// the session from the cache.
func (session *Session) Expire() error {
	session.expiration = time.Now()
	return session.cache.Delete(session.id)
}

// GetSessionKey returns the string representation of
// the session's identifier. 
func (session *Session) GetSessionKey() string {
	return session.id.String()
}

// NewUsersession returns a new Session given a username,
// duration until expiration, and the cache to store the session.
func NewUserSession(user string, cache SessionCache) *Session {
	session := new(Session)
	session.id = NewSessionId()
	session.values = make(map[string]string)
	session.values["user"] = user
	session.expiration = time.Now().Add(SessionDuration)
	session.cache = cache
	return session
}

// GetSession returns a session from the given cache using
// the provided key for lookup.
func GetSession(key string, cache SessionCache) *Session {
	id, err := SessionIdFromString(key)
	if err != nil {
		return nil
	}
	session, err := cache.Retrieve(id)
	if err != nil {
		return nil
	}
	return session
}

// MemorySessionCache provides a SessionCache using an 
// in-memory object.  Sessions will not be persisted
// when an application goes offline.
type MemorySessionCache struct {
	sessions map[SessionId]*Session
}

func NewMemorySessionCache() SessionCache {
	cache := new(MemorySessionCache)
	cache.sessions = make(map[SessionId]*Session)
	return cache
}

func (cache *MemorySessionCache) Retrieve(sessionId SessionId) (*Session, error) {
	session, ok := cache.sessions[sessionId]
	if !ok {
		return nil, nil
	}
	return session, nil
}

func (cache *MemorySessionCache) Store(sessionId SessionId, session *Session) error {
	cache.sessions[sessionId] = session
	return nil
}

func (cache *MemorySessionCache) Delete(sessionId SessionId) error {
	delete(cache.sessions, sessionId)
	return nil
}

func SessionMiddleware(handler RequestHandler, context *RequestContext) {
	cookie, err := context.Request.Cookie("SID")
	if err == nil && len(cookie.Value) > 0 {
		session := GetSession(cookie.Value, context.sessionCache)
		if session == nil {
			cookie.Value = ""
			cookie.Expires = time.Unix(0,0)
			http.SetCookie(context.Writer,cookie)
		} else if time.Now().After(session.expiration) {
			session.Expire()
			cookie.Value = ""
			cookie.Expires = time.Unix(0,0)
			http.SetCookie(context.Writer,cookie)
		} else {
			context.Session = session
			if time.Now().After(session.expiration.Add(SessionUpdateWindow)) {
				session.expiration = time.Now().Add(SessionDuration)
				cookie.Expires = session.expiration
				session.Store()
				http.SetCookie(context.Writer,cookie)
			}
		}
	}
	handler(context)
}
