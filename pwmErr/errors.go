package pwmErr

import "errors"

var InvalidOpt = errors.New("invalid option value")
var FailToCreateRepository = errors.New("fail to create some of repository")
var InsertDB = errors.New("fail to insert data")
var DBInit = errors.New("fail to initialize password storage")
var NoUser = errors.New("not found user")
var GeneratePw = errors.New("fail to generate hashed password")
var Unknown = errors.New("unknown err")
var NotRegisteredKey = errors.New("not registered key")
