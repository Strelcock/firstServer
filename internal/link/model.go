package link

import (
	"math/rand/v2"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqueIndex"`
}

func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: randStringRunes(6),
	}

}

var runesSlice = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	res := make([]rune, n)
	for i := range res {
		res[i] = runesSlice[rand.IntN(len(runesSlice))]
	}
	return string(res)
}

func (l *Link) RegenerateHash() {
	hashLen := len([]rune(l.Hash))
	l.Hash = randStringRunes(hashLen)
}
