package main

import (
	"PlantSite/internal/infra/admins"
	"PlantSite/internal/models/auth"
	authservice "PlantSite/internal/services/auth-service"
	"context"
	"fmt"

	"github.com/spf13/viper"
)

const (
	AdminsList       = "admins"
	AdminLoginKey    = "login"
	AdminPasswordKey = "password"
)

func GetAdminsMap(hasher authservice.PasswdHasher) *admins.AdminMap {
	ctx := context.Background()
	if err := ReadInConfig(); err != nil {
		panic(err)
	}
	m := admins.NewAdminMap()

	var admins []map[string]string
	if err := viper.UnmarshalKey(AdminsList, &admins); err != nil {
		panic(fmt.Errorf("error unmarshalling admins: %w", err))
	}

	for _, admin := range admins {
		login := admin[AdminLoginKey]
		password := admin[AdminPasswordKey]
		passwdHash, err := hasher.Hash([]byte(password))
		if err != nil {
			panic(fmt.Errorf("error hashing admin password: %w", err))
		}
		adm, err := auth.NewAdmin(login, passwdHash)
		if err != nil {
			panic(fmt.Errorf("error creating admin: %w", err))
		}
		if err := m.Add(ctx, *adm); err != nil {
			panic(fmt.Errorf("error adding admin: %w", err))
		}
	}

	return m
}
