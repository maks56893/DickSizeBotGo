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

var listAverage = [...]string{
	"Усреднённые младшие ваши ",
	"Усреднённые члены ",
	"Усреднённые стручки ",
	"Усреднённые волшебные палочки ",
	"Усреднённые козыри в рукаве ",
	"Усреднённые чупачупсы ",
	"Усреднённые песисы ",
	"Усреднённые болты ",
	"Усреднённые хоботки ",
	"Усреднённые лысые Джонни Синсы ",
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
