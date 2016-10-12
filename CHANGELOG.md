# Release 0.3.0 (2016-10-12)

 - sync: Limited condvar support on darwin and windows via spin waiters.
 - mq: FastMq blocking mode for darwin, windows.
 - sync: Default mutex implementation on freebsd uses futex via umtx syscall.

# Release 0.2.0 (2016-06-22)

- Added FastMq implementation. FastMq is a priority queue based on shared memory. See mq package docs for details.
- Huge windows native shm refactor. It has new API now.
- Improved errors reporting: use github.com/pkg/errors to wrap errors.

# Release 0.1.0 (2016-05-24)

- Initial release.