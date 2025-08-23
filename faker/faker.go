package faker

import (
	"fmt"
	"math/rand"
)

var Password = "$2a$04$MhoWw5urR0zStVXEry9JTOhDAJyyminPGaM637MRhRiWOrGxeopOu" // hashed of "password"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandStringLowerRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}
	return string(b)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func UserName() string {
	return RandStringLowerRunes(RandInt(2, 10))
}

func ID() string {
	return fmt.Sprintf("%s-%s-%s-%s", RandStringRunes(4), RandStringRunes(4), RandStringRunes(4), RandStringRunes(4))
}

func Email() string {
	return fmt.Sprintf("%s@%s.com", RandStringLowerRunes(RandInt(5, 10)), RandStringLowerRunes(3))
}
