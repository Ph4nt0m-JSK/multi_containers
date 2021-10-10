package xdd

func init() {
	initDB()
	go func() {
		Save <- &JdCookie{}
	}()
	initContainer()
	//initHandle()
	//intiSky()
}
