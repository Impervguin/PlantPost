package main

import "github.com/spf13/viper"

const (
	HashPrefix  = "hasher"
	HashCostKey = "hash_cost"
)

func GetHashCost() int {
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetInt(Key(HashPrefix, HashCostKey))
}
