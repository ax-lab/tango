use std::{
	fs::File,
	io::{BufRead, BufReader, Error, ErrorKind, Result, Seek, SeekFrom},
	path::Path,
};

use flate2::read::GzDecoder;
use once_cell::sync::Lazy;
use regex::Regex;

use crate::get_path_in_vendor_data;

pub struct Input {
	file: File,
}

impl Input {
	pub fn open<P: AsRef<Path>>(input: P) -> Result<Input> {
		let input = get_path_in_vendor_data(input);
		if let Some(input) = input {
			let file = File::open(input)?;

			let decoder = GzDecoder::new(file);
			decoder.header().ok_or(Error::new(
				ErrorKind::InvalidInput,
				"input must be a gzip file",
			))?;

			let mut file = decoder.into_inner();
			file.seek(SeekFrom::Start(0))?;
			Ok(Input { file })
		} else {
			Err(ErrorKind::NotFound.into())
		}
	}

	pub fn entries(&mut self) -> Result<EntryIterator> {
		let file = &mut self.file;
		file.seek(SeekFrom::Start(0))?;
		Ok(EntryIterator::new(file))
	}
}

pub struct EntryIterator<'a> {
	input: BufReader<GzDecoder<&'a File>>,
	buffer: String,
}

impl<'a> EntryIterator<'a> {
	fn new(input: &'a File) -> Self {
		let input = BufReader::new(GzDecoder::new(input));
		EntryIterator {
			input,
			buffer: String::new(),
		}
	}
}

impl<'a> Iterator for EntryIterator<'a> {
	type Item = Result<Entry>;

	fn next(&mut self) -> Option<Self::Item> {
		static RE_SEQ: Lazy<Regex> = Lazy::new(|| Regex::new(r"<ent_seq>(\d+)").unwrap());

		loop {
			self.buffer.clear();
			let count = match self.input.read_line(&mut self.buffer) {
				Err(err) => {
					return Some(Err(err));
				}
				Ok(count) => count,
			};

			if count == 0 {
				return None;
			}

			if let Some(m) = RE_SEQ.captures(&self.buffer) {
				let seq = m.get(1).unwrap().as_str().parse::<u64>().unwrap();
				return Some(Ok(Entry { seq }));
			}
		}
	}
}

pub struct Entry {
	pub seq: u64,
}

#[cfg(test)]
mod tests {
	use crate::*;
	use jmdict::*;

	#[test]
	fn opens_input() {
		Input::open(InputPath::JMDict).unwrap();
	}

	#[test]
	fn open_returns_error_for_non_existent_file() {
		let dict = Input::open("cfdddac6-59e0-4fd5-834f-fec95a26c378.not-a-file");
		assert!(dict.is_err());
	}

	#[test]
	fn reads_entries() -> Result<()> {
		let mut dict = helper::open();
		let mut entries = dict.entries()?;

		let entry = entries.next().unwrap()?;
		assert_eq!(entry.seq, 1000000);

		let entry = entries.next().unwrap()?;
		assert!(entry.seq > 1000000);

		Ok(())
	}

	mod helper {
		use super::*;

		pub fn open() -> Input {
			Input::open(InputPath::JMDict).unwrap()
		}
	}
}
