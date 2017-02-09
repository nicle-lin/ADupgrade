package update

import "sync"

var Flag uint16
var m *sync.RWMutex = new(sync.RWMutex)
//return false,the caller have to unpack the SSU,and inc Flag
func GetFlag()bool{
	m.RLock()
	defer m.RUnlock()
	if Flag == 0 {
		return false
	}else{
		return true
	}
}
//when unpack SSU done, it should call this function
func IncFlag(){
	m.Lock()
	defer m.Unlock()
	Flag++
}

//when upgrade success, it should call this function
func DecFlag(){
	m.Lock()
	defer m.Unlock()
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