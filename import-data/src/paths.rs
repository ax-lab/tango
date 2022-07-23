use std::path::Path;

#[derive(Clone, Copy, Debug)]
pub enum InputPath {
	JMDict,
}

impl AsRef<Path> for InputPath {
	fn as_ref(&self) -> &Path {
		Path::new(match self {
			InputPath::JMDict => "entries/JMdict.gz",
		})
	}
}

#[cfg(test)]
#[allow(non_snake_case)]
mod tests_InputPath {
	use super::*;
	use crate::get_path_in_vendor_data;

	#[test]
	fn converts_to_existing_path() {
		check_path(InputPath::JMDict);
	}

	fn check_path(input: InputPath) {
		let path = get_path_in_vendor_data(input);
		let msg = format!("checking that {:?}", input);
		assert!(path.is_some(), "{} exists", msg);
		assert!(path.unwrap().is_file(), "{} is a file", msg);
	}
}
