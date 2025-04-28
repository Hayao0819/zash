package shell

func (s *Shell) StartInteractive() {
	if s.started {
		return
	}
	s.started = true
	s.Println("Welcome to Zash!")

	for {
		i, err := s.prompt.WaitInput()
		if err != nil {
			s.Println(err.Error())
			continue
		}
		if err := s.Run(i); err != nil {
			s.Println(err.Error())
		}
		s.prompt.Update()
	}
}
