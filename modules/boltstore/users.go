package boltstore

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"github.com/smolgu/katalog/models/user"
	"github.com/zhuharev/boltutils"
)

var (
	// UsersBucketName reprezent name of boltdb bucket
	UsersBucketName = []byte("users")
)

// UserSave create or update user
func (bs *BoltStore) UserSave(u *user.User) error {
	if uuid.Equal(uuid.Nil, u.ID()) {
		u.Id = uuid.NewV4().Bytes()
	}
	if err := bs.checkUnique(u); err != nil {
		return err
	}
	data, err := u.Encode()
	if err != nil {
		return err
	}
	return bs.db.Put(UsersBucketName, u.Id, data)
}

// UsersGet return user by his id
func (bs *BoltStore) UsersGet(id []byte) (u *user.User, err error) {
	var data []byte
	data, err = bs.db.Get(UsersBucketName, id)
	if err != nil {
		return
	}
	u = new(user.User)
	err = proto.Unmarshal(data, u)
	if err != nil {
		return
	}
	return
}

// UsersGetByVkID return user by his vk_id
func (bs *BoltStore) UsersGetByVkID(vkID int64) (u *user.User, err error) {
	err = bs.db.Iterate(UsersBucketName, func(k []byte, v []byte) error {
		u = new(user.User)
		err = proto.Unmarshal(v, u)
		if err != nil {
			return err
		}
		if u.VkId == vkID {
			return boltutils.ErrBreak
		}
		return boltutils.ErrNotFound
	})
	if u == nil {
		return nil, boltutils.ErrNotFound
	}
	return
}

func (bs *BoltStore) checkUnique(u *user.User) error {
	return bs.db.Iterate(UsersBucketName, func(k []byte, v []byte) error {
		var (
			puser    = new(user.User)
			userUUID = uuid.FromBytesOrNil(u.Id)
		)
		err := proto.Unmarshal(v, puser)
		if err != nil {
			return err
		}
		puserUUID := uuid.FromBytesOrNil(puser.Id)
		if u.VkId != 0 && !uuid.Equal(userUUID, puserUUID) && u.VkId == puser.VkId {
			return fmt.Errorf("vk_id unique failed")
		}
		if u.EmployeeId != 0 && !uuid.Equal(userUUID, puserUUID) && u.EmployeeId == puser.EmployeeId {
			return fmt.Errorf("employee_id unique failed")
		}
		if u.LibraryId != 0 && !uuid.Equal(userUUID, puserUUID) && u.LibraryId == puser.LibraryId {
			return fmt.Errorf("library_id unique failed")
		}
		return nil
	})
}
