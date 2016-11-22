// structure
package SonarRegener

import (
//"fmt"
)

type Sonar struct {
	Header            OIC_Header
	PortSidescan      [5825]int16
	StarboardSidescan [5825]int16
	SubBottom         [6000]int16
}

//initialize sonar struct ,call OIC initialize
func (sonar *Sonar) Init() {
	header := &sonar.Header
	OICInit(header)

}
