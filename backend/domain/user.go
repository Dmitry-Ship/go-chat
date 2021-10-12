package domain

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	Id     uuid.UUID `json:"id"`
	Avatar string    `json:"avatar"`
	Name   string    `json:"name"`
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
		Id:     uuid.New(),
		Avatar: getRandomEmoji(),
		Name:   getRandomWord() + " " + getRandomWord(),
	}
}
