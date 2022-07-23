use std::{
	sync::mpsc::{self, RecvTimeoutError},
	thread,
	time::Duration,
};

/// Runs the given closure in a separate thread. This will panic if the closure
/// does not return within the given timeout.
pub fn panic_after<T, F>(timeout: Duration, f: F) -> T
where
	T: Send + 'static,
	F: FnOnce() -> T,
	F: Send + 'static,
{
	let (tx, rx) = mpsc::channel();
	let handle = thread::Builder::new()
		.name("test callback".into())
		.spawn(move || {
			let val = f();
			tx.send(val).expect("failed to send test results");
		})
		.unwrap();

	match rx.recv_timeout(timeout) {
		Ok(val) => {
			handle.join().expect("test callback join");
			val
		}
		Err(RecvTimeoutError::Disconnected) => {
			handle.join().expect("test callback panicked");
			panic!("test callback disconnected");
		}
		Err(RecvTimeoutError::Timeout) => panic!("test callback timed out after {:?}", timeout),
	}
}

#[cfg(test)]
mod tests {
	use std::thread::sleep;

	use tux::assert_panic;

	use super::*;

	#[test]
	fn returns_the_callback_result() {
		let output = panic_after(Duration::from_millis(10), || "abc");
		assert_eq!(output, "abc");

		let output = panic_after(Duration::from_millis(10), || 123);
		assert_eq!(output, 123);
	}

	#[test]
	fn panics_on_timeout() {
		let do_timeout = || {
			panic_after(Duration::from_millis(1), || {
				sleep(Duration::from_millis(50))
			})
		};
		assert_panic!("timed out after" in do_timeout());
	}

	#[test]
	fn panics_on_callback_panic() {
		let do_panic = || {
			panic_after(Duration::from_millis(50), || {
				panic!("xyz");
			})
		};
		assert_panic!("test callback panicked" in do_panic());
	}
}
