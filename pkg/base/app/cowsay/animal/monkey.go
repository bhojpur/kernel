package animal

type MonkeyTemplate struct {
}

func (mt MonkeyTemplate) Get() string {
	monkeyTemplate := `
	\   
	 \  
	  \
		  .="=.
		_/.-.-.\_     _
	   ( ( o o ) )    ))
		|/  "  \|    //
		 \'---'/    //
		 /'"""'\\  ((
		/ /_,_\ \\  \\
		\_\_'__/  \  ))
		/'  /'~\   |//
	   /   /    \  /
  ,--',--'\/\    /
  '-- "--'  '--'
 `
	return monkeyTemplate
}
