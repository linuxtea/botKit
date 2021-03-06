// outside
package session

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	ConstTimeFormat = "20060102150405"
	ConstFromA      = "A"
	ConstFromB      = "B"
	ConstFromC      = "C"
)

type Session struct {
	From      string `json:"from"` // A B C
	SrcID     int    `json:"src_id"`
	ManagerID int    `json:"manager_id"`
	UserID    int    `json:"user_id"`
	Expire    string `json:"expire"`
	Signature string `json:"signature"`
}

func GenerateSession(from string, srcID, managerID, userID int, duration time.Duration,
	getSecretKey func(srcID int, managerID int) (secretKey string, err error)) (string, error) {
	expire := time.Now().Add(duration)
	secretKey, err := getSecretKey(srcID, managerID)
	if err != nil {
		return "", fmt.Errorf("GenerateSession from:%s srcID:%d managerID:%d userID:%d getSecretKey error %v",
			from, srcID, managerID, userID, err)
	}
	return newSession(from, srcID, managerID, userID, expire, secretKey)
}

func VerifySession(from string, session string,
	getSecretKey func(srcID int, managerID int) (secretKey string, err error)) (*Session, error) {
	jsonSource, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		return nil, fmt.Errorf("VerifySession base64 decode error %v", err)
	}
	s := &Session{}
	err = json.Unmarshal(jsonSource, s)
	if err != nil {
		return nil, fmt.Errorf("VerifySession json unmarshal error %v", err)
	}
	if from != s.From {
		return nil, fmt.Errorf("VerifySession diff from %s:%s", from, s.From)
	}
	expireTime, err := time.Parse(ConstTimeFormat, s.Expire)
	if err != nil {
		return nil, fmt.Errorf("VerifySession parse expire %v error %v", s.Expire, err)
	}
	if expireTime.Before(time.Now()) {
		return nil, fmt.Errorf("VerifySession session is timeout")
	}
	secretKey, err := getSecretKey(s.SrcID, s.ManagerID)
	if err != nil {
		return nil, fmt.Errorf("VerifySession srcID:%d managerID:%d userID:%d getSecretKey error %v",
			s.SrcID, s.ManagerID, s.UserID, err)
	}
	newSignature, err := newSignature(s.From, s.SrcID, s.ManagerID, s.UserID, expireTime, secretKey)
	if err != nil {
		return nil, fmt.Errorf("VerifySession srcID:%d managerID:%d userID:%d NewSignature error %v",
			s.SrcID, s.ManagerID, s.UserID, err)
	}
	if newSignature != s.Signature {
		return nil, fmt.Errorf("VerifySession unauthorized")
	}
	return s, nil
}

func newSession(from string, srcID, managerID, userID int, expire time.Time, secretKey string) (string, error) {
	signature, err := newSignature(from, srcID, managerID, userID, expire, secretKey)
	if err != nil {
		return "", fmt.Errorf("newSignature error %v", err)
	}
	data, err := json.Marshal(&Session{
		From:      from,
		SrcID:     srcID,
		ManagerID: managerID,
		UserID:    userID,
		Expire:    expire.Format(ConstTimeFormat),
		Signature: signature,
	})
	if err != nil {
		return "", fmt.Errorf("newSession json marshal error %v", err)
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func newSignature(from string, srcID, managerID, userID int, expire time.Time, secretKey string) (string, error) {
	switch from {
	case ConstFromA, ConstFromB, ConstFromC:
	default:
		return "", fmt.Errorf("newSignature unknown from %s", from)
	}

	source := from + strconv.Itoa(srcID) + strconv.Itoa(managerID) + strconv.Itoa(userID) + expire.Format(ConstTimeFormat) + secretKey
	sha1er := sha1.New()
	io.WriteString(sha1er, source)
	return fmt.Sprintf("%x", sha1er.Sum(nil)), nil
}
