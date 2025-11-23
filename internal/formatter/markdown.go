package formatter

import (
	"regexp"
	"strings"
)

// ToTelegramMarkdown converts standard markdown to Telegram's MarkdownV2 format
// Telegram MarkdownV2 requires escaping of special characters: _*[]()~`>#+-=|{}.!
// and uses:
// - *bold* for bold
// - _italic_ for italic
// - __underline__ for underline
// - ~strikethrough~ for strikethrough
// - ||spoiler|| for spoiler
// - `code` for inline code
// - ```code block``` for code blocks
// - [text](url) for links
func ToTelegramMarkdown(text string) string {
	// First, protect markdown syntax we want to keep
	protectedPatterns := []struct {
		pattern     *regexp.Regexp
		placeholder string
		matches     []string
	}{
		{regexp.MustCompile(`\*\*([^\*]+)\*\*`), "BOLD_%d_", nil},           // **bold**
		{regexp.MustCompile(`\*([^\*]+)\*`), "ITALIC_%d_", nil},             // *italic*
		{regexp.MustCompile(`_([^_]+)_`), "UNDERSCORE_%d_", nil},            // _italic_
		{regexp.MustCompile("`([^`]+)`"), "CODE_%d_", nil},                  // `code`
		{regexp.MustCompile("```([^`]+)```"), "CODEBLOCK_%d_", nil},         // ```code```
		{regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`), "LINK_%d_", nil}, // [text](url)
	}

	result := text

	// Extract and protect markdown patterns
	for i := range protectedPatterns {
		matches := protectedPatterns[i].pattern.FindAllString(result, -1)
		protectedPatterns[i].matches = matches

		for j, match := range matches {
			placeholder := strings.Replace(protectedPatterns[i].placeholder, "%d", string(rune(j)), 1)
			result = strings.Replace(result, match, placeholder, 1)
		}
	}

	// Escape special characters
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	// Restore markdown patterns
	for i := range protectedPatterns {
		for j, match := range protectedPatterns[i].matches {
			placeholder := strings.Replace(protectedPatterns[i].placeholder, "%d", string(rune(j)), 1)

			// Convert to Telegram format
			var converted string
			switch i {
			case 0: // **bold** -> *bold*
				converted = regexp.MustCompile(`\*\*([^\*]+)\*\*`).ReplaceAllString(match, "*$1*")
			case 1, 2: // *italic* or _italic_ -> _italic_
				converted = regexp.MustCompile(`[\*_]([^\*_]+)[\*_]`).ReplaceAllString(match, "_$1_")
			default:
				converted = match
			}

			result = strings.Replace(result, "\\"+placeholder, converted, 1)
		}
	}

	return result
}
