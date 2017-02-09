package update

var Flag uint16

//return false,the caller have to unpack the SSU,and inc Flag
func SetFlag()bool{
	if Flag == 0 {
		return false
	}else{
		Flag++
		return true
	}
}
