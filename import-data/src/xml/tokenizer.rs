use super::{chars, span, stream::Stream, Span};

/// Low-level XML tokenizer.
///
/// Provides support for consuming an input stream by XML tokens and tracking
/// the [`Span`] for the token.
pub struct Tokenizer<T: Stream> {
	cursor: span::Cursor,
	input: T,
	state: T::State,
}

impl<T: Stream> Tokenizer<T> {
	pub fn new(input: T) -> Self {
		let cursor = Span::default().cursor();
		Tokenizer {
			cursor: cursor,
			input: input,
			state: T::State::default(),
		}
	}

	pub fn span(&self) -> Span {
		return self.cursor.span;
	}

	pub fn set_span(&mut self, span: Span) {
		debug_assert!(span.len == 0);
		self.cursor.span = span;
	}

	pub fn read_name(&mut self) -> Option<Span> {
		self.read_if_while(chars::is_name_start, chars::is_name)
	}

	pub fn peek_char(&mut self) -> Option<char> {
		self.input.peek_char(&mut self.state)
	}

	fn read_if_while<PI, PW>(&mut self, predicate_if: PI, predicate_while: PW) -> Option<Span>
	where
		PI: Fn(char) -> bool,
		PW: Fn(char) -> bool,
	{
		if let Some(char) = self.peek_char() {
			if predicate_if(char) {
				let mut span = self.cursor.span;
				self.skip_chars(predicate_while);
				span.len = self.cursor.span.pos - span.pos;
				return Some(span);
			}
		}
		None
	}

	fn skip_chars<P>(&mut self, predicate: P)
	where
		P: Fn(char) -> bool,
	{
		while let Some(char) = self.peek_char() {
			if predicate(char) {
				self.read_char();
			} else {
				break;
			}
		}
	}

	fn read_char(&mut self) -> Option<char> {
		if let Some((char, length)) = self.input.read_char(&mut self.state) {
			self.cursor.advance(char, length);
			Some(char)
		} else {
			None
		}
	}
}

pub enum Token {}

pub enum TokenType {}

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	fn creates_from_str() {
		let mut tokenizer = Tokenizer::new("abc");
		assert_eq!(tokenizer.peek_char(), Some('a'));
		assert_eq!(tokenizer.span(), Span::default());
	}

	#[test]
	fn set_span_updates_position() {
		let mut tokenizer = Tokenizer::new("");
		tokenizer.set_span(Span {
			pos: 99,
			..Span::default()
		});
		assert_eq!(tokenizer.span().pos, 99);
	}

	/// Helper macro to test read methods
	macro_rules! check_read {
		($expr:ident => $($tail:tt)*) => {
			let header = concat!(concat!("Tokenizer::", stringify!($expr), ""));
			check_read!(impl: header, |t| t.$expr() => $($tail)*);
		};

		(impl: $header:ident, $expr:expr => ) => {};

		(impl: $header:ident, $expr:expr => $text:literal, $($tail:tt)*) => {
			let header = format!("{}({})", $header, stringify!($text));
			helper::check_read_full($text, $expr, &header);
			check_read!(impl: $header, $expr => $($tail)*);
		};

		(impl: $header:ident, $expr:expr => ! $text:literal, $($tail:tt)*) => {
			let header = format!("{}(!{})", $header, stringify!($text));
			helper::check_read_fail($text, $expr, &header);
			check_read!(impl: $header, $expr => $($tail)*);
		};

		(impl: $header:ident, $expr:expr => $text:literal in $input:literal, $($tail:tt)*) => {
			let header = format!("{}({} in {})", $header, stringify!($text), stringify!($input));
			helper::check_read_partial($text, $input, $expr, &header);
			check_read!(impl: $header, $expr => $($tail)*);
		};
	}

	#[test]
	fn parse_name() {
		check_read!(
			read_name =>
			"a",
			"abc",
			"abc123",
			"abc-123",
			"abc_123",
			"abc.123",
			"abc·123",
			"abc:123",
			"言葉",
			!"",
			!" ",
			!"123",
			"abc" in "abc\u{037E}123",
		);
	}

	mod helper {
		use super::*;

		/// Used as the starting point for the tokenizer to test that tokens
		/// return a proper position.
		pub const BASE_SPAN: Span = Span {
			len: 0,
			pos: 10,
			row: 20,
			col: 30,
		};

		/// Check that the tokenizer reads the given input completely using
		/// the provided callback.
		pub fn check_read_full<F>(input: &str, callback: F, header: &str)
		where
			F: FnOnce(&mut Tokenizer<&str>) -> Option<Span> + std::panic::UnwindSafe,
		{
			let expected = Span {
				len: input.len(),
				..BASE_SPAN
			};
			let expected_next_char = Some('\u{0000}');
			let expected_next_span = get_next_span(expected, input);

			let (output, next_char, next_span) = tokenize(input, "\u{0000}", callback, header);

			assert!(output.is_some(), "{}: failed to parse", header);

			let actual = output.unwrap();

			assert!(
				actual.len == expected.len,
				"{}: did not parse expected input (parsed {}, expected {})",
				header,
				actual.len,
				expected.len,
			);

			assert!(
				actual == expected,
				"{}: output `{:?}`, expected `{:?}`",
				header,
				actual,
				expected,
			);

			assert!(
				next_char == expected_next_char,
				"{}: unexpected character after parse: {:?}",
				header,
				next_char,
			);

			assert!(
				next_span == expected_next_span,
				"{}: unexpected position after parse: it was {:?}, expected {:?}",
				header,
				next_span,
				expected_next_span,
			);
		}

		/// Checks that the tokenizer does not read the given input using the
		/// provided callback.
		pub fn check_read_fail<F>(input: &str, callback: F, header: &str)
		where
			F: FnOnce(&mut Tokenizer<&str>) -> Option<Span> + std::panic::UnwindSafe,
		{
			let expected_next_char = input.chars().next();
			let expected_next_span = BASE_SPAN;

			let (output, next_char, next_span) = tokenize(input, "", callback, header);
			assert!(
				output.is_none(),
				"{}: should not parse, instead parsed {:?}",
				header,
				output,
			);

			assert!(
				next_char == expected_next_char,
				"{}: unexpected character after parse: {:?}",
				header,
				next_char,
			);

			assert!(
				next_span == expected_next_span,
				"{}: moved position after failed parse: {:?}",
				header,
				next_span,
			);
		}

		/// Check that the tokenizer partially matches input using the given
		/// callback.
		pub fn check_read_partial<F>(text: &str, input: &str, callback: F, header: &str)
		where
			F: FnOnce(&mut Tokenizer<&str>) -> Option<Span> + std::panic::UnwindSafe,
		{
			let expected = Span {
				len: text.len(),
				..BASE_SPAN
			};
			let expected_next_char = input[text.len()..].chars().next();
			let expected_next_span = get_next_span(expected, text);

			let (output, next_char, next_span) = tokenize(input, "\u{0000}", callback, header);

			assert!(output.is_some(), "{}: failed to parse", header);
			let actual = output.unwrap();

			assert!(
				actual.len == expected.len,
				"{}: did not parse expected input (parsed {}, expected {})",
				header,
				actual.len,
				expected.len,
			);

			assert!(
				actual == expected,
				"{}: output `{:?}`, expected `{:?}`",
				header,
				actual,
				expected,
			);

			assert!(
				next_char == expected_next_char,
				"{}: unexpected {:?} character after parse, expected {:?}",
				header,
				next_char,
				expected_next_char,
			);

			assert!(
				next_span == expected_next_span,
				"{}: unexpected position after parse: it was {:?}, expected {:?}",
				header,
				next_span,
				expected_next_span,
			);
		}

		/// Uses a callback to tokenize the given input.
		fn tokenize<F>(
			input: &str,
			suffix: &str,
			callback: F,
			header: &str,
		) -> (Option<Span>, Option<char>, Span)
		where
			F: FnOnce(&mut Tokenizer<&str>) -> Option<Span> + std::panic::UnwindSafe,
		{
			let guarded_input = format!("{}{}", input, suffix);
			let result = std::panic::catch_unwind(|| {
				let mut tokenizer = Tokenizer::new(guarded_input.as_str());
				tokenizer.set_span(BASE_SPAN);
				let output = callback(&mut tokenizer);
				let next = tokenizer.peek_char();
				(output, next, tokenizer.cursor.span)
			});

			let (output, next, next_pos) = match result {
				Err(_) => panic!("{}: tokenizer panicked", header),
				Ok(result) => result,
			};

			(output, next, next_pos)
		}

		fn get_next_span(span: Span, input: &str) -> Span {
			let mut cursor = span.cursor();
			for chr in input.chars() {
				cursor.advance(chr, 0);
			}
			cursor.span.pos += input.len();
			Span {
				len: 0,
				..cursor.span
			}
		}
	}
}
