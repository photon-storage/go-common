package log

import "github.com/logrusorgru/aurora"

var (
	au = aurora.NewAurora(true)
)

// Formats
//
//
// Bold or increased intensity (1).
func Bold(arg interface{}) aurora.Value {
	return au.Bold(arg)
}

// Faint, decreased intensity (2).
func Faint(arg interface{}) aurora.Value {
	return au.Faint(arg)
}

//
// DoublyUnderline or Bold off, double-underline
// per ECMA-48 (21).
func DoublyUnderline(arg interface{}) aurora.Value {
	return au.DoublyUnderline(arg)
}

// Fraktur, rarely supported (20).
func Fraktur(arg interface{}) aurora.Value {
	return au.Fraktur(arg)
}

//
// Italic, not widely supported, sometimes
// treated as inverse (3).
func Italic(arg interface{}) aurora.Value {
	return au.Italic(arg)
}

// Underline (4).
func Underline(arg interface{}) aurora.Value {
	return au.Underline(arg)
}

//
// SlowBlink, blinking less than 150
// per minute (5).
func SlowBlink(arg interface{}) aurora.Value {
	return au.SlowBlink(arg)
}

// RapidBlink, blinking 150+ per minute,
// not widely supported (6).
func RapidBlink(arg interface{}) aurora.Value {
	return au.RapidBlink(arg)
}

// Blink is alias for the SlowBlink.
func Blink(arg interface{}) aurora.Value {
	return au.Blink(arg)
}

//
// Reverse video, swap foreground and
// background colors (7).
func Reverse(arg interface{}) aurora.Value {
	return au.Reverse(arg)
}

// Inverse is alias for the Reverse
func Inverse(arg interface{}) aurora.Value {
	return au.Inverse(arg)
}

//
// Conceal, hidden, not widely supported (8).
func Conceal(arg interface{}) aurora.Value {
	return au.Conceal(arg)
}

// Hidden is alias for the Conceal
func Hidden(arg interface{}) aurora.Value {
	return au.Hidden(arg)
}

//
// CrossedOut, characters legible, but
// marked for deletion (9).
func CrossedOut(arg interface{}) aurora.Value {
	return au.CrossedOut(arg)
}

// StrikeThrough is alias for the CrossedOut.
func StrikeThrough(arg interface{}) aurora.Value {
	return au.StrikeThrough(arg)
}

//
// Framed (51).
func Framed(arg interface{}) aurora.Value {
	return au.Framed(arg)
}

// Encircled (52).
func Encircled(arg interface{}) aurora.Value {
	return au.Encircled(arg)
}

//
// Overlined (53).
func Overlined(arg interface{}) aurora.Value {
	return au.Overlined(arg)
}

//
// Foreground colors
//
//
// Black foreground color (30)
func Black(arg interface{}) aurora.Value {
	return au.Black(arg)
}

// Red foreground color (31)
func Red(arg interface{}) aurora.Value {
	return au.Red(arg)
}

// Green foreground color (32)
func Green(arg interface{}) aurora.Value {
	return au.Green(arg)
}

// Yellow foreground color (33)
func Yellow(arg interface{}) aurora.Value {
	return au.Yellow(arg)
}

// Brown foreground color (33)
//
// Deprecated: use Yellow instead, following specification
func Brown(arg interface{}) aurora.Value {
	return au.Brown(arg)
}

// Blue foreground color (34)
func Blue(arg interface{}) aurora.Value {
	return au.Blue(arg)
}

// Magenta foreground color (35)
func Magenta(arg interface{}) aurora.Value {
	return au.Magenta(arg)
}

// Cyan foreground color (36)
func Cyan(arg interface{}) aurora.Value {
	return au.Cyan(arg)
}

// White foreground color (37)
func White(arg interface{}) aurora.Value {
	return au.White(arg)
}

//
// Bright foreground colors
//
// BrightBlack foreground color (90)
func BrightBlack(arg interface{}) aurora.Value {
	return au.BrightBlack(arg)
}

// BrightRed foreground color (91)
func BrightRed(arg interface{}) aurora.Value {
	return au.BrightRed(arg)
}

// BrightGreen foreground color (92)
func BrightGreen(arg interface{}) aurora.Value {
	return au.BrightGreen(arg)
}

// BrightYellow foreground color (93)
func BrightYellow(arg interface{}) aurora.Value {
	return au.BrightYellow(arg)
}

// BrightBlue foreground color (94)
func BrightBlue(arg interface{}) aurora.Value {
	return au.BrightBlue(arg)
}

// BrightMagenta foreground color (95)
func BrightMagenta(arg interface{}) aurora.Value {
	return au.BrightMagenta(arg)
}

// BrightCyan foreground color (96)
func BrightCyan(arg interface{}) aurora.Value {
	return au.BrightCyan(arg)
}

// BrightWhite foreground color (97)
func BrightWhite(arg interface{}) aurora.Value {
	return au.BrightWhite(arg)
}

//
// Other
//
// Index of pre-defined 8-bit foreground color
// from 0 to 255 (38;5;n).
//
//       0-  7:  standard colors (as in ESC [ 30–37 m)
//       8- 15:  high intensity colors (as in ESC [ 90–97 m)
//      16-231:  6 × 6 × 6 cube (216 colors): 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
//     232-255:  grayscale from black to white in 24 steps
//
func Index(n uint8, arg interface{}) aurora.Value {
	return au.Index(n, arg)
}

// Gray from 0 to 23.
func Gray(n uint8, arg interface{}) aurora.Value {
	return au.Gray(n, arg)
}

//
// Background colors
//
//
// BgBlack background color (40)
func BgBlack(arg interface{}) aurora.Value {
	return au.BgBlack(arg)
}

// BgRed background color (41)
func BgRed(arg interface{}) aurora.Value {
	return au.BgRed(arg)
}

// BgGreen background color (42)
func BgGreen(arg interface{}) aurora.Value {
	return au.BgGreen(arg)
}

// BgYellow background color (43)
func BgYellow(arg interface{}) aurora.Value {
	return au.BgYellow(arg)
}

// BgBrown background color (43)
//
// Deprecated: use BgYellow instead, following specification
func BgBrown(arg interface{}) aurora.Value {
	return au.BgBrown(arg)
}

// BgBlue background color (44)
func BgBlue(arg interface{}) aurora.Value {
	return au.BgBlue(arg)
}

// BgMagenta background color (45)
func BgMagenta(arg interface{}) aurora.Value {
	return au.BgMagenta(arg)
}

// BgCyan background color (46)
func BgCyan(arg interface{}) aurora.Value {
	return au.BgCyan(arg)
}

// BgWhite background color (47)
func BgWhite(arg interface{}) aurora.Value {
	return au.BgWhite(arg)
}

//
// Bright background colors
//
// BgBrightBlack background color (100)
func BgBrightBlack(arg interface{}) aurora.Value {
	return au.BgBrightBlack(arg)
}

// BgBrightRed background color (101)
func BgBrightRed(arg interface{}) aurora.Value {
	return au.BgBrightRed(arg)
}

// BgBrightGreen background color (102)
func BgBrightGreen(arg interface{}) aurora.Value {
	return au.BgBrightGreen(arg)
}

// BgBrightYellow background color (103)
func BgBrightYellow(arg interface{}) aurora.Value {
	return au.BgBrightYellow(arg)
}

// BgBrightBlue background color (104)
func BgBrightBlue(arg interface{}) aurora.Value {
	return au.BgBrightBlue(arg)
}

// BgBrightMagenta background color (105)
func BgBrightMagenta(arg interface{}) aurora.Value {
	return au.BgBrightMagenta(arg)
}

// BgBrightCyan background color (106)
func BgBrightCyan(arg interface{}) aurora.Value {
	return au.BgBrightCyan(arg)
}

// BgBrightWhite background color (107)
func BgBrightWhite(arg interface{}) aurora.Value {
	return au.BgBrightWhite(arg)
}

//
// Other
//
// BgIndex of 8-bit pre-defined background color
// from 0 to 255 (48;5;n).
//
//       0-  7:  standard colors (as in ESC [ 40–47 m)
//       8- 15:  high intensity colors (as in ESC [100–107 m)
//      16-231:  6 × 6 × 6 cube (216 colors): 16 + 36 × r + 6 × g + b (0 ≤ r, g, b ≤ 5)
//     232-255:  grayscale from black to white in 24 steps
//
func BgIndex(n uint8, arg interface{}) aurora.Value {
	return au.BgIndex(n, arg)
}

// BgGray from 0 to 23.
func BgGray(n uint8, arg interface{}) aurora.Value {
	return au.BgGray(n, arg)
}
