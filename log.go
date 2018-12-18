package kademlia

import "log"

const (
	logLevelDebug = iota
	logLevelInfo
	logLevelDisable
)

const (
	debugPrefix = "[Kademlia:DEBUG] "
	infoPrefix  = "[Kademlia:INFO] "
)

var kadlog = logger{logLevelInfo}

type logger struct {
	loglevel int
}

func (l *logger) debug(v ...interface{}) {
	if l.loglevel <= logLevelDebug {
		log.Print(debugPrefix, v)
	}
}

func (l *logger) debugln(v ...interface{}) {
	if l.loglevel <= logLevelDebug {
		log.Print(debugPrefix, v)
	}
}

func (l *logger) debugf(format string, v ...interface{}) {
	if l.loglevel <= logLevelDebug {
		log.Print(debugPrefix+format, v)
	}
}

func (l *logger) info(v ...interface{}) {
	if l.loglevel <= logLevelInfo {
		log.Print(infoPrefix, v)
	}
}

func (l *logger) infoln(v ...interface{}) {
	if l.loglevel <= logLevelInfo {
		log.Print(infoPrefix, v)
	}
}

func (l *logger) infof(format string, v ...interface{}) {
	if l.loglevel <= logLevelInfo {
		log.Print(infoPrefix+format, v)
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
