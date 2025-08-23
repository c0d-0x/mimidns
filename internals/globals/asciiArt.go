package globals

import "fmt"

func DrawASCIIArt() {
	asciiArt := `
        .__        .__    .___             
  _____ |__| _____ |__| __| _/____   ______
 /     \|  |/     \|  |/ __ |/    \ /  ___/
|  Y Y  \  |  Y Y  \  / /_/ |   |  \\___ \ 
|__|_|  /__|__|_|  /__\____ |___|  /____  >
      \/         \/        \/    \/     \/
`
	fmt.Println(asciiArt)
}
