package grpcutil

import (
	"context"

	commonv1 "server/api/common/v1"
	"server/internal/iam/domain"
	"server/pkg/logger"
	reshelper "server/pkg/util/response_helper"
)

func HandleBusinessError(ctx context.Context, handlerName string, req any, err error) *commonv1.BaseResponse {
	businessError, _ := domain.GetBusinessError(err)
	logger.LogBusinessError(ctx, *businessError, map[string]interface{}{
		"handler": handlerName,
		"request": req,
	})
	return reshelper.BuildErrorResponse(ctx, businessError)
}
