package podtool

import "errors"

func (p *PodTool) CatFile(filePath string) ([]byte, error) {
	catCommand := []string{"cat", filePath}
	content, err := p.ExecCommand(catCommand)
	if err != nil {
		return nil, err
	}
	if detectBinary(content) {
		return nil, errors.New("file can not read")
	}
	return content, nil
}

func detectBinary(buf []byte) bool {
	var whiteByte int = 0
	n := min(1024, len(buf))
	for i := 0; i < n; i++ {
		if (buf[i] >= 0x20) || buf[i] == 9 || buf[i] == 10 || buf[i] == 13 {
			whiteByte++
		} else if buf[i] <= 6 || (buf[i] >= 14 && buf[i] <= 31) {
			return true
		}
	}

	if whiteByte >= 1 {
		return false
	}
	return true
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
