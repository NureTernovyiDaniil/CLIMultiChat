package formatter

import (
	"fmt"
	"regexp"
	"strings"
)

// ToSlackMarkdown converts standard markdown to Slack's mrkdwn format
// Slack uses:
// - *bold* for bold (single asterisk)
// - _italic_ for italic
// - `code` for inline code
// - ```code block``` for code blocks
// - <url|text> for links (or just <url> for plain links)
func ToSlackMarkdown(text string) string {
	result := text

	// Store markdown elements with unique markers
	type replacement struct {
		marker string
		final  string
	}

	var replacements []replacement
	counter := 0

	// Helper function to create unique marker
	makeMarker := func() string {
		marker := fmt.Sprintf("\x00SLACKMARK%d\x00", counter)
		counter++
		return marker
	}

	// 1. Preserve code blocks ```...```
	codeBlockRe := regexp.MustCompile("```([\\s\\S]*?)```")
	result = codeBlockRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		replacements = append(replacements, replacement{marker, match})
		return marker
	})

	// 2. Preserve inline code `...`
	codeRe := regexp.MustCompile("`([^`]+)`")
	result = codeRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		replacements = append(replacements, replacement{marker, match})
		return marker
	})

	// 3. Convert links [text](url) to <url|text>
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	result = linkRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		parts := linkRe.FindStringSubmatch(match)
		if len(parts) == 3 {
			slackLink := fmt.Sprintf("<%s|%s>", parts[2], parts[1])
			replacements = append(replacements, replacement{marker, slackLink})
		} else {
			replacements = append(replacements, replacement{marker, match})
		}
		return marker
	})

	// 4. Convert bold **...** to *...*
	boldRe := regexp.MustCompile(`\*\*([^\*\n]+?)\*\*`)
	result = boldRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		parts := boldRe.FindStringSubmatch(match)
		if len(parts) == 2 {
			converted := "*" + parts[1] + "*"
			replacements = append(replacements, replacement{marker, converted})
		} else {
			replacements = append(replacements, replacement{marker, match})
		}
		return marker
	})

	// 5. Convert remaining single asterisk italic *...* to _..._
	// We need to be careful not to match the bold markers we just created
	italicAsteriskRe := regexp.MustCompile(`(?:^|[^\*])(\*([^\*\n]+?)\*)(?:[^\*]|$)`)
	result = italicAsteriskRe.ReplaceAllStringFunc(result, func(match string) string {
		parts := italicAsteriskRe.FindStringSubmatch(match)
		if len(parts) >= 3 {
			marker := makeMarker()
			converted := "_" + parts[2] + "_"
			// Preserve leading/trailing characters if any
			leadChar := ""
			trailChar := ""
			if len(match) > len(parts[1]) {
				if !strings.HasPrefix(match, parts[1]) {
					leadChar = string(match[0])
				}
				if !strings.HasSuffix(match, parts[1]) {
					trailChar = string(match[len(match)-1])
				}
			}
			replacements = append(replacements, replacement{marker, converted})
			return leadChar + marker + trailChar
		}
		return match
	})

	// 6. Underscores for italic _..._ stay as is (Slack uses same format)
	// No conversion needed

	// 7. Restore preserved elements
	for _, r := range replacements {
		result = strings.Replace(result, r.marker, r.final, 1)
	}

	return result
}

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
		marker string
		final  string
	}

	var replacements []replacement
	counter := 0

	// Helper function to create unique marker
	makeMarker := func() string {
		marker := fmt.Sprintf("\x00MARK%d\x00", counter)
		counter++
		return marker
	}

	// 1. Preserve code blocks ```...``` (don't escape content inside)
	codeBlockRe := regexp.MustCompile("```([\\s\\S]*?)```")
	result = codeBlockRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		// Code blocks content should not be escaped
		replacements = append(replacements, replacement{marker, match})
		return marker
	})

	// 2. Preserve inline code `...` (don't escape content inside)
	codeRe := regexp.MustCompile("`([^`]+)`")
	result = codeRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		// Code content should not be escaped
		replacements = append(replacements, replacement{marker, match})
		return marker
	})

	// 3. Preserve and process links [text](url)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	result = linkRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		// Extract link parts - will escape later
		parts := linkRe.FindStringSubmatch(match)
		if len(parts) == 3 {
			// Store the link structure, escape later
			finalLink := fmt.Sprintf("\x01LINK:%s\x02%s\x03", parts[1], parts[2])
			replacements = append(replacements, replacement{marker, finalLink})
		} else {
			replacements = append(replacements, replacement{marker, match})
		}
		return marker
	})

	// 4. Preserve bold text **...**
	boldRe := regexp.MustCompile(`\*\*([^\*\n]+?)\*\*`)
	result = boldRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		// Convert **text** -> *text*, will escape content later
		parts := boldRe.FindStringSubmatch(match)
		if len(parts) == 2 {
			// Mark bold content for later escaping
			converted := fmt.Sprintf("\x01BOLD:%s\x03", parts[1])
			replacements = append(replacements, replacement{marker, converted})
		} else {
			replacements = append(replacements, replacement{marker, match})
		}
		return marker
	})

	// 5. Preserve italic *...* (single asterisk, not double)
	italicRe := regexp.MustCompile(`(?:^|[^\*])(\*([^\*\n]+?)\*)(?:[^\*]|$)`)
	result = italicRe.ReplaceAllStringFunc(result, func(match string) string {
		// Find the actual italic part
		parts := italicRe.FindStringSubmatch(match)
		if len(parts) >= 3 {
			marker := makeMarker()
			// Mark italic content for later escaping
			converted := fmt.Sprintf("\x01ITALIC:%s\x03", parts[2])
			// Preserve leading/trailing characters if any
			leadChar := ""
			trailChar := ""
			if len(match) > len(parts[1]) {
				if !strings.HasPrefix(match, parts[1]) {
					leadChar = string(match[0])
				}
				if !strings.HasSuffix(match, parts[1]) {
					trailChar = string(match[len(match)-1])
				}
			}
			replacements = append(replacements, replacement{marker, converted})
			return leadChar + marker + trailChar
		}
		return match
	})

	// 6. Preserve italic _..._
	underscoreItalicRe := regexp.MustCompile(`_([^_\n]+?)_`)
	result = underscoreItalicRe.ReplaceAllStringFunc(result, func(match string) string {
		marker := makeMarker()
		parts := underscoreItalicRe.FindStringSubmatch(match)
		if len(parts) == 2 {
			// Mark italic content for later escaping
			converted := fmt.Sprintf("\x01ITALIC:%s\x03", parts[1])
			replacements = append(replacements, replacement{marker, converted})
		} else {
			replacements = append(replacements, replacement{marker, match})
		}
		return marker
	})

	// 7. Escape all remaining special characters in plain text
	result = escapeSpecialChars(result)

	// 8. Restore preserved elements and process marked content
	for _, r := range replacements {
		final := r.final

		// Process BOLD markers
		if strings.HasPrefix(final, "\x01BOLD:") {
			content := strings.TrimPrefix(final, "\x01BOLD:")
			content = strings.TrimSuffix(content, "\x03")
			escapedContent := escapeSpecialCharsSimple(content)
			final = "*" + escapedContent + "*"
		}

		// Process ITALIC markers
		if strings.HasPrefix(final, "\x01ITALIC:") {
			content := strings.TrimPrefix(final, "\x01ITALIC:")
			content = strings.TrimSuffix(content, "\x03")
			escapedContent := escapeSpecialCharsSimple(content)
			final = "_" + escapedContent + "_"
		}

		// Process LINK markers
		if strings.HasPrefix(final, "\x01LINK:") {
			parts := strings.Split(final, "\x02")
			if len(parts) == 2 {
				text := strings.TrimPrefix(parts[0], "\x01LINK:")
				url := strings.TrimSuffix(parts[1], "\x03")
				escapedText := escapeSpecialCharsSimple(text)
				escapedURL := escapeSpecialCharsSimple(url)
				final = "[" + escapedText + "](" + escapedURL + ")"
			}
		}

		result = strings.Replace(result, r.marker, final, 1)
	}

	return result
}

// escapeSpecialCharsSimple escapes special characters without checking for markers
func escapeSpecialCharsSimple(text string) string {
	specialChars := map[rune]bool{
		'_': true, '*': true, '[': true, ']': true,
		'(': true, ')': true, '~': true, '`': true,
		'>': true, '#': true, '+': true, '-': true,
		'=': true, '|': true, '{': true, '}': true,
		'.': true, '!': true,
	}

	result := strings.Builder{}
	for _, r := range text {
		if specialChars[r] {
			result.WriteRune('\\')
		}
		result.WriteRune(r)
	}

	return result.String()
}

// escapeSpecialChars escapes special characters for Telegram MarkdownV2
// but preserves our markers
func escapeSpecialChars(text string) string {
	specialChars := map[rune]bool{
		'_': true, '*': true, '[': true, ']': true,
		'(': true, ')': true, '~': true, '`': true,
		'>': true, '#': true, '+': true, '-': true,
		'=': true, '|': true, '{': true, '}': true,
		'.': true, '!': true,
	}

	result := strings.Builder{}
	runes := []rune(text)
	i := 0

	for i < len(runes) {
		// Skip markers (null byte followed by 'M')
		if i < len(runes)-1 && runes[i] == 0 && runes[i+1] == 'M' {
			// Found start of marker, find end (next null byte)
			result.WriteRune(runes[i])
			i++
			result.WriteRune(runes[i])
			i++

			// Copy everything until the closing null byte
			for i < len(runes) && runes[i] != 0 {
				result.WriteRune(runes[i])
				i++
			}
			if i < len(runes) {
				result.WriteRune(runes[i]) // Write closing null byte
				i++
			}
			continue
		}

		// Check if current character needs escaping
		if specialChars[runes[i]] {
			result.WriteRune('\\')
		}
		result.WriteRune(runes[i])
		i++
	}

	return result.String()
}
