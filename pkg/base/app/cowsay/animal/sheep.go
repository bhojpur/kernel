package animal

type SheepTemplate struct {
}

func (st SheepTemplate) Get() string {
	sheepTemplate := `
	 \
	  \  __     
		UooU\.'@@@@@@'.
		\__/(@@@@@@@@@@)
			 (@@@@@@@@)
			 'YY~~~~YY'
			  ||    ||
  `
	return sheepTemplate
}
