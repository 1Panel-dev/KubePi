package podtool

import "strings"

func (p *PodTool) EditFile(path string, content string) error {
	catCommand := []string{"cat", path}
	oldContent, err := p.ExecCommand(catCommand)
	if err != nil {
		return err
	}
	defer func() {
		p.ExecConfig.Stdin = nil
	}()
	p.ExecConfig.Stdin = strings.NewReader(content)
	_, err = p.ExecCommand([]string{"tee", path})
	if err != nil {
		p.ExecConfig.Stdin = strings.NewReader(string(oldContent))
		_, _ = p.ExecCommand([]string{"tee", path})
		return err
	}
	return nil
}
