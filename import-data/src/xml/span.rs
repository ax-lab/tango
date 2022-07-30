/// Represents a span of text from the input.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Span {
	/// Line number, starting at 1, for the beginning of the span.
	pub row: usize,
	/// Column number, starting at 1, for the beginning of the span.
	pub col: usize,
	/// Offset in bytes for the beginning of the span in the input.
	pub pos: usize,
	/// Total length in bytes for the span.
	pub len: usize,
}

impl Default for Span {
	fn default() -> Self {
		Span {
			row: 1,
			col: 1,
			pos: 0,
			len: 0,
		}
	}
}

impl Span {
	/// Returns a cursor helper to modify the span position based on input
	/// characters.
	pub fn cursor(self) -> Cursor {
		Cursor {
			span: self,
			last_was_cr: false,
		}
	}
}

/// Helper to track a span position based on input characters.
pub struct Cursor {
	pub span: Span,
	last_was_cr: bool,
}

impl Cursor {
	pub fn advance(&mut self, input: char, len: usize) {
		let is_cr = input == '\r';
		if input == '\n' || is_cr {
			if is_cr || !self.last_was_cr {
				self.span.col = 1;
				self.span.row += 1;
			}
			self.last_was_cr = is_cr;
		} else {
			self.span.col += 1;
		}
		self.span.pos += len;
	}
}

#[allow(unused_macros)]
macro_rules! span {
	($row:literal : $col:literal @ $pos:literal) => {
		Span {
			row: $row,
			col: $col,
			pos: $pos,
			len: 0,
		}
	};
}

#[cfg(test)]
mod tests {
	use super::*;

	#[test]
	pub fn defaults_to_first_position() {
		let span = Span::default();
		assert_eq!(
			span,
			Span {
				row: 1,
				col: 1,
				pos: 0,
				len: 0,
			}
		)
	}

	#[test]
	pub fn macro_creates_span() {
		let span = span!(2:3 @4);
		assert_eq!(
			span,
			Span {
				row: 2,
				col: 3,
				pos: 4,
				len: 0
			}
		)
	}

	#[test]
	pub fn cursor_tracks_characters() {
		let mut cursor = Span::default().cursor();

		// advances in the column
		cursor.advance('a', 1);
		assert_eq!(cursor.span, span!(1:2 @1));

		// advances by the given offset
		cursor.advance('気', 2);
		assert_eq!(cursor.span, span!(1:3 @3));

		// advances line
		cursor.advance('\n', 1);
		assert_eq!(cursor.span, span!(2:1 @4));

		// advances line with CR
		cursor.advance(' ', 1);
		assert_eq!(cursor.span, span!(2:2 @5));

		cursor.advance('\r', 1);
		assert_eq!(cursor.span, span!(3:1 @6));

		// supports CRLF sequence
		cursor.advance('\n', 1);
		assert_eq!(cursor.span, span!(3:1 @7));

		// LR sequence after CRLF
		cursor.advance('\n', 1);
		assert_eq!(cursor.span, span!(4:1 @8));
	}
}
