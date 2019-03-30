package model

import (
	"crypto/md5"
	"fmt"
	"gochat/config"
	"math/rand"
	"strconv"
	"time"
)

// 使用token验证，即无需上线后再验证，自动获取用户的Id等信息

type AuthInfo struct {
	Id       string
	Username string
	Roomname string
	Token    string
	Expires  time.Time
}

func NewInfo(id, username, roomname string) *AuthInfo {
	token := NewToken()
	fmt.Println(token)
	return &AuthInfo{
		Id:       id,
		Username: username,
		Roomname: roomname,
		Token:    token,
		Expires:  time.Now().Add(30 * time.Minute),
	}
}

func NewId(username, roomname string) string {
	nowTime := strconv.Itoa(int(time.Now().Unix()))
	ctx := md5.New()
	ctx.Write([]byte(nowTime + username + roomname))
	return fmt.Sprintf("%x", ctx.Sum(nil))
}

func NewToken() string {
	now := time.Now().UnixNano()
	rand.Seed(now)
	var letterRunes = []rune(config.TokenSecureCode)
	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	fmt.Println(string(b))
	nowTime := strconv.Itoa(int(now))
	ctx := md5.New()
	ctx.Write([]byte(nowTime + string(b)))
	return fmt.Sprintf("%x", ctx.Sum(nil))
}
