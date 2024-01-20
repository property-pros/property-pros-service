package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/vireocloud/property-pros-service/constants"
)

func GetUserIdFromContext(ctx context.Context) (string, error) {
	userIdFromContext := ctx.Value(constants.UserIdKey)

	if userIdFromContext == nil {
		return "", errors.New("unresolved userid")
	}

	usrID := fmt.Sprintf("%v", userIdFromContext)

	return usrID, nil
}
