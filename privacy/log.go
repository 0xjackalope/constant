package privacy

import "github.com/ninjadotorg/cash/common"

type PrivacyLogger struct {
	log common.Logger
}

func (self *PrivacyLogger) Init(inst common.Logger) {
	self.log = inst
}

// Global instant to use
var Logger = PrivacyLogger{}
