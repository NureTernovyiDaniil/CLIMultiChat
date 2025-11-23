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
	result := text

	// Store markdown elements with unique markers
	type replacement struct {
		original string
		marker   string
		final    string
	}

	var replacements []replacement
	counter := 0

	// 1. Preserve code blocks ```...```
	codeBlockRe := regexp.MustCompile("```([^`]+)```")
	for _, match := range codeBlockRe.FindAllString(result, -1) {
		marker := "\x00CODEBLOCK" + string(rune(counter)) + "\x00"
		replacements = append(replacements, replacement{match, marker, match})
		result = strings.Replace(result, match, marker, 1)
		counter++
	}

	// 2. Preserve inline code `...`
	codeRe := regexp.MustCompile("`([^`]+)`")
	for _, match := range codeRe.FindAllString(result, -1) {
		marker := "\x00CODE" + string(rune(counter)) + "\x00"
		replacements = append(replacements, replacement{match, marker, match})
		result = strings.Replace(result, match, marker, 1)
		counter++
	}

	// 3. Preserve links [text](url)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	for _, match := range linkRe.FindAllString(result, -1) {
		marker := "\x00LINK" + string(rune(counter)) + "\x00"
		replacements = append(replacements, replacement{match, marker, match})
		result = strings.Replace(result, match, marker, 1)
		counter++
	}

	// 4. Preserve bold text **...**
	boldRe := regexp.MustCompile(`\*\*([^\*]+)\*\*`)
	for _, match := range boldRe.FindAllString(result, -1) {
		marker := "\x00BOLD" + string(rune(counter)) + "\x00"
		// Convert **text** -> *text*
		converted := boldRe.ReplaceAllString(match, "*$1*")
		replacements = append(replacements, replacement{match, marker, converted})
		result = strings.Replace(result, match, marker, 1)
		counter++
	}

	// 5. Preserve italic _..._ or *...*
	italicRe := regexp.MustCompile(`(?:^|[^\*])(\*[^\*]+\*)|(_[^_]+_)`)
	for _, match := range italicRe.FindAllString(result, -1) {
		// Remove possible leading character (not *)
		cleanMatch := strings.TrimLeft(match, " \t\n\r")
		if strings.HasPrefix(cleanMatch, "*") || strings.HasPrefix(cleanMatch, "_") {
			marker := "\x00ITALIC" + string(rune(counter)) + "\x00"
			// Convert to _text_
			var converted string
			if strings.HasPrefix(cleanMatch, "*") {
				converted = regexp.MustCompile(`\*([^\*]+)\*`).ReplaceAllString(cleanMatch, "_$1_")
			} else {
				converted = cleanMatch
			}
			prefix := strings.TrimSuffix(match, cleanMatch)
			replacements = append(replacements, replacement{cleanMatch, marker, converted})
			result = strings.Replace(result, match, prefix+marker, 1)
			counter++
		}
	}

	// 6. Escape all special characters
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	// 7. Restore preserved elements
	for _, r := range replacements {
		result = strings.Replace(result, r.marker, r.final, 1)
	}

	return result
}
