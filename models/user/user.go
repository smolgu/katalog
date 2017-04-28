package user

import (
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)

// New return pointer of new user
func New() *User {
	u := &User{}
	return u
}

// ID return uuid of user
func (u *User) ID() uuid.UUID {
	return uuid.FromBytesOrNil(u.Id)
}

// Encode marshall proto message
func (u *User) Encode() ([]byte, error) {
	return proto.Marshal(u)
}

// Encode marshall proto message
func (us *Users) Encode() ([]byte, error) {
	return proto.Marshal(us)
}
