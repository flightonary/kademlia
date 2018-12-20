package kademlia

import "log"

const (
	logLevelDebug = iota
	logLevelInfo
	logLevelDisable
)

const (
	debugPrefix = "[Kademlia:DEBUG]"
	infoPrefix  = "[Kademlia:INFO]"
)

var kadlog = logger{logLevelInfo}

type logger struct {
	loglevel int
}

func (l *logger) debug(v ...interface{}) {
	if l.loglevel <= logLevelDebug {
		v = append(v[0:1], v[0:]...)
		v[0] = debugPrefix
		log.Print(v...)
	}
}

func (l *logger) debugf(format string, v ...interface{}) {
	if l.loglevel <= logLevelDebug {
		log.Printf(debugPrefix+format, v...)
	}
}

func (l *logger) info(v ...interface{}) {
	if l.loglevel <= logLevelInfo {
		v = append(v[0:1], v[0:]...)
		v[0] = infoPrefix
		log.Print(v...)
	}
}

func (l *logger) infof(format string, v ...interface{}) {
	if l.loglevel <= logLevelInfo {
		log.Printf(infoPrefix+format, v...)
	}
}

func SetLogLevelDebug() {
	kadlog.loglevel = logLevelDebug
}

func SetLogLevelInfo() {
	kadlog.loglevel = logLevelInfo
}

func DisableLog() {
	kadlog.loglevel = logLevelDisable
}
