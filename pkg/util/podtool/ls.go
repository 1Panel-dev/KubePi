package podtool

import "strings"

func (p *PodTool) ListFiles(path string) ([]File, error) {
	var files []File
	commands := []string{"ls", "-lF", "--full-time", path}

	res, err := p.ExecCommand(commands)
	if err != nil {
		return nil, err
	}
	result := string(res)
	array := strings.Split(result, "\n")
	i := 0
	if strings.Contains(array[0], "total") {
		i = 1
	}
	for ; i < len(array); i++ {
		line := array[i]
		fArray := strings.Fields(line)
		if len(fArray) >= 9 {
			f := File{
				Name:      fArray[8],
				Size:      fArray[4],
				Mode:      fArray[0],
				ModTime:   fArray[5] + " " + fArray[6],
				User:      fArray[2],
				UserGroup: fArray[3],
				IsDir:     strings.HasSuffix(fArray[8], "/"),
			}
			files = append(files, f)
		}
	}
	return files, nil
}
