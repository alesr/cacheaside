package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/alesr/cacheaside/internal/repository/cache/memcache"
	"github.com/alesr/cacheaside/internal/repository/memdb"
	"github.com/alesr/cacheaside/internal/service"
)

func main() {
	logger := slog.Default()

	svc := service.New(
		memcache.New(logger, memdb.New()),
	)

	item1 := svc.Create()
	item2 := svc.Create()
	_ = svc.Create()
	_ = svc.Create()
	_ = svc.Create()

	items := svc.List()

	for _, item := range items {
		fmt.Println(item.ID)
	}

	item, err := svc.Fetch(item1.ID)
	if err != nil {
		logger.Error("could not fetch item", slog.String("error", err.Error()))
		os.Exit(1)
	}

	fmt.Println(item.ID)

	item, err = svc.Fetch(item2.ID)
	if err != nil {
		logger.Error("could not fetch item", slog.String("error", err.Error()))
		os.Exit(1)
	}

	fmt.Println(item.ID)
}
