package main

import (
	"math/rand"
	"strconv"
	"time"
)

var listMeasure = [...]string{
	"Мясная сигара у тебя ",
	"Младший у тебя ",
	"Питон в кустах у тебя ",
	"Чупачупс у тебя ",
	"Нагибатель у тебя ",
	"Бур у тебя ",
	"Лоллипап у тебя ",
	"Гуч у тебя ",
	"Младший у тебя ",
	"Волшебная палочка у тебя ",
	"Пенис у тебя ",
	"Лысый Джонни Синс у тебя ",
	"Писюлька у тебя ",
	"Писюн у тебя ",
	"Пиструн у тебя ",
	"Талант у тебя ",
	"Болт у тебя ",
	"Стручок у тебя ",
	"Член у тебя ",
	"Хоботок у тебя ",
	"Писюндель у тебя ",
	"Грибочек у тебя ",
	"21 палец у тебя ",
	"Пипка у тебя ",
	"Малыш у тебя ",
	"Младший дружок у тебя ",
}

var listOfSmallDickEmoji = [...]string{
	"😔",
	"😒",
	"😢",
	"🥲",
	"😞",
	"😔",
	"😒",
	"😣",
	"😖",
	"😨",
	"😒",
	"🌚",
}

var listOfBigDickEmoji = [...]string{
	"😎",
	"🤭",
	"😮",
	"😱",
	"😯",
	"😏",
	"😁",
	"🌝",
	"🍆",
	//	"🚀",
}

var listOfTodayDicks = [...]string{
	"члены",
	"морковки",
	"удавы",
	"питоны",
	"чупачупсы",
	"стручки",
	"волшебные палочки",
	"пенисы",
	"шишки",
	"лысые Джонни Синсы",
}

var listAverage = [...]string{
	"<i>Усреднённые младшие ваши</i>\n\n ",
	"<i>Усреднённые члены</i>\n\n ",
	"<i>Усреднённые стручки</i>\n\n ",
	"<i>Усреднённые волшебные палочки</i>\n\n ",
	"<i>Усреднённые козыри в рукаве</i>\n\n ",
	"<i>Усреднённые чупачупсы</i>\n\n ",
	"<i>Усреднённые песисы</i>\n\n ",
	"<i>Усреднённые болты</i>\n\n ",
	"<i>Усреднённые хоботки</i>\n\n ",
	"<i>Усреднённые лысые Джонни Синсы</i>\n\n ",
}

func GetRandMeasureReplyPattern(dickSize int) string {
	rand.Seed(time.Now().UnixNano())
	var res string
	if dickSize < 20 {
		res = listMeasure[rand.Intn(len(listMeasure))] + strconv.Itoa(dickSize) + "см" + listOfSmallDickEmoji[rand.Intn(len(listOfSmallDickEmoji))]
	} else {
		res = listMeasure[rand.Intn(len(listMeasure))] + strconv.Itoa(dickSize) + "см" + listOfBigDickEmoji[rand.Intn(len(listOfBigDickEmoji))]
	}

	return res
}

func GetRandAverageReplyText() string {
	rand.Seed(time.Now().UnixNano())
	return listAverage[rand.Intn(len(listAverage))]
}

func GetRandTodayReplyText() string {
	rand.Seed(time.Now().UnixNano())
	return "<b>Свежие (и не очень ❗️) " + listOfTodayDicks[rand.Intn(len(listOfTodayDicks))] + "\n\n</b>"
}
