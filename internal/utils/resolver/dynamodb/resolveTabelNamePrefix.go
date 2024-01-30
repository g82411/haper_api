package dynamodb

import (
	"context"
	"fmt"
	con "hyper_api/internal/config"
)

func ResolveTableNameWithPrefix(ctx context.Context, tableName string) string {
	config := ctx.Value("config").(*con.Config)

	return fmt.Sprintf("%s_%s", config.Env, tableName)
}
