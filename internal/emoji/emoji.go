/*
 * owobot - Your server's guardian and entertainer
 * Copyright (C) 2023 owobot Contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package emoji

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/rivo/uniseg"
)

var (
	customEmojiRegex  = regexp.MustCompile(`^<(a?):(\w+):(\d+)>$`)
	unicodeEmojiTable = &unicode.RangeTable{
		R32: []unicode.Range32{
			{ // Enclosed Alphanumeric Supplement
				Lo:     0x1F100,
				Hi:     0x1F1FF,
				Stride: 1,
			},
			{ // Miscellaneous Symbols and Pictographs
				Lo:     0x1F300,
				Hi:     0x1F5FF,
				Stride: 1,
			},
			{ // Emoticons
				Lo:     0x1F600,
				Hi:     0x1F64F,
				Stride: 1,
			},
			{ // Transport and Map Symbols
				Lo:     0x1F680,
				Hi:     0x1F6FF,
				Stride: 1,
			},
			{ // Geometric Shapes Extended
				Lo:     0x1F780,
				Hi:     0x1F7FF,
				Stride: 1,
			},
			{ // Supplemental Symbols and Pictographs
				Lo:     0x1F900,
				Hi:     0x1F9FF,
				Stride: 1,
			},
			{ // Symbols and Pictographs Extended-A
				Lo:     0x1FA70,
				Hi:     0x1FAFF,
				Stride: 1,
			},
		},
		R16: []unicode.Range16{
			{ // Zero-width characters
				Lo:     0x200B,
				Hi:     0x200D,
				Stride: 1,
			},
			{ // Miscellaneous Technical
				Lo:     0x2300,
				Hi:     0x23FF,
				Stride: 1,
			},
			{ // Miscellaneous Symbols
				Lo:     0x2600,
				Hi:     0x26FF,
				Stride: 1,
			},
			{ // Dingbats
				Lo:     0x2700,
				Hi:     0x27BF,
				Stride: 1,
			},
			{ // Miscellaneous Symbols and Arrows
				Lo:     0x2B00,
				Hi:     0x2BFF,
				Stride: 1,
			},
			{ // Variation Selectors
				Lo:     0xFE00,
				Hi:     0xFE0F,
				Stride: 1,
			},
		},
	}
)

// Emoji represents a Discord emoji.
type Emoji struct {
	Name       string
	ID         string
	IsAnimated bool
	IsCustom   bool
}

// APIFormat returns a string that represents the emoji
// in discord API requests.
func (e *Emoji) APIFormat() string {
	if e.IsCustom {
		return e.Name + ":" + e.ID
	} else {
		return e.Name
	}
}

// MessageFormat returns a string that represents the emoji
// in discord messages.
func (e *Emoji) MessageFormat() string {
	if e.IsCustom {
		var sb strings.Builder
		sb.WriteByte('<')
		if e.IsAnimated {
			sb.WriteByte('a')
		}
		sb.WriteByte(':')
		sb.WriteString(e.Name)
		sb.WriteByte(':')
		sb.WriteString(e.ID)
		sb.WriteByte('>')
		return sb.String()
	} else {
		return e.Name
	}
}

// Parse parses a single Discord emoji. It can handle both
// unicode emoji and custom emoji.
func Parse(s string) (Emoji, bool) {
	s = strings.TrimSpace(s)

	if isEmoji(s) {
		// This string should only contain a single emoji.
		// If it has more, return false.
		if uniseg.GraphemeClusterCount(s) != 1 {
			return Emoji{}, false
		}
		return Emoji{Name: s}, true
	}

	if customEmojiRegex.MatchString(s) {
		matches := customEmojiRegex.FindStringSubmatch(s)
		return Emoji{
			Name:       matches[2],
			ID:         matches[3],
			IsAnimated: matches[1] == "a",
			IsCustom:   true,
		}, true
	}

	return Emoji{}, false
}

// isEmoji checks to make sure all the characters in the string
// are within the unicode emoji range table.
func isEmoji(s string) bool {
	for _, char := range s {
		if !unicode.In(char, unicodeEmojiTable) {
			return false
		}
	}
	return true
}
