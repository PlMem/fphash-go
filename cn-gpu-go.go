package cngpugo

// #cgo CXXFLAGS: -std=c++11 -mavx2 -msse2 -maes
// #cgo LDFLAGS: -static-libstdc++ -static-libgcc
// #include "cn-gpu.hxx"
import "C"

import (
	"sync"
	"unsafe"
)

type Ctx struct {
	cn_ctx unsafe.Pointer
}

// Ctx pool
var ctxPool = sync.Pool{
	New: func() any {
		return &Ctx{cn_ctx: C.new_ctx()}
	},
}

// Singleton ctx
var ctx *Ctx
var once sync.Once

func getCtx() *Ctx {
	once.Do(func() {
		ctx = &Ctx{cn_ctx: C.new_ctx()}
	})
	return ctx
}

func Variant_version() int {
	return int(C.variant_version(getCtx().cn_ctx))
}

func Hash(input []uint8) [32]uint8 {
	out := [32]uint8{}
	var in_ptr unsafe.Pointer
	if len(input) != 0 {
		in_ptr = unsafe.Pointer(&input[0])
	}
	ctx := ctxPool.Get().(*Ctx)
	C.cn_hash(
		ctx.cn_ctx,
		in_ptr,
		C.size_t(len(input)),
		unsafe.Pointer(&out[0]),
	)
	ctxPool.Put(ctx)
	return out
}
