package boltstore

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/smolgu/katalog/models/user"
	"github.com/smolgu/katalog/types"
	"github.com/zhuharev/boltutils"
)

var (
	// TrainingDirectionBucketName reprezent name of boltdb bucket
	TrainingDirectionBucketName = []byte("td")

	// ErrUniqueConstraintViolated error if insert value with same unique field
	// and different primary keys
	ErrUniqueConstraintViolated = fmt.Errorf("%s", "unique constraint violated")
)

var (
	_ types.Store = &BoltStore{}
)

// BoltStore wrap bolt.DB
type BoltStore struct {
	db *boltutils.DB
}

func makeCreateBucketsFunc() func(tx *bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(UsersBucketName)
		if err != nil {
			return err
		}
		return nil
	}
}

// New returns BoltStore
func New(path string) (*BoltStore, error) {
	db, err := bolt.Open(path, 0777, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(makeCreateBucketsFunc())
	if err != nil {
		return nil, err
	}
	return &BoltStore{
		db: boltutils.New(db),
	}, nil
}

// GetTrainingDirectionByShortName returns td by his shortName
func (bs *BoltStore) GetTrainingDirectionByShortName(shortName string) (*user.TrainingDirection, error) {
	td := new(user.TrainingDirection)
	err := bs.db.Iterate(TrainingDirectionBucketName, func(k, v []byte) error {
		err := proto.Unmarshal(v, td)
		if err != nil {
			return err
		}
		if td.ShortName == shortName {
			return boltutils.ErrBreak
		}
		return nil
	})
	if err != nil && err != boltutils.ErrBreak {
		return nil, err
	} else if err == boltutils.ErrBreak {
		return td, nil
	}
	return nil, boltutils.ErrNotFound
}

func (bs *BoltStore) hasTrainingDirectionWithShortName(shortName string) (bool, error) {
	err := bs.db.Iterate(TrainingDirectionBucketName, func(k, v []byte) error {
		td := new(user.TrainingDirection)
		err := proto.Unmarshal(v, td)
		if err != nil {
			return err
		}
		if td.ShortName == shortName {
			return boltutils.ErrBreak
		}
		return boltutils.ErrNotFound
	})
	if err == boltutils.ErrBreak {
		return true, nil
	}
	return false, err
}

// SetTrainingDirection set value into db
func (bs *BoltStore) SetTrainingDirection(td *user.TrainingDirection) error {
	// check td exists
	oldTd, err := bs.GetTrainingDirectionByShortName(td.ShortName)
	if err != nil && err != boltutils.ErrNotFound {
		return err
	}
	if err == nil {
		if !bytes.Equal(oldTd.Id, td.Id) {
			return ErrUniqueConstraintViolated
		}
	}
	bts, err := proto.Marshal(td)
	if err != nil {
		return err
	}
	return bs.db.Put(TrainingDirectionBucketName, td.Id, bts)
}
