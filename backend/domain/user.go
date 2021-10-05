package domain

import (
	"math/rand"
	"strings"
	"time"
)

type Notification struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type User struct {
	Id     int64             `json:"id"`
	Avatar string            `json:"avatar"`
	Name   string            `json:"name"`
	Send   chan Notification `json:"-"`
}

func getRandomEmoji() string {
	emojis := []string{
		"🧑",
		"💇‍♀️",
		"💇‍♂️",
		"💇",
		"💆‍♀️",
		"💆‍♂️",
		"💆",
		"🧟‍♀️",
		"🧟‍♂️",
		"🧟",
		"🧝‍♀️",
		"🧝‍♂️",
		"🧝",
		"🧛‍♀️",
		"🧛‍♂️",
		"🦹‍♀️",
		"🦹‍♂️",
		"🦸‍♀️",
		"👱",
		"👨",
		"🧔",
		"👨‍🦰",
		"👨‍🦱",
		"👨‍🦳",
		"👨‍🦲",
		"👩",
		"👩‍🦰",
		"🧑‍🦰",
		"👩‍🦱",
		"🧑‍🦱",
		"👩‍🦳",
		"🧑‍🦳",
		"👩‍🦲",
		"🧑‍🦲",
		"👱‍♀️",
		"👱‍♂️",
		"🧓",
		"👴",
		"👵",
		"🙍",
		"🙍‍♂️",
		"🙍‍♀️",
		"🙎",
		"🙎‍♂️",
		"🙎‍♀️",
		"🙅",
		"🙅‍♂️",
		"🙅‍♀️",
		"🙆",
		"🙆‍♂️",
		"🙆‍♀️",
		"💁",
		"💁‍♂️",
		"💁‍♀️",
		"🙋",
		"🙋‍♂️",
		"🙋‍♀️",
		"🧏‍♂️",
		"🧏‍♀️",
		"🙇‍♂️",
		"🙇‍♀️",
		"🤦‍♂️",
		"🤦‍♀️",
		"🤷‍♂️",
		"🤷‍♀️",
		"👨‍⚕️",
		"👩‍⚕️",
		"👨‍🎓",
		"👩‍🎓",
		"👨‍🏫",
		"👩‍🏫",
		"👨‍⚖️",
		"👩‍⚖️",
		"👨‍🌾",
		"👩‍🌾",
		"👨‍🍳",
		"👩‍🍳",
		"👨‍🔧",
		"👩‍🔧",
		"👨‍🏭",
		"👩‍🏭",
		"👨‍💼",
		"👩‍💼",
		"👨‍🔬",
		"👩‍🔬",
		"👨‍💻",
		"👩‍💻",
		"👨‍🎤",
		"👩‍🎤",
		"👨‍🎨",
		"👩‍🎨",
		"👨‍✈️",
		"👩‍✈️",
		"👨‍🚀",
		"👩‍🚀",
		"👨‍🚒",
		"👩‍🚒",
		"👮‍♂️",
		"👮‍♀️",
		"🕵️‍♂️",
		"🕵️‍♀️",
		"💂‍♂️",
		"💂‍♀️",
		"🥷",
		"👷‍♂️",
		"👷‍♀️",
		"🤴",
		"👸",
		"👳‍♂️",
		"👳‍♀️",
		"🧕",
		"🤵‍♂️",
		"🤵‍♀️",
		"👰",
		"👰‍♂️",
		"👰‍♀️",
		"🤰",
		"🤱",
		"👩‍🍼",
		"👨‍🍼",
		"🧑‍🍼",
		"👼",
		"🎅",
		"🤶",
		"🧑‍🎄",
		"🦸",
		"🦸‍♂️",
		"🦸‍♀️",
	}
	randomIndex := rand.Intn(len(emojis))
	return emojis[randomIndex]
}

func getRandomWord() string {
	randomLetters := ""
	for i := 0; i < 6; i++ {
		randomLetters += string(rune(rand.Intn(26) + 97))

		if i == 0 {
			randomLetters = strings.ToUpper(randomLetters)
		}

	}

	return randomLetters
}

func NewUser() *User {
	return &User{
		Id:     int64(time.Now().UnixNano()),
		Avatar: getRandomEmoji(),
		Name:   getRandomWord() + " " + getRandomWord(),
		Send:   make(chan Notification, 256),
	}
}

func (c *User) NewNotification(notificationType string, data interface{}) Notification {
	return Notification{
		Type: notificationType,
		Data: data,
	}
}
