package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-9")
	ctx = WithUserID(ctx, "alice")
	ctx = WithTenantID(ctx, "acme")

	rid, _ := RequestID(ctx)
	uid, _ := UserID(ctx)
	tid, _ := TenantID(ctx)
	fmt.Printf("request=%s user=%s tenant=%s\n", rid, uid, tid)
}
