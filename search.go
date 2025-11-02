package main

import (
    "regexp"
    "strings"
    "unicode"
)

func NormalizePokemonName(name string) string {
    n := strings.ToLower(name)

    // Handle gender signs
    n = strings.ReplaceAll(n, "♀", "-f")
    n = strings.ReplaceAll(n, "♂", "-m")

    // Replace spaces and punctuation with hyphens
    n = strings.ReplaceAll(n, " ", "-")
    n = strings.ReplaceAll(n, ".", "")
    n = strings.ReplaceAll(n, ":", "")
    n = strings.ReplaceAll(n, "'", "")
    n = strings.ReplaceAll(n, "’", "")
    n = strings.ReplaceAll(n, "–", "-")
    n = strings.ReplaceAll(n, "—", "-")

    // Replace accented characters
    replacer := strings.NewReplacer(
        "é", "e",
        "è", "e",
        "á", "a",
        "à", "a",
        "í", "i",
        "ó", "o",
        "ú", "u",
    )
    n = replacer.Replace(n)

    // Remove any character not a-z, 0-9, or '-'
    re := regexp.MustCompile(`[^a-z0-9\-]`)
    n = re.ReplaceAllStringFunc(n, func(s string) string {
        if unicode.IsLetter([]rune(s)[0]) || unicode.IsDigit([]rune(s)[0]) {
            return s
        }
        return ""
    })

    // Collapse multiple hyphens
    n = strings.ReplaceAll(n, "--", "-")

    return n
}
