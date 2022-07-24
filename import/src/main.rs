use std::time::Instant;

use import_data::*;

fn main() {
	println!("\nImporting JMdict data...");

	let now = Instant::now();
	let mut dict = jmdict::Input::open(InputPath::JMDict).unwrap();
	let entries = dict.entries().unwrap();

	let elapsed = now.elapsed();
	println!("... opened file in {:?}", elapsed);

	let now = Instant::now();
	let count = entries.count();
	let elapsed = now.elapsed();
	println!("... read {} entries in {:?}", count, elapsed);
}
