#[macro_export]
macro_rules! defer {
	($e: expr) => {
		let _scope = $crate::defer::DeferGuard::new(|| -> () {
			let _ = $e;
		});
	};
}

pub struct DeferGuard<F: FnMut()> {
	callback: F,
}

impl<F: FnMut()> DeferGuard<F> {
	pub fn new(callback: F) -> Self {
		DeferGuard { callback }
	}
}

impl<F: FnMut()> Drop for DeferGuard<F> {
	fn drop(&mut self) {
		(self.callback)();
	}
}

#[cfg(test)]
mod tests {
	use std::{
		cell::RefCell,
		rc::Rc,
		sync::{Arc, Mutex},
	};

	use tux::assert_panic;

	#[test]
	pub fn executes_the_expression_once() {
		let counter = Rc::new(RefCell::new(0));
		let dropper = || {
			let my_count = counter.clone();
			defer!(*my_count.borrow_mut() += 1);
		};

		dropper();
		assert_eq!(*counter.borrow(), 1);
	}

	#[test]
	pub fn executes_all_defers() {
		let counter = Rc::new(RefCell::new(0));
		let dropper = || {
			let my_count = counter.clone();
			defer!(*my_count.borrow_mut() += 1);
			defer!(*my_count.borrow_mut() += 2);
		};

		dropper();
		assert_eq!(*counter.borrow(), 3);
	}

	#[test]
	pub fn executes_on_panic() {
		let counter = Arc::new(Mutex::new(0));
		let dropper = || {
			let my_count = counter.clone();
			defer!(*my_count.lock().unwrap() += 1);
			defer!(*my_count.lock().unwrap() += 2);
			defer!(*my_count.lock().unwrap() += 4);
			panic!("some panic");
		};

		assert_panic!("some panic" in dropper());
		assert_eq!(*counter.lock().unwrap(), 7);
	}
}
