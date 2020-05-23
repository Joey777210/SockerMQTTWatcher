package command

// Create a new instance of the logger. You can have any number of instances.
//var Mylog = log.New()
//
//func SetMylog() {
//	// The API for setting attributes is a little different than the package level
//	// exported logger. See Godoc.
//
//
//	// You could set this to any `io.Writer` such as a file
//	file, err := os.OpenFile("/var/run/socker/mqttLog", os.O_CREATE|os.O_WRONLY, 0777)
//	if err == nil {
//	 Mylog.Out = file
//	} else {
//	 log.Info("Failed to log to file, using default stderr")
//	}
//
//	Mylog.WithFields(log.Fields{
//		"animal": "walrus",
//		"size":   10,
//	}).Info("A group of walrus emerges from the ocean")
//}
