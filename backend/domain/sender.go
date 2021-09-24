package domain

import (
	"math/rand"
	"strings"
)

type Sender struct {
	Id     string `json:"id"`
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

func NewSender(id string) *Sender {
	return &Sender{
		Id:     id,
		Avatar: getRandomEmoji(),
		Name:   getRandomWord() + " " + getRandomWord(),
	}
}
