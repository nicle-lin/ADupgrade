package update

var Flag uint16

//return false,the caller have to unpack the SSU,and inc Flag
func GetFlag()bool{
	if Flag == 0 {
		return false
	}else{
		return true
	}
}
//when unpack SSU done, it should call this function
func IncFlag(){
	Flag++
}

//when upgrade success, it should call this function
func DecFlag(){
	if Flag > 0{
		Flag--
	}
}


func UnpackSSU(){
	if !GetFlag(){
		IncFlag()
		//don't have to unpack SSU,because it has been unpacked
		return
	}



	IncFlag()
}