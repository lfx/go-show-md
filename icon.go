package main

import "encoding/base64"

// getIconData returns the byte slice for the menu bar icon
func getIconData() []byte {
	// 16x16 black dot PNG base64
	b64 := "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAABoSURBVDhP7YxBCsAgDEX9/6a7U/AChZ5EfJ0G2o048C08QshI+k5K6QQA+AA9W1cCKkIu76dA7w1M1d6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N6C1N5/wzcsHClC1qQAAAAASUVORK5CYII="

	data, _ := base64.StdEncoding.DecodeString(b64)
	return data
}
