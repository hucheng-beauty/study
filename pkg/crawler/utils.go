package crawler

// TODO: restart is not storage data
var visitedUrls = make(map[string]bool, 0)

func isDuplicate(url string) bool {
    if visitedUrls[url] {
        return true
    }

    visitedUrls[url] = true
    return false
}
