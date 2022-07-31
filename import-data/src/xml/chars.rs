pub fn is_valid(c: char) -> bool {
	match c {
		'\u{0009}'
		| '\u{000A}'
		| '\u{000D}'
		| '\u{0020}'..='\u{D7FF}'
		| '\u{E000}'..='\u{FFFD}'
		| '\u{10000}'..='\u{10FFFF}' => true,
		_ => false,
	}
}

pub fn is_whitespace(c: char) -> bool {
	matches!(c, '\x20' | '\x09' | '\x0D' | '\x0A')
}

pub fn is_name_start(c: char) -> bool {
	match c {
		':'
		| 'A'..='Z'
		| '_'
		| 'a'..='z'
		| '\u{C0}'..='\u{D6}'
		| '\u{D8}'..='\u{F6}'
		| '\u{F8}'..='\u{2FF}'
		| '\u{370}'..='\u{37D}'
		| '\u{37F}'..='\u{1FFF}'
		| '\u{200C}'..='\u{200D}'
		| '\u{2070}'..='\u{218F}'
		| '\u{2C00}'..='\u{2FEF}'
		| '\u{3001}'..='\u{D7FF}'
		| '\u{F900}'..='\u{FDCF}'
		| '\u{FDF0}'..='\u{FFFD}'
		| '\u{10000}'..='\u{EFFFF}' => true,
		_ => false,
	}
}

pub fn is_name(c: char) -> bool {
	matches!(c, '-' | '.' | '0'..='9' | '\u{B7}' | '\u{0300}'..='\u{036F}' | '\u{203F}'..='\u{2040}')
		|| is_name_start(c)
}

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	fn valid() {
		helper::check_range(' ', '\u{00FF}', is_valid, "is valid");
		helper::check_char('\t', is_valid, "is valid");
		helper::check_char('\r', is_valid, "is valid");
		helper::check_char('\n', is_valid, "is valid");
		helper::check_char('\u{D7FF}', is_valid, "is valid");
		helper::check_char('\u{E000}', is_valid, "is valid");
		helper::check_char('\u{FFFD}', is_valid, "is valid");
		helper::check_char('\u{10000}', is_valid, "is valid");
		helper::check_char('\u{10FFFF}', is_valid, "is valid");

		let invalid = |c: char| !is_valid(c);

		helper::check_char('\u{000B}', invalid, "is not valid");
		helper::check_char('\u{000C}', invalid, "is not valid");
		helper::check_char('\u{000E}', invalid, "is not valid");
		helper::check_char('\u{000F}', invalid, "is not valid");

		helper::check_char('\u{FFFE}', invalid, "is not valid");
		helper::check_char('\u{FFFF}', invalid, "is not valid");

		helper::check_range('\u{0000}', '\u{0008}', invalid, "is not valid");
		helper::check_range('\u{0010}', '\u{001F}', invalid, "is not valid");

		let s1 = unsafe { char::from_u32_unchecked(0xD800) };
		let s2 = unsafe { char::from_u32_unchecked(0xDBFF) };
		helper::check_range(s1, s2, invalid, "is not valid");

		let s1 = unsafe { char::from_u32_unchecked(0xDC00) };
		let s2 = unsafe { char::from_u32_unchecked(0xDFFF) };
		helper::check_range(s1, s2, invalid, "is not valid");
	}

	#[test]
	fn whitespace() {
		helper::check_char(' ', is_whitespace, "is whitespace");
		helper::check_char('\t', is_whitespace, "is whitespace");
		helper::check_char('\r', is_whitespace, "is whitespace");
		helper::check_char('\n', is_whitespace, "is whitespace");
		helper::check_range('a', 'z', |c| !is_whitespace(c), "is not whitespace");
	}

	#[test]
	fn name_start() {
		helper::check_name_start(is_name_start, "is name start");
		helper::check_char(' ', |c| !is_name_start(c), "is not name start");
	}

	#[test]
	fn name() {
		helper::check_char('-', is_name, "is name");
		helper::check_char('-', is_name, "is name");
		helper::check_char('·', is_name, "is name");
		helper::check_range('0', '9', is_name, "is name");
		helper::check_range('\u{0300}', '\u{036F}', is_name, "is name");
		helper::check_range('\u{203F}', '\u{2040}', is_name, "is name");
		helper::check_name_start(is_name, "is name");
		helper::check_char(' ', |c| !is_name(c), "is not name");
	}

	mod helper {
		pub fn check_name_start<F: Fn(char) -> bool>(check: F, name: &'static str) {
			check_char(':', &check, name);
			check_char('_', &check, name);
			check_range('A', 'Z', &check, name);
			check_range('a', 'z', &check, name);
			check_range('\u{C0}', '\u{D6}', &check, name);
			check_range('\u{D8}', '\u{F6}', &check, name);
			check_range('\u{F8}', '\u{2FF}', &check, name);
			check_range('\u{370}', '\u{37D}', &check, name);
			check_range('\u{37F}', '\u{1FFF}', &check, name);
			check_range('\u{200C}', '\u{200D}', &check, name);
			check_range('\u{2070}', '\u{218F}', &check, name);
			check_range('\u{2C00}', '\u{2FEF}', &check, name);
			check_range('\u{3001}', '\u{D7FF}', &check, name);
			check_range('\u{F900}', '\u{FDCF}', &check, name);
			check_range('\u{FDF0}', '\u{FFFD}', &check, name);
			check_range('\u{10000}', '\u{10100}', &check, name);
			check_range('\u{EFF00}', '\u{EFFFF}', &check, name);
		}

		pub fn check_char<F: Fn(char) -> bool>(chr: char, check: F, name: &'static str) {
			assert!(
				check(chr),
				"checking `{}` (U+{:04X}) {}",
				chr,
				chr as u32,
				name
			);
		}

		pub fn check_range<F: Fn(char) -> bool>(a: char, b: char, check: F, name: &'static str) {
			for chr in a..=b {
				check_char(chr, &check, name);
			}
		}
	}
}
