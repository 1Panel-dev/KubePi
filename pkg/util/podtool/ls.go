package podtool

import "strings"

func (p *PodTool) ListFiles(path string) ([]File, error) {
	var files []File
	commands := []string{"ls", "-l", "--full-time", path}

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
				Size:      fArray[4],
				Mode:      fArray[0],
				ModTime:   fArray[5] + " " + fArray[6],
				User:      fArray[2],
				UserGroup: fArray[3],
			}
			var name string
			if strings.HasPrefix(f.Mode, "l") && len(fArray) > 10 {
				f.Link = getLink(line)
				name = getLinkName(fArray[8], line)
			} else {
				name = getName(fArray[8], line)
			}
			if strings.HasPrefix(f.Mode, "d") {
				f.IsDir = true
			}

			f.Name = name
			files = append(files, f)
		}
	}
	return files, nil
}

func getName(sub string, line string) string {
	return strings.TrimSpace(line[strings.LastIndex(line, sub):])
}

func getLink(line string) string {
	const linkTag = "->"
	in := strings.Index(line, linkTag)
	return strings.TrimSpace(line[in+len(linkTag):])
}

func getLinkName(sub string, line string) string {
	const linkTag = "->"
	linkIn := strings.Index(line, linkTag)
	nameIn := strings.Index(line, sub)
	return strings.TrimSpace(line[nameIn:linkIn])
}
