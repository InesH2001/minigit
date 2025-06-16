package utils

func GetLine(lines []string, i int) string {
    if i < len(lines) {
        return lines[i]
    }
    return ""
}