package shell

func (s *Shell) StartInteractive() {
	if s.started {
		return
	}
	s.started = true
	s.Println("Welcome to Zash!")

	for {
		i := s.WaitInputWithPrompt(s.Context())
		if err := s.Run(i); err != nil {
			s.Println(err.Error())
		}
		s.prompt.Update()
	}
}
