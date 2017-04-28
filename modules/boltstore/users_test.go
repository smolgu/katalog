package boltstore

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smolgu/katalog/models/user"
)

var (
	testDb     *BoltStore
	testDbPath = "testdb"
)

func createTestDb() {
	var err error
	testDb, err = New(testDbPath)
	if err != nil {
		panic(err)
	}
}

func cleanTestDb() {
	os.RemoveAll(testDbPath)
}

func TestUserSave(t *testing.T) {
	createTestDb()
	defer cleanTestDb()

	u := user.New()
	var vkID int64 = 1
	u.VkId = vkID

	Convey("Given some integer with a starting value", t, func() {

		err := testDb.UserSave(u)
		Convey("save user without error", func() {
			So(err, ShouldBeNil)
		})

		u2, err := testDb.UsersGetByVkID(vkID)
		Convey("Err nust be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("user must exists", func() {
			So(u2, ShouldNotEqual, nil)
		})

		Convey("user ids should be equal", func() {
			So(u2.Id, ShouldResemble, u.Id)
		})

		Convey("user vk_ids should be equal", func() {
			So(u2.VkId, ShouldResemble, u.VkId)
		})

		u3 := user.New()
		u3.VkId = vkID
		err = testDb.UserSave(u3)
		Convey("Create existing user", func() {
			So(err, ShouldNotBeNil)
		})

	})
}
