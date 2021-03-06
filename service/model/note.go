package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"


	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2/bson"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/database"
)

type Note struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"`
	Content   string        `db:"content" bson:"content"`
	UserID    bson.ObjectId `bson:"user_id"`
	UID       uint32        `db:"user_id" bson:"userid,omitempty"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}

func (u *Note) NoteID() string {
	r := ""

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		r = fmt.Sprintf("%v", u.ID)
	case database.TypeMongoDB:
		r = u.ObjectID.Hex()
	case database.TypeBolt:
		r = u.ObjectID.Hex()
	}

	return r
}

func NoteByID(userID string, noteID string) (Note, error) {
	var err error

	result := Note{}

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Get(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM note WHERE id = ? AND user_id = ? LIMIT 1", noteID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

			if bson.IsObjectIdHex(noteID) {
				err = c.FindId(bson.ObjectIdHex(noteID)).One(&result)
				if result.UserID != bson.ObjectIdHex(userID) {
					result = Note{}
					err = ErrUnauthorized
				}
			} else {
				err = ErrNoResult
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		err = database.View("note", userID+noteID, &result)
		if err != nil {
			err = ErrNoResult
		}
		if result.UserID != bson.ObjectIdHex(userID) {
			result = Note{}
			err = ErrUnauthorized
		}
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

func NotesByUserID(userID string) ([]Note, error) {
	var err error

	var result []Note

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		err = database.SQL.Select(&result, "SELECT id, content, user_id, created_at, updated_at, deleted FROM note WHERE user_id = ?", userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

			if bson.IsObjectIdHex(userID) {
				err = c.Find(bson.M{"user_id": bson.ObjectIdHex(userID)}).All(&result)
			} else {
				err = ErrNoResult
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		err = database.BoltDB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("note"))
			if b == nil {
				return bolt.ErrBucketNotFound
			}

			c := b.Cursor()

			prefix := []byte(userID)
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				var single Note

				err := json.Unmarshal(v, &single)
				if err != nil {
					log.Println(err)
					continue
				}

				result = append(result, single)
			}

			return nil
		})
	default:
		err = ErrCode
	}

	return result, standardizeError(err)
}

func NoteCreate(content string, userID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO note (content, user_id) VALUES (?,?)", content, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

			note := &Note{
				ObjectID:  bson.NewObjectId(),
				Content:   content,
				UserID:    bson.ObjectIdHex(userID),
				CreatedAt: now,
				UpdatedAt: now,
				Deleted:   0,
			}
			err = c.Insert(note)
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		note := &Note{
			ObjectID:  bson.NewObjectId(),
			Content:   content,
			UserID:    bson.ObjectIdHex(userID),
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}

		err = database.Update("note", userID+note.ObjectID.Hex(), &note)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

func NoteUpdate(content string, userID string, noteID string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("UPDATE note SET content=? WHERE id = ? AND user_id = ? LIMIT 1", content, noteID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("note")
			var note Note
			note, err = NoteByID(userID, noteID)
			if err == nil {
				if note.UserID.Hex() == userID {
					note.UpdatedAt = now
					note.Content = content
					err = c.UpdateId(bson.ObjectIdHex(noteID), &note)
				} else {
					err = ErrUnauthorized
				}
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		var note Note
		note, err = NoteByID(userID, noteID)
		if err == nil {
			if note.UserID.Hex() == userID {
				note.UpdatedAt = now
				note.Content = content
				err = database.Update("note", userID+note.ObjectID.Hex(), &note)
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}

func NoteDelete(userID string, noteID string) error {
	var err error

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("DELETE FROM note WHERE id = ? AND user_id = ?", noteID, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("note")

			var note Note
			note, err = NoteByID(userID, noteID)
			if err == nil {
				if note.UserID.Hex() == userID {
					err = c.RemoveId(bson.ObjectIdHex(noteID))
				} else {
					err = ErrUnauthorized
				}
			}
		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:
		var note Note
		note, err = NoteByID(userID, noteID)
		if err == nil {
			if note.UserID.Hex() == userID {
				err = database.Delete("note", userID+note.ObjectID.Hex())
			} else {
				err = ErrUnauthorized
			}
		}
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
