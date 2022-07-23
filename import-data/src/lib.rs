use std::path::{Path, PathBuf};

use once_cell::sync::Lazy;

mod paths;
pub use paths::*;

pub mod jmdict;

/// Finds the absolute path for a file or subfolder inside the vendor data
/// folder.
pub fn get_path_in_vendor_data<P: AsRef<Path>>(subpath: P) -> Option<PathBuf> {
	static VENDOR_DATA: Lazy<PathBuf> =
		Lazy::new(|| ["vendor", "data"].into_iter().collect::<PathBuf>());

	let subpath = subpath.as_ref();
	if subpath.has_root() {
		return None;
	}

	let mut search_path = VENDOR_DATA.clone();
	search_path.push(subpath);

	let mut basepath = std::env::current_exe().expect("failed to get current executable path");
	basepath.pop();
	loop {
		let vendor_path = basepath.join(&search_path);
		if let Ok(stat) = std::fs::metadata(&vendor_path) {
			if stat.is_dir() || stat.is_file() {
				return Some(vendor_path);
			}
		}

		if !basepath.pop() {
			return None;
		}
	}
}

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	fn returns_none_for_non_existing_path() {
		let path = get_path_in_vendor_data(
			"for/sure/this/cant/exist/cfdddac6-59e0-4fd5-834f-fec95a26c378.not",
		);
		assert!(path.is_none());
	}

	#[test]
	fn returns_none_for_rooted_path() {
		let path = get_path_in_vendor_data("/");
		assert!(path.is_none());
	}

	#[test]
	fn finds_vendor_folder_for_file() {
		let path = get_path_in_vendor_data("README.md");
		let path = check_vendor_path(path);

		assert!(path.is_file());
		assert_eq!(path.file_name().unwrap(), "README.md");
	}

	#[test]
	fn finds_vendor_folder_for_file_in_a_folder() {
		let path = get_path_in_vendor_data("entries/JMdict.gz");
		let path = check_vendor_path(path);
		assert!(path.is_file());
		assert_eq!(path.file_name().unwrap(), "JMdict.gz");
		assert_eq!(path.parent().unwrap().file_name().unwrap(), "entries");
	}

	#[test]
	fn finds_vendor_folder_for_subfolder() {
		let path = get_path_in_vendor_data("entries");
		let path = check_vendor_path(path);
		assert!(path.is_dir());
		assert_eq!(path.file_name().unwrap(), "entries");

		let jmdict = path.join("JMdict.gz");
		assert!(jmdict.is_file());
	}

	fn check_vendor_path(path: Option<PathBuf>) -> PathBuf {
		assert!(path.is_some());
		let path = path.unwrap();
		assert!(path.is_absolute());
		assert!(path.exists());
		return path;
	}
}
