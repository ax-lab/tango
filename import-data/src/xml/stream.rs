/// Abstracts an underlying Reader and provides char-based methods to consume
/// the XML input.
pub trait Stream {
	type State: Default;

	fn peek_char(&mut self, state: &mut Self::State) -> Option<char>;

	fn read_char(&mut self, state: &mut Self::State) -> Option<(char, usize)>;
}

impl<'a> Stream for &'a str {
	type State = usize;

	fn peek_char(&mut self, state: &mut Self::State) -> Option<char> {
		let input = &self[*state..];
		input.chars().next()
	}

	fn read_char(&mut self, state: &mut Self::State) -> Option<(char, usize)> {
		let input = &self[*state..];
		let mut chars = input.char_indices();
		match chars.next() {
			None => None,
			Some((_, char)) => {
				let char_length = match chars.next() {
					None => input.len(),
					Some((next_offset, _)) => next_offset,
				};
				*state += char_length;
				Some((char, char_length))
			}
		}
	}
}
