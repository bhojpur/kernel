package animal

type CowTemplate struct {
}

func (cowTemplate CowTemplate) Get() string {
	cowTempalte := `
	 \   ^__^
	  \  (oo)\_______
		 (__)\       )\/\
			 ||----w |
			 ||     ||
 `
	return cowTempalte
}
