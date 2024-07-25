package main

import (
	"fmt"
	"io"
	"net/http"
    "github.com/tidwall/gjson"
    "crypto/sha1"
)

type Gamejolt struct {
    gameid int
    privatekey string
}

func (gamejolt *Gamejolt) LoginPlayer(username string, gametoken string) bool {

    baseURL := fmt.Sprintf("https://api.gamejolt.com/api/game/v1_2/users/auth?game_id=%d&username=%s&user_token=%s", gamejolt.gameid, username, gametoken)

    signURL := fmt.Sprintf("%s%s", baseURL, gamejolt.privatekey)

    sha := sha1.New()
    sha.Write([]byte(signURL))

    requestURL := fmt.Sprintf("%s&signature=%x", baseURL, sha.Sum(nil))
    res, err := http.Get(requestURL)

    if err != nil {
        println(err)
        return false
    }

    body, err := io.ReadAll(res.Body)

    if err != nil {
        println(err)
        return false
    }

    answer := gjson.Get(string(body), "response.success").Bool()

    return answer
}

func (gamejolt *Gamejolt) AddTrophy(username string, gametoken string, trophyid int) (bool) {
    baseURL := fmt.Sprintf("https://api.gamejolt.com/api/game/v1_2/trophies/add-achieved/?game_id=%d&username=%s&user_token=%s&trophy_id=%d", gamejolt.gameid, username, gametoken, trophyid)

    signURL := fmt.Sprintf("%s%s", baseURL, gamejolt.privatekey)

    sha := sha1.New()
    sha.Write([]byte(signURL))

    requestURL := fmt.Sprintf("%s&signature=%x", baseURL, sha.Sum(nil))
    res, err := http.Get(requestURL)

    if err != nil {
        println(err)
        return false
    }

    body, err := io.ReadAll(res.Body)

    if err != nil {
        println(err)
        return false
    }

    answer := gjson.Get(string(body), "response.success").Bool()

    return answer
}
