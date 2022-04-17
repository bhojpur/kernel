package animal

type DemonTemplate struct {
}

func (dt DemonTemplate) Get() string {
	demonTemplate := `
	\         ,        ,
	 \       /(        )'
	  \      \\ \\___   / |
			 /- _  '-/  '
			(/\\/ \\ \\   /\\
			/ /   | '    \
			O O   ) /    |
			'-^--''<     '
		   (_.)  _  )   /
			'.___/'    /
			  '-----' /
 <----.     __ / __   \\
 <----|====O)))==) \\) /====
 <----'    '--' '.__,' \\
			  |        |
			   \\       /
		 ______( (_  / \\______
	   ,'  ,-----'   |        \\
	   '--{__________)        \\
  `
	return demonTemplate
}
