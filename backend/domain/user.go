package domain

import (
	"math/rand"
	"strings"
)

type User struct {
	Id     int32  `json:"id"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
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
		Id:     int32(rand.Int31()),
		Avatar: getRandomEmoji(),
		Name:   getRandomWord() + " " + getRandomWord(),
	}
}
