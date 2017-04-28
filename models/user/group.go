package user

import (
	"time"

	types "github.com/golang/protobuf/ptypes"
	"github.com/satori/go.uuid"
)

// NewGroup returns Group
func NewGroup() *Group {
	now, _ := types.TimestampProto(time.Now())
	return &Group{
		Id:      uuid.NewV4().Bytes(),
		Created: now,
		Updated: now,
	}
}

// ID returns uuid
func (g *Group) ID() uuid.UUID {
	return uuid.FromBytesOrNil(g.Id)
}
