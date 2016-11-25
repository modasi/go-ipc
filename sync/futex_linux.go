// Copyright 2016 Aleksandr Demakin. All rights reserved.

package sync

import (
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"

	"bitbucket.org/avd/go-ipc/internal/allocator"
	"bitbucket.org/avd/go-ipc/internal/common"
)

const (
	cFUTEX_WAIT        = 0
	cFUTEX_WAKE        = 1
	cFUTEX_REQUEUE     = 3
	cFUTEX_CMP_REQUEUE = 4
	cFUTEX_WAKE_OP     = 5

	// FUTEX_PRIVATE_FLAG is used to optimize futex usage for process-private futexes.
	FUTEX_PRIVATE_FLAG = 128
	// FUTEX_CLOCK_REALTIME is used to tell the kernel, that is must treat timeouts for
	// FUTEX_WAIT_BITSET, FUTEX_WAIT_REQUEUE_PI, and FUTEX_WAIT as an absolute time based on CLOCK_REALTIME
	FUTEX_CLOCK_REALTIME = 256
)

// FutexWait checks if the the value equals futex's value.
// If it doesn't, Wait returns EWOULDBLOCK.
// Otherwise, it waits for the Wake call on the futex for not longer, than timeout.
func FutexWait(addr unsafe.Pointer, value uint32, timeout time.Duration, flags int32) error {
	fun := func(tm time.Duration) error {
		var ptr unsafe.Pointer
		if flags&FUTEX_CLOCK_REALTIME != 0 {
			ptr = unsafe.Pointer(common.AbsTimeoutToTimeSpec(tm))
		} else {
			ptr = unsafe.Pointer(common.TimeoutToTimeSpec(tm))
		}
		_, err := sys_futex(addr, cFUTEX_WAIT|flags, value, ptr, nil, 0)
		return err
	}
	return common.UninterruptedSyscallTimeout(fun, timeout)
}

// FutexWake wakes count threads waiting on the futex.
// Returns number of woken threads.
func FutexWake(addr unsafe.Pointer, count uint32, flags int32) (int, error) {
	for {
		woken, err := sys_futex(addr, cFUTEX_WAKE|flags, count, nil, nil, 0)
		if !common.IsInterruptedSyscallErr(err) {
			return int(woken), err
		}
	}
}

func sys_futex(addr unsafe.Pointer, op int32, val uint32, ts, addr2 unsafe.Pointer, val3 uint32) (int32, error) {
	r1, _, err := unix.Syscall6(unix.SYS_FUTEX,
		uintptr(addr),
		uintptr(op),
		uintptr(val),
		uintptr(ts),
		uintptr(addr2),
		uintptr(val3))
	allocator.Use(addr)
	allocator.Use(addr2)
	if err != 0 {
		return 0, os.NewSyscallError("FUTEX", err)
	}
	return int32(r1), nil
}
