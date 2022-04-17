package animal

type PigTemplate struct {
}

func (pt PigTemplate) Get() string {
	pigTemplate := `
	  \
	   \       
		 _//| .-~~~-.
	   _/oo  }       }-@
	  ('')_  }       |
	   '--'| { }--{  }
			//_/  /_/ 
   `
	return pigTemplate
}
