package utils

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func GenUUID() string {
	uuidWithDate := fmt.Sprintf("%s-%v", time.Now().Local().Format("20060102"), uuid.New().String())
	return uuidWithDate
}
