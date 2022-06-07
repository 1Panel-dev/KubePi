package podtool

func (p *PodTool) EditFile(path string, content string) error {
	catCommand := []string{"cat", path}
	oldContent, err := p.ExecCommand(catCommand)
	if err != nil {
		return err
	}
	_, err = p.ExecCommand([]string{"sh", "-c", "cat /dev/null > " + path})
	if err != nil {
		return err
	}
	command := "echo '" + content + "' >> " + path
	_, err = p.ExecCommand([]string{"sh", "-c", command})
	if err != nil {
		command2 := "echo '" + string(oldContent) + "' >> " + path
		_, _ = p.ExecCommand([]string{"sh", "-c", command2})
		return err
	}
	return nil
}
